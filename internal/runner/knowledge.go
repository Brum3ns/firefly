package runner

import (
	"context"
	"crypto/sha1"
	"fmt"
)

func (r *Runner) handleKnowledge(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case jobScanner := <-r.channel.knowledge:
			r.wg.scanner.Done()

			method := jobScanner.job.httpRawRequest.Method
			url := jobScanner.job.httpRawRequest.Url

			r.knowledge.AppendKnowledge(
				makeTargetHash(method, url),
				jobScanner.job.payload,
				jobScanner.job.httpResponse,
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
