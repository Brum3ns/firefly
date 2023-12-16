package runner

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/files"
	"github.com/Brum3ns/firefly/pkg/firefly/config"
	"github.com/Brum3ns/firefly/pkg/firefly/verbose"
	"github.com/Brum3ns/firefly/pkg/knowledge"
	"github.com/Brum3ns/firefly/pkg/output"
	"github.com/Brum3ns/firefly/pkg/payloads"
	"github.com/Brum3ns/firefly/pkg/prepare"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/scan"
	"github.com/Brum3ns/firefly/pkg/statistics"
)

// The runner should contain the structures needed for all the processes.
// It must NOT contain structures that need to be modified and/or dynamicly changed once the process is running.
type Runner struct {
	Count           int
	VerifyMode      bool
	OutputOK        bool
	Conf            *config.Configure
	Design          *design.Design
	RequestTasks    *request.TaskStorage
	statistic       statistics.Statistic
	channel         Channel
	handler         Handler
	StoredKnowledge map[string]knowledge.Knowledge
}

type Handler struct {
	HTTP    request.Handler
	Scanner scan.Handler
}

type Channel struct {
	Skip            chan statistics.Skip
	ListenerScanner chan scan.Result
	ListenerHTTP    chan request.Result
	ResultScanner   chan scan.Result
	ResultHTTP      chan request.Result
}

// Setup a new runner. The runner can run in a verification mode, in that case the argument "knowledgeStorage" MUST be set to "nil".
// The other mode is the attack mode and need the "knowledgeStorage" to contain knowledge (data) about the target to attack to be run successfully.
func NewRunner(conf *config.Configure, knowledgeStorage map[string]knowledge.Knowledge) *Runner {
	var verifyMode = (knowledgeStorage == nil)
	return &Runner{
		Count:      0,
		Conf:       conf,
		VerifyMode: verifyMode,
		OutputOK:   (len(conf.Output) > 0 && knowledgeStorage != nil),
		Design:     design.NewDesign(),
		statistic:  statistics.NewStatistic(verifyMode),
		channel: Channel{
			Skip:            make(chan statistics.Skip),
			ListenerScanner: make(chan scan.Result),
			ListenerHTTP:    make(chan request.Result),
			ResultScanner:   make(chan scan.Result),
			ResultHTTP:      make(chan request.Result),
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
	)

	// Start the request handler and send requests to the target:
	go r.handler.HTTP.Run(r.channel.ListenerHTTP)

	// Start the scanner in the background:
	go r.handler.Scanner.Run(r.channel.ListenerScanner)

	// [Handler - Scanner]
	// Listen for results from the HTTP handler and preform a scan for each intercepted HTTP result:
	go func() {
		var (
			mutex       sync.Mutex
			progressBar = statistics.NewProgressBar(101, &r.statistic)
		)

		for {
			scanResult := <-r.channel.ListenerScanner
			if scanResult.Error != nil {
				verbose.Show(scanResult.Error)
				r.statistic.AddScannerError(scanResult.Error)
				continue
			}
			r.statistic.CountScanner()

			//Collect and store verify data (verification mode):
			if r.VerifyMode {
				mutex.Lock()
				learnt[scanResult.Output.TargetHashId] = append(learnt[scanResult.Output.TargetHashId], knowledge.Learnt{
					Payload:  scanResult.Output.Payload,
					Extract:  scanResult.Output.Scanner.Extract,
					HTMLNode: prepare.GetHTMLNode(scanResult.Output.Response.Body),
					Response: scanResult.Output.Response,
				})
				mutex.Unlock()

				//Analyze if the result is an unkown behavior:
			} else if scanResult.UnkownBehavior {
				r.statistic.CountBehavior()

				//Send the result to the output file specified by the user:
				if r.OutputOK {
					err := output.WriteJSON(r.statistic.Output, outputFileWriter, scanResult.Output)
					if err != nil {
						log.Println(design.STATUS.ERROR, "Request ID:", scanResult.Output.RequestId, err)
					}
					r.statistic.CountOutput()
				}

				//Display the final result to the screen (CLI)
				if !r.Conf.NoDisplay {
					display.ToScreen(scanResult.Output)
				}
			}
			r.statistic.CountComplete()

			progressBar.Print()
		}
	}()

	// [Handler - HTTP]
	// Intercept HTTP request/response results from the request handler and add the response as a job to the scanner handler:
	go func() {
		var storedKnowledge = knowledge.Knowledge{}
		for {
			resultHTTP := <-r.channel.ListenerHTTP
			r.statistic.CountRequest()

			//Check if we got a valid HTTP response from our requested target or if any error appeared:
			if resultHTTP.Error != nil {
				r.statistic.CountError()
				verbose.Show(resultHTTP.Error)
				continue
			}

			r.statistic.CountResponse()

			//Filter the HTTP response (if set):
			if r.Conf.Filter.Run(resultHTTP.Response) {
				r.statistic.CountFilter()
				continue
			}

			//Give the scanner handler job related to the Http result (request/response):
			go func() {
				if !r.VerifyMode {
					storedKnowledge = r.StoredKnowledge[resultHTTP.TargetHashId]
				}
				r.handler.Scanner.AddJob(r.VerifyMode, resultHTTP, storedKnowledge)
			}()
		}
	}()

	// Give all the request job to the HTTP handler and wait until the handlers are completed with all the jobs:
	jobHandlerAmount := r.jobToHandler(&r.handler.HTTP)
	r.waitForHandlers(jobHandlerAmount)

	//Close the output file (if any output have been handled):
	if r.OutputOK {
		if err := outputFileWriter.Close(); err != nil {
			log.Fatal(err)
		}
	}

	return knowledge.GetKnowledge(learnt), r.statistic, nil
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
				insert := request.NewInsert(r.Conf.InsertKeyword, payload)

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

				requestSettings := request.RequestSettings{
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

func (r *Runner) waitForHandlers(jobHandlerAmount int) {
	for {
		time.Sleep(30 * time.Millisecond)
		if jobHandlerAmount > 0 && jobHandlerAmount == r.handler.HTTP.GetJobAmount() && r.handler.HTTP.GetInProcess() == 0 {
			break
		}
	}
}
