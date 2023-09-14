package scan

import (
	"fmt"
	"sync"

	"github.com/Brum3ns/firefly/pkg/firefly/config"
	"github.com/Brum3ns/firefly/pkg/output"
	"github.com/Brum3ns/firefly/pkg/payloads"
	"github.com/Brum3ns/firefly/pkg/prepare"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/scan/difference"
	"github.com/Brum3ns/firefly/pkg/scan/extract"
	"github.com/Brum3ns/firefly/pkg/scan/transformation"
	"github.com/Brum3ns/firefly/pkg/verify"
)

type Engine struct {
	Process   process
	WaitGroup sync.WaitGroup
	JobQueue  chan Job
	Pool      chan chan Job
	quit      chan bool
	Properties
}

// Properties given by user input to adapt the scanning process
type Properties struct {
	Threads       int
	PayloadVerify string
	//VerifiedStorage map[string][]verify.TargetKnowledge //!Note : (This map *MUST* be static and not modifed)

	//The scanner contains the points to a base structure that contains the base structure
	//of all the scanner techniques the engine need. This save memory and gain better preformence in the overall preformance.
	//Note : (Static data stored. Read struct DESC)
	Scanner *config.Scanner
}

// process represents the process that executes the job
type process struct {
	PayloadVerify string
	jobChannel    chan Job
	pool          chan chan Job
	//Knowledge  map[string][]verify.TargetKnowledge //Note : (Should be a pointer of "Properties.VerifyStorage")
	Scanner *config.Scanner //!Note : (Static data stored. Read struct DESC)
	Result  processResult   //<-Returned
}

type processResult struct {
	UnkownBehavior bool
	Http           request.Result
	Extract        extract.Result
	Difference     difference.Result
	Transformation transformation.Result
}

type Job struct {
	Knowledge []verify.TargetKnowledge
	Http      request.Result
}

// Note : (Alias of structure "output.ResultFinal")
type Result struct {
	Output         output.ResultFinal
	Error          error
	UnkownBehavior bool
}

// Start the handler for the workers by giving the tasks to preform and the amount of workers.
func NewEngine(properties Properties) *Engine {
	return &Engine{
		Properties: properties,
		JobQueue:   make(chan Job),
		Pool:       make(chan chan Job, properties.Threads),
	}
}

// Start all the processes and assign tasks (jobs) to the scanners that are listening. Use the method "Stop()" to stop the scanner.
// Note : (The scanner engine *MUST* run inside a [go]rutine. It can only stop from the method "Stop()" that do send a stop signal to the engine)
func (e *Engine) Run(listener chan<- Result) {
	var pResult = make(chan processResult)

	//Validate process amount:
	if e.Threads <= 0 {
		e.Threads = 1
	}

	//Start the amount of processes related to the amount of given threads:
	for i := 0; i < e.Threads; i++ {
		e.Process = newProcess(e.PayloadVerify, e.Scanner, e.Pool)
		e.Process.spawnProcess(pResult)
	}

	//Listen for new jobs from the queue and send it to the job channel for the workers to handle it:
	go func() {
		for {
			select {
			case job := <-e.JobQueue:
				go func(job Job) {
					//Get an available job channel from any running process:
					jobChannel := <-e.Pool

					//Give the available process the job:
					jobChannel <- job
				}(job)

				//Listen for result from any process, if a result is recived, then send it to the listener [chan]nel:
			case r := <-pResult:
				//Detect unkown/new behavior
				r.UnkownBehavior = detectBehavior(r)

				listener <- makeResult(r)
				e.WaitGroup.Done()
			}
		}
	}()

	//Listen a stop signal then wait until all background processes are completed:
	if <-e.quit {
		e.WaitGroup.Wait()
		fmt.Println(":: Scanner Engine stopped")
		return
	}
}

// Add new jobs (tasks) to be performed by the engine processes:
func (e *Engine) AddJob(verifyMode bool, httpResult request.Result, knowledge []verify.TargetKnowledge) {
	e.WaitGroup.Add(1)
	e.JobQueue <- Job{
		Http:      httpResult,
		Knowledge: knowledge,
	}
}

func (e *Engine) Stop() {
	e.quit <- true
}

// Create a new process
func newProcess(payloadVerify string, scanner *config.Scanner, pool chan chan Job) process {
	return process{
		PayloadVerify: payloadVerify,
		pool:          pool,
		Scanner:       scanner,
		//Knowledge:  verifiedStorage,
		jobChannel: make(chan Job),
	}
}

// Spawn a new process
func (p process) spawnProcess(result chan processResult) {
	go func() {
		for {
			// Add the current spawned process into the process queue:
			p.pool <- p.jobChannel

			//A job was given, start processing it
			select {
			case job := <-p.jobChannel:
				result <- p.start(job)
			}
		}
	}()
}

// Start a new process
func (p process) start(job Job) processResult {
	var (
		wg                   sync.WaitGroup
		extResult            extract.Result
		diffResult           difference.Result
		transformationResult transformation.Result
	)
	//Check if we should preform scanner techniques or not:
	if !p.Scanner.DisablesTechniques {

		//Transformation process:
		if job.Http.Tag == payloads.TAG_TRANSFORMATION && p.Scanner.OK_Transformation {
			wg.Add(1)
			go func() {
				tfmt := p.Scanner.Transformation
				transformationResult = tfmt.Detect(job.Http.Response.Body, job.Http.Payload)
				wg.Done()
			}()
		}

		//Extract process:
		if p.Scanner.OK_Extract {
			wg.Add(1)
			go func() {
				ext := p.Scanner.Extract
				ext.AddJob(
					job.Http.Response.Body,
					job.Http.Response.HeaderString,
					p.Scanner.Extract.Known,
				)
				extResult = ext.Run()
				wg.Done()
			}()
		}

		//[Diff]erence process:
		if p.Scanner.OK_Diff && job.Knowledge != nil {
			wg.Add(1)
			go func() {
				for _, storage := range job.Knowledge {
					if job.Http.Response.Body == storage.Response.Body {
						wg.Done()
						return
					}
				}
				//Make a new difference instant and provided the current HTTP response body and headers:
				diff := difference.NewDifference(
					difference.Properties{
						Payload:       job.Http.Payload,
						PayloadVerify: p.PayloadVerify,
						ResponseBody: &difference.ResponseBody{
							Body:     job.Http.Response.Body,
							HtmlNode: prepare.GetHTMLNode(job.Http.Response.Body),
						},
						ResponseHeaders: &difference.ResponseHeaders{
							HeaderString: job.Http.Response.HeaderString,
							Headers:      job.Http.Response.Header,
						},
					},
				)
				//Add all the known storage
				for _, knownStorage := range job.Knowledge {
					diff.AppendKnownHTMLNode(knownStorage.HTMLNode)
					diff.AppendKnownHeaders(knownStorage.Response.Headers)
				}
				diffResult = diff.Run()
				wg.Done()
			}()
		}

		//Wait for all the scanners to finish:
		wg.Wait()
	}

	return processResult{
		Http:           job.Http,
		Extract:        extResult,
		Difference:     diffResult,
		Transformation: transformationResult,
	}
}

// Start the extract scanning process
func makeResult(r processResult) Result {
	req := r.Http.Request
	resp := r.Http.Response

	return Result{
		UnkownBehavior: r.UnkownBehavior,
		Output: output.ResultFinal{
			TargetHashId: r.Http.TargetHashId,
			RequestId:    r.Http.RequestId,
			Tag:          r.Http.Tag,
			Date:         r.Http.Date,
			Payload:      r.Http.Payload,
			OK:           true,

			Request: output.Request{
				URL:         req.RequestURI,
				URLOriginal: req.URLOriginal,
				Host:        req.URL.Host,
				Scheme:      req.URL.Scheme,
				Method:      req.Method,
				PostBody:    req.Body,
				Proto:       req.Proto,
				Headers:     req.HeadersOriginal,
			},
			Response: output.Response{
				Time:          resp.Time,
				Host:          resp.Request.Host,
				Body:          resp.Body,
				Title:         resp.Title,
				Proto:         resp.Proto,
				IPAddress:     resp.IPAddress,
				StatusCode:    resp.StatusCode,
				WordCount:     resp.WordCount,
				LineCount:     resp.LineCount,
				ContentType:   resp.ContentType,
				ContentLength: int(resp.ContentLength),
				HeaderAmount:  len(resp.Header),
				Headers:       resp.Header,
			},
			Scanner: output.Scanner{
				Extract:        r.Extract,
				Diff:           r.Difference,
				Transformation: r.Transformation,
				//Data...
			},

			Error: nil,
		},
		Error: nil,
	}
}
