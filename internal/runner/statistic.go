package runner

import (
	"context"
)

func (r *Runner) handlerStatistic(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case result := <-r.channel.statistic:
			r.wg.result.Done()

			if !r.option.Process.Silence {
				PrintVerbose(Verbose{
					time:   getTime(),
					result: result,
				})
			}

			r.wg.verbose.Done()

		case <-ctx.Done():
			return
		}
	}
}
