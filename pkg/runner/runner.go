package runner

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/files"
	"github.com/Brum3ns/firefly/pkg/firefly/config"
	"github.com/Brum3ns/firefly/pkg/firefly/verbose"
	"github.com/Brum3ns/firefly/pkg/output"
	"github.com/Brum3ns/firefly/pkg/prepare"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/scan"
	"github.com/Brum3ns/firefly/pkg/verify"
)

// The runner should contain the structures needed for all the processes.
// It must NOT contain structures that need to be modified and/or dynamicly changed once the process is running.
type Runner struct {
	Count            int
	VerifyMode       bool
	OutputOK         bool
	Conf             *config.Configure
	Design           *design.Design
	RequestTasks     *request.TaskStorage
	ClientProperties *request.ClientProperties
	Knowledge        Knowledge //<-Returned
}
type skipProcess struct {
	tag string
	err error
}

// The Knowledge structure contains variables that will be modified during the runner process
// Note : (Runner's return data)
type Knowledge struct {
	Storage map[string][]verify.TargetKnowledge //(Original URL|Data of Verified behavior)
}

// Firefly verify/fuzz runner
// The runner is the core process for all other child processes. It's preforming the requests and listen for responses from the target.
// When a response has been collected it's sent to the engine that handle the hardware processes. It do so by spinning up a task for
// each analyze process (aka: tasks).
func Run(conf *config.Configure, knowledgeStorage map[string][]verify.TargetKnowledge) (map[string][]verify.TargetKnowledge, error) {
	runner := &Runner{
		Count:      0,
		Conf:       conf,
		VerifyMode: (knowledgeStorage == nil),
		OutputOK:   (len(conf.Output) > 0 && knowledgeStorage != nil),
		Design:     design.NewDesign(),
		ClientProperties: &request.ClientProperties{
			Timeout: conf.Timeout,
			Proxy:   conf.Proxy,
			HTTP2:   conf.HTTP2,
		},
		RequestTasks: &request.TaskStorage{
			URLs:            conf.URLs,
			Schemes:         conf.Scheme,
			Methods:         conf.Methods,
			Headers:         conf.Headers,
			PostData:        conf.PostData,
			InsertPoint:     conf.Insert,
			RandomUserAgent: conf.RandomAgent,
			Payloads:        conf.Wordlist.GetAll(),
		},
		//Check if we already have verified data stored:
		Knowledge: Knowledge{
			Storage: verify.Prepare(knowledgeStorage),
		},
	}

	//Pre-Declare the output file and JSON encoder that will be used within the output process: (if set):
	var (
		outputFileWriter = &os.File{}
		err              error
	)

	//Create output file and create a file writer (*if output file set*):
	if runner.OutputOK {
		if !files.FileExist(runner.Conf.Output) || runner.Conf.Overwrite {
			outputFileWriter, err = os.OpenFile(runner.Conf.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			if err = outputFileWriter.Truncate(0); err != nil {
				log.Fatal(err)
			}
			if _, err = outputFileWriter.Seek(0, 0); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("%s The specified output file already exists (\033[33m%s\033[0m), use the overwrite option to overwrite it.", design.STATUS.FAIL, runner.Conf.Output)
		}
		verbose.Show("Save result to output file: " + runner.Conf.Output)
	}

	var (
		//Setup the [chan]nels and waitgroups for all the processes that will be preformed by the runner:
		done            = make(chan bool)
		doneRequests    = make(chan bool)
		RequestAmount   = make(chan int)
		skip            = make(chan skipProcess)
		listenerScanner = make(chan scan.Result) //<- (listen for the scanner result for each HTTP response)
		InterceptHTTP   = make(chan request.Result)
	)

	//Track the statistic of the runners core and nested processes:
	stats := newStatistic()
	go func() {
		stats.jobTotal = <-RequestAmount
	}()

	//Start the request handler and send requests to the target:
	requestHandler := request.NewHandler(runner.ClientProperties, runner.RequestTasks, runner.Conf.Threads, runner.Conf.Delay, runner.VerifyMode)
	go func() {
		requestHandler.Run(InterceptHTTP, RequestAmount, doneRequests)
		close(InterceptHTTP)
		close(doneRequests)
	}()

	//[Listener] : Listen for the scanner result and make the final result:
	go func() {
		var mutex sync.Mutex
		for stats.inProcess() {
			//TODO - Loadingbar

			select {
			case scanResult := <-listenerScanner:
				if scanResult.Error != nil {
					stats.appendScannerError(scanResult.Error)
					continue
				}
				stats.countScanner()

				//Collect and store verify data (verification mode):
				if runner.VerifyMode {
					runner.AppendKnowledge(&mutex, scanResult)

					//Analyze if the result is an unkown behavior:
				} else if true {
					//Send the result to the output file specified by the user:
					if runner.OutputOK {
						err := output.WriteJSON(stats.Output, outputFileWriter, scanResult.Output)
						if err != nil {
							log.Println(design.STATUS.ERROR, "Request ID:", scanResult.Output.RequestId, err)
						}
						stats.countOutput()
					}
					//Display the final result to the screen (CLI)
					if !runner.Conf.NoDisplay {
						output.DisplayCLI(stats.Completed, runner.Design, scanResult.Output)
					}
				}
				stats.countComplete()

			case s := <-skip:
				stats.handleSkipped(s)
				if s.err != nil {
					runner.verbose(s.err)
				}
			}
		}
		//Send signal that all the runners process are completed:
		done <- true
	}()

	// Prepare a new scanner engine with all the base properties:
	scanEngine := scan.NewEngine(scan.Properties{
		Scanner:       runner.Conf.Scanner,
		Threads:       runner.Conf.ThreadsEngine,
		PayloadVerify: conf.Options.VerifyPayload,
	})
	//Start the scanner in the background:
	go scanEngine.Run(listenerScanner)

	//Intercept HTTP request/response results from the request handler and add the response as a job to the scanner engine:
	for sendRequest := true; sendRequest; {
		select {
		case HttpResult := <-InterceptHTTP:
			stats.countRequest()

			//Check if any error appeared in the request/response process:
			if HttpResult.Error != nil {
				skip <- skipProcess{
					tag: "error",
					err: HttpResult.Error,
				}
				break

				//Filter the HTTP response (if set):
			} else if conf.Filter.Run(HttpResult.Response) {
				skip <- skipProcess{
					tag: "filter",
				}
				break
			}

			//Give the scanner engine job related to the Http result (request/response):
			go func() {
				//Add job to the scanner engine that runs in the background:
				if runner.VerifyMode {
					scanEngine.AddJob(runner.VerifyMode, HttpResult, nil)
				} else {
					scanEngine.AddJob(runner.VerifyMode, HttpResult, knowledgeStorage[HttpResult.TargetHashId])
				}
			}()

		case <-doneRequests:
			sendRequest = false
		}
	}

	//Wait until the runner is completed:
	<-done

	//Close the output file (if any output have been handled):
	if runner.OutputOK {
		if err := outputFileWriter.Close(); err != nil {
			log.Fatal(err)
		}
	}

	return runner.Knowledge.Storage, nil
}

func (r *Runner) verbose(msg any) {
	if !r.Conf.Verbose {
		return
	}
	icon := r.Design.INFO
	switch msg.(type) {
	case error:
		icon = r.Design.ERROR
	}
	fmt.Fprintln(os.Stderr, icon, fmt.Sprintf("%v", msg))
}

func (r *Runner) AppendKnowledge(mutex *sync.Mutex, rs scan.Result) {
	verifyStorage := verify.TargetKnowledge{
		Payload:  rs.Output.Payload,
		HTMLNode: prepare.GetHTMLNode(rs.Output.Response.Body),
		Response: rs.Output.Response,
		//Request:  rs.Output.Request,
		//Behavior: rs.Output.Behavior,
		//Scanner:  rs.Output.Scanner,
	}

	mutex.Lock()
	r.Knowledge.Storage[rs.Output.TargetHashId] = append(r.Knowledge.Storage[rs.Output.TargetHashId], verifyStorage)
	mutex.Unlock()
}
