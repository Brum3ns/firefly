package scan

import (
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

// Extract performs the extraction of data from the HTTP response.
func (scan *Scan) Extract() Extract {
	var (
		resultExtractBody   extract.Result
		resultExtractHeader extract.Result
	)
	// Check if body extraction should be performed
	if !scan.skipBodyExtract() {
		// Execute body extraction
		resultExtractBody = scan.config.Extract.Run(scan.config.HTTPResponse.Body)
	}
	// Check if header extraction should be performed
	if !scan.skipHeaderExtract() {
		// Execute header extraction after converting headers to string
		resultExtractHeader = scan.config.Extract.Run(rhttp.HeaderToString(scan.config.HTTPResponse.Headers))
	}

	// Get the unique hits

	// Return the combined extraction results
	return Extract{
		OK:     resultExtractBody.OK || resultExtractHeader.OK, // Indicate if either extraction was successful
		Body:   resultExtractBody,
		Header: resultExtractHeader,
	}
}

// skipBodyExtract determines whether the body extraction should be skipped.
func (scan *Scan) skipBodyExtract() bool {
	return scan.config.Skip.SkipExtractBody
}

// skipHeaderExtract determines whether the header extraction should be skipped.
func (scan *Scan) skipHeaderExtract() bool {
	return scan.config.Skip.SkipExtractHeader
}
