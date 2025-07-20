package httpdiff

import (
	"slices"

	"github.com/Brum3ns/firefly/pkg/httpnode"
)

type HeaderResult struct {
	OK         bool                     `json:"ok"`
	HeaderHits int                      `json:"headerhits"`
	Appear     httpnode.HeaderMergeNode `json:"appear"`
	Disappear  httpnode.HeaderMergeNode `json:"disappear"`
}

// Filter HTTP header name
func (hdiff *HttpDiff) FilterHeader(header string /*value string*/) bool {
	if len(hdiff.config.Filter.Headers) == 0 {
		return false
	}
	return slices.Contains(hdiff.config.Filter.Headers, header)
}

// Take two prepared header structures and compare their differences
func (hdiff *HttpDiff) GetHeadersDiff(HeaderNode httpnode.HeaderMergeNode) HeaderResult {
	var (
		appear      = httpnode.NewHeader()
		disappear   = httpnode.NewHeader()
		testedItems = make(map[string]struct{})
		totalHits   = 0
	)

	for currentHeader, currentHeaderData := range HeaderNode {
		// Diff Filter check
		if hdiff.FilterHeader(currentHeader) {
			continue
		}

		// Check if the information is unique in relation to the shared header names
		if knownHeaderData, ok := hdiff.config.Merge.HeaderMergeNode[currentHeader]; ok {

			// Add to tested header names
			testedItems[currentHeader] = struct{}{}

			// Search for unique values inside the current and known header data values
			if (!lstIntShareItem(knownHeaderData.Amount, currentHeaderData.Amount) ||
				!lstStringShareItem(knownHeaderData.Values, currentHeaderData.Values)) &&
				!slices.Contains(knownHeaderData.Values, hdiff.config.Payload) { // TODO : we want to detect the payload reflected in the headers, no?

				appear[currentHeader] = currentHeaderData
			}
		} else {
			appear[currentHeader] = currentHeaderData
		}
	}

	// Add headers that disappear in the current response compare to the original
	for knownHeader, knownHeaderData := range hdiff.config.Merge.HeaderMergeNode {
		// Diff Filter check
		if hdiff.FilterHeader(knownHeader) {
			continue
		}

		if _, ok := testedItems[knownHeader]; !ok {
			// Extract the highest difference from the known values and add it
			disappear[knownHeader] = httpnode.HeaderInfo{
				Amount: []int{highestLstIntValue(knownHeaderData.Amount)},
				Values: knownHeaderData.Values,
			}
		}

	}
	totalHits = len(disappear) + len(appear)
	return HeaderResult{
		OK:         (totalHits > 0),
		HeaderHits: totalHits,
		Appear:     appear,
		Disappear:  disappear,
	}
}
