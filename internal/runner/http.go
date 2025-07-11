package runner

import (
	"context"
	"log"
	"time"

	"github.com/Brum3ns/firefly/pkg/faker"
	"github.com/Brum3ns/firefly/pkg/httpfilter"
	"github.com/Brum3ns/firefly/pkg/payload"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

func (r *Runner) handlerJobRequest(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case job := <-r.channel.httpRequest:
			r.workerpool.request.Submit(func() {
				// Insert faker data
				rawHttpRequest := string(job.httpRawRequest.TemplateRawRequest)
				if !r.option.Placeholder.Disable {
					var err error
					rawHttpRequest, err = faker.Generate(rawHttpRequest)
					if err != nil {
						log.Println(err)
						r.wg.httpRequest.Done()
						return
					}
				}

				// Payload data insert
				rawRequest := payload.Insert(
					rawHttpRequest,
					r.option.Payload.Placeholder,
					job.payload,
				)

				timer := time.Now()

				// Jitter delay if in knowledge mode
				if r.mode == mode_knowledge && r.option.Knowledge.Jitter > 0 {
					time.Sleep(rhttp.GetJitter(r.option.Knowledge.Jitter))
				}

				resp, err := r.request.Client.SendRawRequest(r.request.GetURL(), []byte(rawRequest))
				if err != nil {
					log.Printf("HTTP request failed, error : %v", err)
					r.wg.httpRequest.Done()
					return
				}
				respTime := time.Since(timer)
				if err := job.setResponse(resp, respTime); err != nil {
					log.Printf("could not make response instance for job, id:[%v], error %v", job.id, err)
					r.wg.httpRequest.Done()
					return
				}

				// Forward job to workers
				// Code...
				r.wg.httpResponse.Add(1)
				r.channel.httpResponse <- job
				// HTTP request delay
				time.Sleep(time.Duration(r.option.Http.Delay) * time.Millisecond)
			})
		case <-ctx.Done():
			return
		}
	}
}

func (r *Runner) handlerResponse(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case job := <-r.channel.httpResponse:
			r.wg.httpRequest.Done()
			//Filter the HTTP response (if set):
			filterResp := httpfilter.Response{
				Body:         []byte(job.httpResponse.Body),
				StatusCode:   job.httpResponse.StatusCode,
				ResponseSize: job.httpResponse.LineCount,
				WordCount:    job.httpResponse.WordCount,
				LineCount:    job.httpResponse.LineCount,
				ResponseTime: job.httpResponse.Time,
				Headers:      job.httpResponse.Headers,
			}

			// HTTP Filter filter/match (if set)
			if r.httpfilter.Run(filterResp) || (r.httpmatch.IsSet() && !r.httpmatch.Run(filterResp)) {
				r.wg.httpResponse.Done()
				continue
			}
			// Check filters
			// Code...

			// Make scanner job and send to scanner
			// Code...
			r.wg.scanner.Add(1)
			r.channel.scanner <- newJobScanner(job)
		case <-ctx.Done():
			return
		}
	}
}
