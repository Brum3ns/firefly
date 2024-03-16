package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Brum3ns/firefly/internal/config"
	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/internal/knowledge"
	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/internal/scan"
	"github.com/Brum3ns/firefly/internal/ui"
	"github.com/Brum3ns/firefly/internal/verbose"
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/files"
	"github.com/Brum3ns/firefly/pkg/insertpoint"
	"github.com/Brum3ns/firefly/pkg/payloads"
	"github.com/Brum3ns/firefly/pkg/prepare"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/statistics"
	"github.com/Brum3ns/firefly/pkg/waitgroup"
)

// The runner should contain the structures needed for all the processes.
// It must NOT contain structures that need to be modified and/or dynamicly changed once the process is running.
type Runner struct {
	Count           int
	OutputOK        bool
	VerifyMode      bool
	TerminalUIMode  bool
	Conf            *config.Configure
	Design          *design.Design
	RequestTasks    *request.TaskStorage
	stats           statistics.Statistic
	channel         Channel
	handler         Handler
	StoredKnowledge map[string]knowledge.Knowledge
}

type Handler struct {
	HTTP    request.Handler
	Scanner scan.Handler
}

type Channel struct {
	ListenerScanner chan scan.Result
	ListenerHTTP    chan request.Result
	Result          chan output.ResultFinal
	Statistic       chan bool
}

// Setup a new runner. The runner can run in a verification mode, in that case the argument "knowledgeStorage" MUST be set to "nil".
// The other mode is the attack mode and need the "knowledgeStorage" to contain knowledge (data) about the target to attack to be run successfully.
func NewRunner(conf *config.Configure, knowledgeStorage map[string]knowledge.Knowledge) *Runner {
	var verifyMode = (knowledgeStorage == nil)
	return &Runner{
		Count:          0,
		Conf:           conf,
		VerifyMode:     verifyMode,
		TerminalUIMode: (!verifyMode && conf.TerminalUI),
		OutputOK:       (len(conf.Output) > 0 && knowledgeStorage != nil),
		Design:         design.NewDesign(),
		stats:          statistics.NewStatistic(verifyMode),
		channel: Channel{
			ListenerScanner: make(chan scan.Result),
			ListenerHTTP:    make(chan request.Result),
			Result:          make(chan output.ResultFinal),
			Statistic:       make(chan bool),
		},
		handler: Handler{
			// Setup the HTTP handler:
			HTTP: request.NewHandler(request.HandlerSettings{
				Delay:      conf.Delay,
				Threads:    conf.Threads,
				VerifyMode: verifyMode,
				Client: request.NewClient(request.ClientSettings{
					Timeout: conf.Timeout,
					Proxy:   conf.Proxy,
					HTTP2:   conf.HTTP2,
				}),
				RequestBase: request.RequestBase{
					RandomUserAgent:      conf.RandomAgent,
					HeadersOriginalArray: conf.Headers,
					PostBody:             conf.PostData,
					InsertPoint:          conf.InsertKeyword,
				},
			}),

			// Setup the HTTP scanner handler:
			Scanner: scan.NewHandler(scan.Settings{
				Scanner:       conf.Scanner,
				Threads:       conf.ThreadsScanner,
				PayloadVerify: conf.Options.VerifyPayload,
				Knowledge:     knowledgeStorage,
			}),
		},
		StoredKnowledge: knowledgeStorage,
	}
}

// Firefly verify/fuzz runner
// The runner is the core process for all other child processes. It's preforming the requests and listen for HTTP results to be scanned analyzed.
func (r *Runner) Run() (map[string]knowledge.Knowledge, statistics.Statistic, error) {
	var (
		outputFileWriter = r.MustValidateOutput()
		learnt           = make(map[string][]knowledge.Learnt)
		display          = output.NewDisplay(r.Design)
		terminalUI       = ui.NewProgram()
		wg               waitgroup.WaitGroup
	)

	// Start terminal UI
	if r.TerminalUIMode {
		wg.Add(1)
		go func() {
			if _, err := terminalUI.Run(); err != nil {
				log.Fatalf("terminal UI - %s", err)
			}
			wg.Done()
		}()
	}

	// Start the request and scanner handlers
	go r.handler.HTTP.Run(r.channel.ListenerHTTP)
	go r.handler.Scanner.Run(r.channel.ListenerScanner)

	//Runner listener
	go func() {
		var (
			progressbar = ui.NewProgressBar(100, &r.stats)
			//progressBar = statistics.NewProgressBar(100, &r.stats)
			mutex sync.Mutex
		)
		for {
			select {
			case <-r.channel.Statistic:
				if !r.Conf.NoDisplay && r.TerminalUIMode {
					terminalUI.Send(r.stats)
				}

			case result := <-r.channel.Result:
				r.stats.Count()

				if r.VerifyMode {
					mutex.Lock()
					learnt[result.TargetHashId] = append(learnt[result.TargetHashId], knowledge.Learnt{
						Payload:  result.Payload,
						Extract:  result.Scanner.Extract,
						HTMLNode: prepare.GetHTMLNode(result.Response.Body),
						Response: result.Response,
					})
					mutex.Unlock()
				} else if result.UnkownBehavior {
					r.stats.Behavior.Count()

					// Send the result to the output file specified by the user:
					if r.OutputOK {
						err := output.WriteJSON(r.stats.Output.GetCount(), outputFileWriter, result)
						if err != nil {
							log.Println(design.STATUS.ERROR, "Request ID:", result.RequestId, err)
						}
						r.stats.Output.Count()
					}

					// Display the final result to the screen (CLI)
					if !r.Conf.NoDisplay {
						if r.TerminalUIMode {
							terminalUI.Send(r.stats)
							terminalUI.Send(result)
						} else {
							display.ToScreen(result, r.Conf.TerminalUI)
							progressbar.Print()
						}
					}
				}
			}
		}
	}()

	//Listeners
	go r.listenerScanner()
	go r.listenerHTTP()

	// Give all the request jobs to the HTTP handler and wait until the handlers are completed with all the jobs:
	jobHandlerAmount := r.jobToHandler(&r.handler.HTTP)
	r.waitForHandlers(jobHandlerAmount)

	// Close the output file (if any output have been handled):
	if r.OutputOK {
		if err := outputFileWriter.Close(); err != nil {
			log.Fatal(err)
		}
	}

	if r.TerminalUIMode {
		terminalUI.Quit()
		wg.Wait()
	}

	return knowledge.GetKnowledge(learnt), r.stats, nil
}

// Listen for results from the HTTP handler and preform a scan for each intercepted HTTP result:
func (r *Runner) listenerScanner() {
	for {
		scanResult := <-r.channel.ListenerScanner
		if scanResult.Error != nil {
			verbose.Show(scanResult.Error)
		} else {
			r.stats.Scanner.Count()
			r.channel.Result <- scanResult.Output
		}
	}
}

// Listen for HTTP request/response results from the request handler and add the response as a job to the scanner handler:
func (r *Runner) listenerHTTP() {
	var storedKnowledge = knowledge.Knowledge{}
	for {
		resultHTTP := <-r.channel.ListenerHTTP
		r.stats.Request.Count()

		//Check if we got a valid HTTP response from our requested target or if any error appeared:
		if resultHTTP.Error != nil {
			r.stats.Response.CountError()
			r.channel.Statistic <- true
			verbose.Show(resultHTTP.Error)
			continue
		}
		r.stats.Response.Count()
		r.stats.Response.UpdateTime(resultHTTP.Response.Time)

		//Filter the HTTP response (if set):
		if r.Conf.Filter.Run(resultHTTP.Response) {
			r.stats.Response.CountFilter()
			r.channel.Statistic <- true
			continue
		}

		//Give the scanner handler job related to the Http result (request/response):
		go func() {
			if !r.VerifyMode {
				storedKnowledge = r.StoredKnowledge[resultHTTP.TargetHashId]
			}
			// Note: payload included in the HTTP result:
			r.handler.Scanner.AddJob(r.VerifyMode, resultHTTP, storedKnowledge)
		}()
	}
}

// Validate and verify the output to store the result to (if set):
// Note : (will panic in case an error is triggered)
func (r *Runner) MustValidateOutput() *os.File {
	var (
		fileWriter = &os.File{}
		err        error
	)
	//Create output file and create a file writer (*if output file set*):
	if r.OutputOK {
		if !files.FileExist(r.Conf.Output) || r.Conf.Overwrite {
			fileWriter, err = os.OpenFile(r.Conf.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Panicln(err)

			} else if err = fileWriter.Truncate(0); err != nil {
				log.Panicln(err)

			} else if _, err = fileWriter.Seek(0, 0); err != nil {
				log.Panicln(err)
			}
		} else {
			err = fmt.Errorf("%s The specified output file already exists (\033[33m%s\033[0m), use the overwrite option to overwrite it", design.STATUS.FAIL, r.Conf.Output)
			log.Panicln(err)
		}
		verbose.Show("Save result to output file: " + r.Conf.Output)
	}
	return fileWriter
}

func (r *Runner) jobToHandler(requestHandler *request.Handler) int {
	var (
		payloadWordlist = r.Conf.Wordlist.GetAll()
		headersArray    = r.Conf.Headers
		postbody        = r.Conf.PostData
		jobAmount       = 0
	)
	for hash, host := range r.Conf.Hosts {
		param := r.Conf.Params[hash]
		rawURL := host.URL

		for _, tag := range payloads.TAGS {
			// Check if we should adapt to "behavior verification mode":
			if (r.VerifyMode && tag != payloads.TAG_VERIFY) || (!r.VerifyMode && tag == payloads.TAG_VERIFY) {
				continue
			}

			wordlist := payloadWordlist[tag]
			for _, payload := range wordlist {
				// Prepare the request by inserting the current payload into the request:
				// !Note : (Some variables given will be modified)
				insert := insertpoint.NewInsert(r.Conf.InsertKeyword, payload)

				URLStruct, _ := url.Parse(rawURL)

				if param.AutoQueryURL {
					URLStruct.RawQuery = param.URL.RawQueryInsertPoint
					rawURL = URLStruct.String()
				}

				if param.AutoQueryBody {
					postbody = param.Body.RawQueryInsertPoint
				}

				if param.AutoQueryCookie {
					headersArray = request.SetNewHeaderValue(headersArray, "cookie", param.Cookie.RawQueryInsertPoint)
				}

				randomUserAgents, err := getRandomUserAgent(global.FILE_RANDOMAGENT)
				if err != nil {
					log.Fatalf("Random User-Agent:", err)
				}

				requestSettings := request.RequestSettings{
					UserAgents:   randomUserAgents,
					TargetHashId: hash,
					Tag:          tag,
					Payload:      payload,
					URLOriginal:  rawURL,
					Parameter:    r.Conf.Params[hash],
					URL:          insert.SetURL(rawURL),
					Method:       insert.SetMethod(host.Method),
					RequestBase: request.RequestBase{
						Headers:              insert.SetHeaders(headersArray),
						PostBody:             insert.SetPostBody(postbody),
						RandomUserAgent:      r.Conf.RandomAgent,
						HeadersOriginalArray: r.Conf.Headers,
					},
				}
				jobAmount++
				requestHandler.AddJob(requestSettings)
			}
		}
	}
	return jobAmount
}

// Take a file containing user agents
func getRandomUserAgent(file string) ([]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("User-Agent file error :", err)
	}
	return strings.Split(string(content), "\n"), nil
}

func (r *Runner) waitForHandlers(jobHandlerAmount int) {
	for {
		time.Sleep(100 * time.Millisecond)
		if jobHandlerAmount > 0 && jobHandlerAmount == r.handler.HTTP.GetJobAmount() && r.handler.HTTP.GetInProcess() == 0 {
			break
		}
	}
}
