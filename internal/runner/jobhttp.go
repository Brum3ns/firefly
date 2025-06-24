package runner

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Brum3ns/firefly/pkg/rhttp"
	"github.com/google/uuid"
)

type jobHTTP struct {
	id             string
	payload        string
	target         string
	httpRawRequest rhttp.HttpRaw
	httpResponse   rhttp.Response
}

func newJobHTTP(payload string) jobHTTP {
	return jobHTTP{
		id:      fmt.Sprintf("http-%s", uuid.NewString()),
		payload: payload,
	}
}

func (job *jobHTTP) setHttpRawRequest(httpRaw rhttp.HttpRaw) {
	job.httpRawRequest = httpRaw
}

func (job *jobHTTP) setResponse(resp *http.Response, respTime time.Duration) error {
	r, err := rhttp.NewResponse(resp, respTime)
	if err != nil {
		return err
	}
	job.httpResponse = r
	return nil
}

func (r *Runner) sendCoreJobs() {
	var payloads []string
	// Select payloads based on runner mode
	switch r.mode {
	case mode_fuzz:
		payloads = r.payload.GetWordlist()
	case mode_knowledge:
		// Add the verify payload the amount of time to verify the target knowledge
		for i := 0; i < r.option.VerifyKnowledge; i++ {
			payloads = append(payloads, r.option.PayloadVerify)
		}
	}
	// Send the jobs
	for _, httpReq := range r.request.GetCoreRawRequests() {
		for _, pyld := range payloads {
			// Make HTTP request and insert payload
			job := newJobHTTP(pyld)
			job.setHttpRawRequest(httpReq)

			r.wg.httpRequest.Add(1)
			r.channel.httpRequest <- job
		}
	}
}
