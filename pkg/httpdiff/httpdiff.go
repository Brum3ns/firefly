package httpdiff

import (
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/randomness"
)

type HttpDiff struct {
	config Config
}

type Config struct {
	Payload       string
	PayloadVerify string
	Randomness    randomness.Randomness
	Filter        Filter
	Merge         Merge
}

type Merge struct {
	HeaderMergeNode httpnode.HeaderMergeNode
	HTMLMergeNode   httpnode.HTMLNodeMerge
	JSONMergeNode   httpnode.JSONNodeMerge
}

type Result struct {
	OK     bool         `json:"ok"`
	Header HeaderResult `json:"headerresult"`
	HTML   HTMLResult   `json:"htmlresult"`
	JSON   JSONResult   `json:"jsonresult"`
}

type diffNode struct {
	hit             int
	checkRandomness bool
	data            map[string]int
}

func NewHttpDiff(config Config) *HttpDiff {
	return &HttpDiff{
		config: config,
	}
}

func newDiffNode() diffNode {
	return diffNode{
		data: make(map[string]int),
	}
}

// TODO
/* func (hdiff *HttpDiff) GetJSONDiff(jsonNode httpnode.json) JsonResult {
	return JsonResult{}
} */

func (hdiff *HttpDiff) nodeDiff(current diffNode, known map[string][]int, payload string) (diffNode, diffNode) {
	var (
		appear      = newDiffNode()
		disappear   = newDiffNode()
		testedItems = make(map[string]struct{})
	)

	for currentItem, currentValue := range current.data {
		// Set the isDiff as true by default
		isDiff := true
		amountDiff := 0

		// Check if the current item exists in the known
		if knownValues, ok := known[currentItem]; ok {

			// Add the item to the tested map to be used to discover items that disappear in the response body
			testedItems[currentItem] = struct{}{}

			// Compare with the known values
			for _, knownValue := range knownValues {
				if currentValue == knownValue {
					isDiff = false
					break

					// If the current value isen't in the list check the amount diff and update to the highest amount diff
				} else if v := lengthMinMaxDiff(currentValue, knownValue); v > amountDiff {
					amountDiff = v
				}
			}
		} else if amountDiff == 0 {
			amountDiff = currentValue
		}

		if isDiff && currentItem != payload {
			// Check randomness (false positive)
			if !current.checkRandomness || (current.checkRandomness && !hdiff.config.Randomness.IsRandom(currentItem)) {
				appear.data[currentItem] = amountDiff
				appear.hit += amountDiff
			}
		}
	}

	// Check known item and see if any of them where not included in the current response, then add them as a valid diff
	for knownItem, knownValues := range known {
		if _, ok := testedItems[knownItem]; !ok {
			// Check randomness (false positive)
			if !current.checkRandomness || (current.checkRandomness && !hdiff.config.Randomness.IsRandom(knownItem)) {
				value := highestLstIntValue(knownValues)
				disappear.data[knownItem] = value
				disappear.hit += value
			}
		}
	}
	return appear, disappear
}
