package runner

import (
	"fmt"

	"github.com/Brum3ns/firefly/internal/scan"
	"github.com/google/uuid"
)

type jobScanner struct {
	id         string
	job        jobHTTP
	scanResult scan.Result
}

func newJobScanner(job jobHTTP) jobScanner {
	return jobScanner{
		id:  fmt.Sprintf("scan-%s", uuid.NewString()),
		job: job,
	}
}

func (j jobScanner) setScanResult(result scan.Result) {
	j.scanResult = result
}
