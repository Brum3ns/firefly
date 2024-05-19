package scan

import (
	"log"

	"github.com/Brum3ns/firefly/internal/config"
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpdiff"
	"github.com/Brum3ns/firefly/pkg/httpprepare"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/transformation"
)

// scan represents the scan that executes the job
type scan struct {
	jobChannel chan Job
	pool       chan chan Job

	Scanner *config.Scanner //!Note : (Static data stored. Read struct DESC)
	Result  scanResult
}

type scanResult struct {
	UnkownBehavior bool
	Http           request.Result
	Extract        extract.Result
	Difference     httpdiff.Result
	Transformation transformation.Result
}

// Create a new scan
func newScan(scanner *config.Scanner, pool chan chan Job) scan {
	return scan{
		pool:       pool,
		Scanner:    scanner,
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
		//Behavior contains the methods that check unknown behavior along with the behavioral status of the current job:
		behavior = NewBehavior()

		// Scanning techniques
		ResultExtract        extract.Result
		ResultDifference     httpdiff.Result
		ResultTransformation transformation.Result
	)

	//Quick basic behavior checks:
	if job.OK_knowledge {
		behavior.status = behavior.QuickDetect(job)
	}

	//Check if we should preform scanner techniques or not:
	if !s.Scanner.DisablesTechniques {

		if s.Scanner.OK_Diff {
			ResultDifference = s.Difference(job)
		}
		if s.Scanner.OK_Extract {
			ResultExtract = s.Extract(job)
		}
		if s.Scanner.OK_Transformation {
			ResultTransformation = s.Transformation(job)
		}
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

// Scan for errors, patterns in the response that have been triggered by the payload
func (s scan) Extract(job Job) extract.Result {
	e := s.Scanner.Extract
	e.AddJob(
		job.Http.Response.Body,
		job.Http.Response.HeaderString,
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

// Scan for differences in the current compare to the known HTTP responses
func (s scan) Difference(job Job) httpdiff.Result {
	//Make a new difference instant and provided the current HTTP response body and headers:
	diff := httpdiff.NewDifference(
		httpdiff.Config{
			Payload:       job.Http.Payload,
			PayloadVerify: job.Knowledge.PayloadVerify,
			Compare: httpdiff.Compare{
				HTMLMergeNode:   job.Knowledge.Combine.HTMLNode,
				HeaderMergeNode: job.Knowledge.Combine.HeaderNode,
			},
			Randomness: s.Scanner.Randomness,
		},
	)

	headerResult := diff.GetHeadersDiff(httpprepare.GetHeaderNode(job.Http.Response.Header))
	htmlResult := diff.GetHTMLNodeDiff(httpprepare.GetHTMLNode(job.Http.Response.Body))

	return httpdiff.Result{
		OK:           (headerResult.OK || htmlResult.OK),
		HeaderResult: headerResult,
		HTMLResult:   htmlResult,
	}
}

// Scan for transformations within the payload
func (s scan) Transformation(job Job) transformation.Result {
	tfmt := s.Scanner.Transformation
	return tfmt.Detect(job.Http.Response.Body, job.Http.Payload)
}
