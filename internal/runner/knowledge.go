package runner

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/Brum3ns/firefly/internal/knowledge"
)

func (r *Runner) handleKnowledge(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case jobScanner := <-r.channel.knowledge:
			r.wg.scanner.Done()

			method := jobScanner.coreJob.httpRawRequest.Method
			url := jobScanner.coreJob.httpRawRequest.Url

			r.knowledge.AppendKnowledge(
				knowledge.KnowledgeMeta{
					TargetHash:          makeTargetHash(method, url),
					Payload:             jobScanner.coreJob.payload,
					HTTPResponse:        jobScanner.coreJob.httpResponse,
					ExtractResultBody:   jobScanner.scanResult.Extract.Body,
					ExtractResultHeader: jobScanner.scanResult.Extract.Header,
				},
			)
			r.wg.knowledge.Done()

		case <-ctx.Done():
			return
		}
	}
}

func makeTargetHash(method, url string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(method+url)))
}

func (r *Runner) hasKnowledge() bool {
	return len(r.knowledge.Merged) > 0
}
