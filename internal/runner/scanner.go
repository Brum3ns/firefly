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

		case job := <-r.channel.scanner:
			r.wg.httpResponse.Done()
			// Make adapted struct for scanner
			// Code...

			// Forward job to the scanner engine
			// Code...

			scanner, err := scan.NewScan(
				scan.Config{
					HTTPResponse:   job.job.httpResponse,
					HTTPDiffFilter: httpdiff.Filter{}, // TODO : keep or not?
					Randomness:     r.randomness,
				},
			)
			if err != nil {
				log.Printf("could not create scanner instance, error %v", err)
			}

			result, err := scanner.Run()
			if err != nil {
				log.Printf("scan failed, error %v", err)
				r.wg.scanner.Done()
				continue
			}
			job.setScanResult(result) // TODO : Fix knowledge handling first

			switch r.getMode() {
			case mode_fuzz:
				r.wg.result.Add(1)
				r.channel.result <- job

			case mode_knowledge:
				r.wg.knowledge.Add(1)
				r.channel.knowledge <- job
			}

		case <-ctx.Done():
			return
		}
	}
}
