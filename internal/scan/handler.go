package scan

import (
	"fmt"
	"log"
	"sync"

	"github.com/Brum3ns/firefly/internal/config"
	"github.com/Brum3ns/firefly/internal/knowledge"
	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/difference"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/payloads"
	"github.com/Brum3ns/firefly/pkg/prepare"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/transformation"
)

type Handler struct {
	Process   process
	WaitGroup sync.WaitGroup
	JobQueue  chan Job
	Pool      chan chan Job
	quit      chan bool
	Settings
}

// Settings given by user input to adapt the scanning process
type Settings struct {
	Threads       int
	PayloadVerify string

	//The scanner contains the points to a base structure that contains the base structure
	//of all the scanner techniques the handler need. This save memory and gain better preformence in the overall preformance.
	//Note : (Static data stored. Read struct DESC)
	Scanner *config.Scanner

	//This map holds all the knowledge of all the targets
	//!Note : (This map *MUST* be static and not modifed)
	Knowledge map[string]knowledge.Knowledge
}

// process represents the process that executes the job
type process struct {
	jobChannel chan Job
	pool       chan chan Job
	Scanner    *config.Scanner //!Note : (Static data stored. Read struct DESC)
	Result     processResult   //<-Returned

	//Knowledge  map[string][]verify.TargetKnowledge //Note : (Should be a pointer of "Properties.VerifyStorage")
	Knowledge map[string]knowledge.Knowledge
}

type processResult struct {
	UnkownBehavior bool
	Http           request.Result
	Extract        extract.Result
	Difference     difference.Result
	Transformation transformation.Result
}

type Job struct {
	Knowledge knowledge.Knowledge
	Http      request.Result
}

// Note : (Alias of structure "output.ResultFinal")
type Result struct {
	Output output.ResultFinal
	Error  error
}

// Start the handler for the workers by giving the tasks to preform and the amount of workers.
func NewHandler(properties Settings) Handler {
	return Handler{
		Settings: properties,
		JobQueue: make(chan Job),
		Pool:     make(chan chan Job, properties.Threads),
	}
}

// Start all the processes and assign tasks (jobs) to the scanners that are listening. Use the method "Stop()" to stop the scanner.
// Note : (The scanner handler *MUST* run inside a [go]rutine. It can only stop from the method "Stop()" that do send a stop signal to the handler)
func (e *Handler) Run(listener chan<- Result) {
	var pResult = make(chan processResult)

	//Validate process amount:
	if e.Threads <= 0 {
		e.Threads = 1
	}

	//Start the amount of processes related to the amount of given threads:
	for i := 0; i < e.Threads; i++ {
		e.Process = newProcess(e.Knowledge, e.Scanner, e.Pool)
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
				listener <- makeResult(r)
				e.WaitGroup.Done()
			}
		}
	}()

	//Listen a stop signal then wait until all background processes are completed:
	if <-e.quit {
		e.WaitGroup.Wait()
		fmt.Println(":: Scanner handler stopped")
		return
	}
}

// Add new jobs (tasks) to be performed by the handler processes:
func (e *Handler) AddJob(verifyMode bool, httpResult request.Result, knowledge knowledge.Knowledge) {
	e.WaitGroup.Add(1)
	e.JobQueue <- Job{
		Http:      httpResult,
		Knowledge: knowledge,
	}
}

func (e *Handler) Stop() {
	e.quit <- true
}

// Create a new process
func newProcess(knowl map[string]knowledge.Knowledge, scanner *config.Scanner, pool chan chan Job) process {
	return process{
		Knowledge: knowl,
		pool:      pool,
		Scanner:   scanner,
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

		//Confirm that we do have knowledge about our current target:
		targetKnowledge, ok_knowledge = p.Knowledge[job.Http.TargetHashId]

		//Behavior contains the methods that check unknown behavior along with the behavioral status of the current job:
		behavior = NewBehavior()
	)

	//Quick basic behavior checks:
	if ok_knowledge {
		behavior.status = behavior.QuickDetect(job.Http.Response, targetKnowledge)
	}

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
					//p.Scanner.Extract.Known,
				)
				extResult = ext.Run()

				//In case we do have knowledge of the target.
				if ok_knowledge {
					//Extract all the new unique regex and patterns discovered.
					//Note: Current map result MUST be first. Return the same order of "currentMaps" as the result.
					currentMaps := []map[string]int{
						extResult.RegexBody,
						extResult.RegexHeaders,
						extResult.PatternBody,
						extResult.PatternHeaders,
					}
					knownMaps := []map[string][]int{
						targetKnowledge.Combine.Extract.RegexBody,
						targetKnowledge.Combine.Extract.RegexHeaders,
						targetKnowledge.Combine.Extract.PatternBody,
						targetKnowledge.Combine.Extract.PatternHeaders,
					}
					ExtractMapDiff, totalhits := extract.GetMultiUnique(currentMaps, knownMaps, job.Http.Payload)

					//This shall not be happening, then it's a bug (critical)
					if len(ExtractMapDiff) != len(currentMaps) {
						log.Fatal(design.STATUS.CRITICAL,
							" The current maps used in the handler - extract process containing the list of maps did not match the diff list. Please report this to the official Firefly Github repository",
						)
					}
					//Note : (The same order as in "compareMaps")
					extResult = extract.Result{
						OK:             (totalhits > 0),
						TotalHits:      totalhits,
						RegexBody:      ExtractMapDiff[0],
						RegexHeaders:   ExtractMapDiff[1],
						PatternBody:    ExtractMapDiff[2],
						PatternHeaders: ExtractMapDiff[3],
					}
				}
				wg.Done()
			}()
		}

		//[Diff]erence process:
		if ok_knowledge && p.Scanner.OK_Diff /* && job.Knowledge != nil */ {
			wg.Add(1)
			go func() {
				//Make a new difference instant and provided the current HTTP response body and headers:
				diff := difference.NewDifference(
					difference.Properties{
						Payload:          job.Http.Payload,
						PayloadVerify:    targetKnowledge.PayloadVerify,
						CompareHTMLNodes: targetKnowledge.Combine.HTMLNode,
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
				diffResult = diff.Run()
				wg.Done()
			}()
		}
		//Wait for all the scanners to finish:
		wg.Wait()
	}

	//Collect all the result for each scan:

	//Confirm the unexpected behavior
	if (ok_knowledge && !behavior.status) && (diffResult.OK || extResult.OK || transformationResult.OK) {
		behavior.status = true
	}

	return processResult{
		UnkownBehavior: behavior.status,
		Http:           job.Http,
		Extract:        extResult,
		Difference:     diffResult,
		Transformation: transformationResult,
	}
}

// Start the extract scanning process
func makeResult(pResult processResult) Result {
	req := pResult.Http.Request
	resp := pResult.Http.Response

	return Result{
		Output: output.ResultFinal{
			TargetHashId:   pResult.Http.TargetHashId,
			RequestId:      pResult.Http.RequestId,
			Tag:            pResult.Http.Tag,
			Date:           pResult.Http.Date,
			Payload:        pResult.Http.Payload,
			UnkownBehavior: pResult.UnkownBehavior,
			OK:             true,

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
				ContentLength: resp.ContentLength,
				HeaderAmount:  resp.HeaderAmount,
				Headers:       resp.Header,
			},
			Scanner: output.Scanner{
				Extract:        pResult.Extract,
				Diff:           pResult.Difference,
				Transformation: pResult.Transformation,
				//Data...
			},

			Error: nil,
		},
		Error: nil,
	}
}
