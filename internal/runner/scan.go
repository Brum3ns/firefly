package runner

import (
	"context"
	"log"

	"github.com/Brum3ns/firefly/internal/scan"
	"github.com/Brum3ns/firefly/pkg/httpdiff"
)

func (r *Runner) handleJobScanner(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {

		case jobScanner := <-r.channel.scanner:
			r.wg.httpResponse.Done()
			// Make adapted struct for scanner
			// Code...

			// Make target hash
			targetHash := makeTargetHash(
				jobScanner.coreJob.httpRawRequest.Method,
				jobScanner.coreJob.httpRawRequest.Url,
			)

			// Setup the scanner engine for the current job
			scanner, err := scan.NewScan(scan.Config{
				Payload:       jobScanner.coreJob.payload,
				VerifyPayload: r.option.Payload.VerifyCanary,
				HTTPResponse:  jobScanner.coreJob.httpResponse,
				Randomness:    r.randomness,
				Extract:       r.extract,
				Skip: scan.Skip{
					SkipBodyDiff:      r.option.Skip.DiffBody,
					SkipHeaderDiff:    r.option.Skip.DiffHeader,
					SkipExtractBody:   r.option.Skip.ExtractBody,
					SkipExtractHeader: r.option.Skip.ExtractHeader,
				},
				HTTPDiffFilter: httpdiff.Filter{
					Headers: r.option.Filter.FilterDiffHeaders,
				}, // TODO
				Knowledge: r.knowledge.Merged[targetHash],
			})
			if err != nil {
				log.Println(err)
				continue
			}

			// Run the scanner
			result, err := scanner.Run()
			if err != nil {
				log.Printf("scan failed, error %v", err)
				r.wg.scanner.Done()
				continue
			}
			jobScanner.setScanResult(result) // TODO : Fix knowledge handling first

			switch r.getMode() {
			case mode_fuzz:
				r.wg.result.Add(1)
				r.channel.result <- jobScanner

			case mode_knowledge:
				r.wg.knowledge.Add(1)
				r.channel.knowledge <- jobScanner
			}

		case <-ctx.Done():
			return
		}
	}
}
