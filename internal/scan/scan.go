package scan

import (
	"log"

	"github.com/Brum3ns/firefly/internal/config"
	"github.com/Brum3ns/firefly/internal/knowledge"
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/difference"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/prepare"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/transformation"
)

// scan represents the scan that executes the job
type scan struct {
	jobChannel chan Job
	pool       chan chan Job
	Scanner    *config.Scanner //!Note : (Static data stored. Read struct DESC)
	Result     scanResult      //<-Returned

	//Knowledge  map[string][]verify.TargetKnowledge //Note : (Should be a pointer of "Properties.VerifyStorage")
	Knowledge map[string]knowledge.Knowledge
}

type scanResult struct {
	UnkownBehavior bool
	Http           request.Result
	Extract        extract.Result
	Difference     difference.Result
	Transformation transformation.Result
}

// Create a new scan
func newScan(knowl map[string]knowledge.Knowledge, scanner *config.Scanner, pool chan chan Job) scan {
	return scan{
		Knowledge: knowl,
		pool:      pool,
		Scanner:   scanner,
		//Knowledge:  verifiedStorage,
		jobChannel: make(chan Job),
	}
}

// Spawn a new scan process
func (s scan) spawnScan(result chan scanResult) {
	go func() {
		for {
			// Add the current spawned scan into the scanning queue:
			s.pool <- s.jobChannel

			//A job was given, start processing it
			select {
			case job := <-s.jobChannel:
				result <- s.scan(job)
			}
		}
	}()
}

// Start a new process
func (s scan) scan(job Job) scanResult {
	var (
		//wg sync.WaitGroup

		// Confirm that we do have knowledge about our current target:
		//targetKnowledge, ok_knowledge = s.Knowledge[job.Http.TargetHashId]

		//Behavior contains the methods that check unknown behavior along with the behavioral status of the current job:
		behavior = NewBehavior()

		ResultExtract        extract.Result
		ResultDifference     difference.Result
		ResultTransformation transformation.Result
	)

	//Quick basic behavior checks:
	/* if ok_knowledge {
		behavior.status = behavior.QuickDetect(job.Http.Response, targetKnowledge)
	} */

	//Check if we should preform scanner techniques or not:
	if !s.Scanner.DisablesTechniques {
		ResultDifference = s.Difference(&job)
		ResultExtract = s.Extract(&job)
		ResultTransformation = s.Transformation(&job)

		/* 	// Difference scan
		if s.Scanner.OK_Diff {
			wg.Add(1)
			go func() {
				ResultDifference = s.Difference(&job)
				wg.Done()
			}()
		}
		// Extract scan
		if s.Scanner.OK_Extract {
			wg.Add(1)
			go func() {
				ResultExtract = s.Extract(&job)
				wg.Done()
			}()
		}
		// Transformation scan
		if s.Scanner.OK_Transformation {
			wg.Add(1)
			go func() {
				ResultTransformation = s.Transformation(&job)
				wg.Done()
			}()
		}
		// Wait until all the scans are done
		wg.Wait() */
	}

	//Confirm the unexpected behavior
	if (job.OK_knowledge && !behavior.status) && (ResultDifference.OK || ResultExtract.OK || ResultTransformation.OK) {
		behavior.status = true
	}

	return scanResult{
		UnkownBehavior: behavior.status,
		Http:           job.Http,
		Extract:        ResultExtract,
		Difference:     ResultDifference,
		Transformation: ResultTransformation,
	}
}

func (s scan) Extract(job *Job) extract.Result {
	e := s.Scanner.Extract
	e.AddJob(
		job.Http.Response.Body,
		job.Http.Response.HeaderString,
		//p.Scanner.Extract.Known,
	)

	result := e.Run()

	//In case we do have knowledge of the target.
	if job.OK_knowledge {
		//Extract all the new unique regex and patterns discovered.
		//Note: Current map result MUST be first. Return the same order of "currentMaps" as the result.
		currentMaps := []map[string]int{
			result.RegexBody,
			result.RegexHeaders,
			result.PatternBody,
			result.PatternHeaders,
		}
		knownMaps := []map[string][]int{
			job.Knowledge.Combine.Extract.RegexBody,
			job.Knowledge.Combine.Extract.RegexHeaders,
			job.Knowledge.Combine.Extract.PatternBody,
			job.Knowledge.Combine.Extract.PatternHeaders,
		}
		ExtractMapDiff, totalhits := extract.GetMultiUnique(currentMaps, knownMaps, job.Http.Payload)

		//This shall not be happening, then it's a bug (critical)
		if len(ExtractMapDiff) != len(currentMaps) {
			log.Fatal(design.STATUS.CRITICAL,
				" The current maps used in the handler - extract process containing the list of maps did not match the diff list. Please report this to the official Firefly Github repository",
			)
		}
		//Note : (The same order as in "compareMaps")
		result = extract.Result{
			OK:             (totalhits > 0),
			TotalHits:      totalhits,
			RegexBody:      ExtractMapDiff[0],
			RegexHeaders:   ExtractMapDiff[1],
			PatternBody:    ExtractMapDiff[2],
			PatternHeaders: ExtractMapDiff[3],
		}
	}
	return result
}

func (s scan) Difference(job *Job) difference.Result {
	//Make a new difference instant and provided the current HTTP response body and headers:
	diff := difference.NewDifference(
		difference.Properties{
			Payload:          job.Http.Payload,
			PayloadVerify:    job.Knowledge.PayloadVerify,
			CompareHTMLNodes: job.Knowledge.Combine.HTMLNode,
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
	return diff.Run()
}

func (s scan) Transformation(job *Job) transformation.Result {
	tfmt := s.Scanner.Transformation
	return tfmt.Detect(job.Http.Response.Body, job.Http.Payload)
}
