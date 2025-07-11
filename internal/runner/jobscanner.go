package runner

import (
	"fmt"

	"github.com/Brum3ns/firefly/internal/scan"
	"github.com/google/uuid"
)

type jobScanner struct {
	id         string
	coreJob    jobHTTP
	scanResult scan.Result
}

func newJobScanner(job jobHTTP) jobScanner {
	return jobScanner{
		id:      fmt.Sprintf("scan-%s", uuid.NewString()),
		coreJob: job,
	}
}

func (j *jobScanner) setScanResult(result scan.Result) {
	j.scanResult = result
}
