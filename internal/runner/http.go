package runner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Brum3ns/firefly/pkg/faker"
	"github.com/Brum3ns/firefly/pkg/httpfilter"
	"github.com/Brum3ns/firefly/pkg/payload"
)

func (r *Runner) handlerJobRequest(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case job := <-r.channel.httpRequest:
			r.workerpool.request.Submit(func() {
				// Insert faker data
				rawHttpRequest := string(job.httpRawRequest.TemplateRawRequest)
				if !r.option.DisablePlaceholders {
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
					r.option.FuzzPlaceholder,
					job.payload,
				)

				timer := time.Now()
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
				time.Sleep(time.Duration(r.option.HttpDelay) * time.Millisecond)
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

			fmt.Println(r.wg.httpRequest.GetCount(), job.httpResponse.Time)

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
