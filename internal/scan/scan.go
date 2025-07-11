package scan

import (
	"github.com/Brum3ns/firefly/internal/knowledge"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpdiff"
	"github.com/Brum3ns/firefly/pkg/randomness"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

// Scan represents a scanning operation configuration and behavior.
type Scan struct {
	config Config
}

// Config holds the configuration for a Scan.
type Config struct {
	HTTPResponse   rhttp.Response
	Randomness     randomness.Randomness
	Knowledge      knowledge.MergedKnowledge
	HTTPDiffFilter httpdiff.Filter
	Payload        string
	VerifyPayload  string
	Extract        extract.Extract
	Skip           Skip
}

// Skip specifies which checks to skip during the scan.
type Skip struct {
	SkipBodyDiff      bool
	SkipHeaderDiff    bool
	SkipExtractBody   bool
	SkipExtractHeader bool
}

// Result holds the outcome of the scanning process.
type Result struct {
	OK       bool            `json:"ok"`
	HTTPDiff httpdiff.Result `json:"httpdiff"`
	Extract  Extract         `json:"extract"`
}

// Extract contains the results of the extraction operation.
type Extract struct {
	OK     bool           `json:"ok"`
	Body   extract.Result `json:"body"`
	Header extract.Result `json:"header"`
}

// NewScan creates a new Scan instance with the provided configuration.
func NewScan(config Config) (Scan, error) {
	return Scan{
		config: config,
	}, nil
}

// Run initiates the scanning process and returns the overall result.
func (scan *Scan) Run() (Result, error) {
	resultDiff := scan.Diff()
	resultExtract := scan.Extract()
	return Result{
		OK:       resultDiff.OK || resultExtract.OK,
		HTTPDiff: resultDiff,
		Extract:  resultExtract,
	}, nil
}
