package scan

import (
	"github.com/Brum3ns/firefly/pkg/httpdiff"
	"github.com/Brum3ns/firefly/pkg/httpnode"
)

// Diff compares the HTTP response against the expectations and returns the differences.
func (scan *Scan) Diff() httpdiff.Result {
	diff := httpdiff.NewHttpDiff(
		httpdiff.Config{
			Payload:       scan.config.Payload,
			PayloadVerify: scan.config.VerifyPayload,
			Merge: httpdiff.Merge{
				HTMLMergeNode:   scan.config.Knowledge.Merge.HTMLNode,
				HeaderMergeNode: scan.config.Knowledge.Merge.HeaderNode,
				JSONMergeNode:   scan.config.Knowledge.Merge.JSONNode,
			},
			Randomness: scan.config.Randomness,
			Filter:     scan.config.HTTPDiffFilter,
		},
	)

	var (
		headerResult httpdiff.HeaderResult
		htmlResult   httpdiff.HTMLResult
		jsonResult   httpdiff.JSONResult
	)

	if !scan.skipBodyDiff() {
		switch httpnode.DetectResponseType(scan.config.HTTPResponse.Body) {
		case httpnode.ResponseTypeJSON:
			if r, err := httpnode.GetJSONNode(scan.config.HTTPResponse.Body); err != nil {
				jsonResult.Malformed = true
			} else {
				jsonResult = diff.GetJSONNodeDiff(r)
			}
		default:
			htmlResult = diff.GetHTMLNodeDiff(httpnode.GetHTMLNode(scan.config.HTTPResponse.Body))
		}
	}

	if !scan.skipHeaderDiff() {
		headerResult = diff.GetHeadersDiff(httpnode.GetHeaderNode(scan.config.HTTPResponse.Headers))
	}

	return httpdiff.Result{
		OK:     (headerResult.OK || htmlResult.OK),
		Header: headerResult,
		HTML:   htmlResult,
		JSON:   jsonResult,
	}
}

// skipBodyDiff determines whether the body diff check should be skipped.
func (scan *Scan) skipBodyDiff() bool {
	return scan.config.Skip.SkipBodyDiff
}

// skipHeaderDiff determines whether the header diff check should be skipped.
func (scan *Scan) skipHeaderDiff() bool {
	return scan.config.Skip.SkipHeaderDiff
}
