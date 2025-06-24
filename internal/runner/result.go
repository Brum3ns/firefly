package runner

import (
	"context"
)

type result struct{}

func (r *Runner) handlerResult(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case result := <-r.channel.result:
			r.wg.scanner.Done()
			// Analyze result
			//Code ...

			// New job / Finish (nothing more todo)
			//Code ...

			if result.id == "" {

			}

			// Fix output to file
			// Code...

			r.wg.result.Done()

		case <-ctx.Done():
			return
		}
	}
}
