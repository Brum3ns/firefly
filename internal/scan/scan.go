package scan

import (
	"github.com/Brum3ns/firefly/internal/knowledge"
	"github.com/Brum3ns/firefly/pkg/httpdiff"
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/randomness"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

type Scan struct {
	config Config
}

type Config struct {
	HTTPResponse   rhttp.Response
	Randomness     randomness.Randomness
	Knowledge      knowledge.MergedKnowledge
	HTTPDiffFilter httpdiff.Filter
	Payload        string
	VerifyPayload  string
}
type Result struct {
	HTTPDiff httpdiff.Result
}

func NewScan(config Config) (Scan, error) {
	return Scan{
		config: config,
	}, nil
}

func (scan *Scan) Run() (Result, error) {
	resultHttpDiff := scan.Diff()

	return Result{
		HTTPDiff: resultHttpDiff,
	}, nil
}

func (scan *Scan) Diff() httpdiff.Result {
	//Make a new difference instant and provided the current HTTP response body and headers:
	diff := httpdiff.NewDifference(
		httpdiff.Config{
			Payload:       scan.config.Payload,
			PayloadVerify: scan.config.VerifyPayload,
			Compare: httpdiff.Compare{
				HTMLMergeNode:   scan.config.Knowledge.Merge.HTMLNode,
				HeaderMergeNode: scan.config.Knowledge.Merge.HeaderNode,
			},
			Randomness: scan.config.Randomness,
			Filter:     scan.config.HTTPDiffFilter,
		},
	)

	headerResult := diff.GetHeadersDiff(httpnode.GetHeaderNode(scan.config.HTTPResponse.Headers))
	htmlResult := diff.GetHTMLNodeDiff(httpnode.GetHTMLNode(scan.config.HTTPResponse.Body))

	return httpdiff.Result{
		OK:           (headerResult.OK || htmlResult.OK),
		HeaderResult: headerResult,
		HTMLResult:   htmlResult,
	}
}
