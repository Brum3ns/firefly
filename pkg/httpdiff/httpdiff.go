package httpdiff

import (
	"slices"

	"github.com/Brum3ns/firefly/pkg/httpprepare"
	"github.com/Brum3ns/firefly/pkg/randomness"
)

type Difference struct {
	*Config
}

type Config struct {
	Payload       string
	PayloadVerify string
	Randomness    randomness.Randomness
	Filter
	Compare
}

type Compare struct {
	HeaderMergeNode httpprepare.Header
	HTMLMergeNode   httpprepare.HTMLNodeCombine
}

type Result struct {
	OK bool
	HeaderResult
	HTMLResult
}

type HeaderResult struct {
	OK         bool
	HeaderHits int
	Appear     httpprepare.Header
	Disappear  httpprepare.Header
}

type HTMLResult struct {
	OK        bool
	Appear    HTMLNodeDiff
	Disappear HTMLNodeDiff
}

type HTMLNodeDiff struct {
	TagStartHits       int
	TagEndHits         int
	TagSelfCloseHits   int
	WordsHits          int
	CommentHits        int
	AttributeHits      int
	AttributeValueHits int
	httpprepare.HTMLNode
}
type Filter struct {
	HeaderFilter
	//HTMLFilter
}

type HeaderFilter struct {
	Header httpprepare.Header
}

type diffNode struct {
	hit             int
	checkRandomness bool
	data            map[string]int
}

func NewDifference(config Config) *Difference {
	return &Difference{
		Config: &config,
	}
}

func newDiffNode() diffNode {
	return diffNode{
		data: make(map[string]int),
	}
}

// Run the [diff]erence enumiration process for the HTML node
func (diff *Difference) GetHTMLNodeDiff(htmlNode httpprepare.HTMLNode) HTMLResult {
	totalHits := 0
	storage := struct {
		appearHits    int
		disappearHits int
		appear        []diffNode
		disappear     []diffNode
	}{}

	//!Note : Order for "current" and "known" togther with the list length *MUST* be the same:
	current := [7]diffNode{
		{data: htmlNode.TagStart},
		{data: htmlNode.TagEnd},
		{data: htmlNode.TagSelfClose},
		{data: htmlNode.Words, checkRandomness: true},
		{data: htmlNode.Comment, checkRandomness: true},
		{data: htmlNode.Attribute},
		{data: htmlNode.AttributeValue, checkRandomness: true},
	}
	known := [7]map[string][]int{
		diff.Config.Compare.HTMLMergeNode.TagStart,
		diff.Config.Compare.HTMLMergeNode.TagEnd,
		diff.Config.Compare.HTMLMergeNode.TagSelfClose,
		diff.Config.Compare.HTMLMergeNode.Words,
		diff.Config.Compare.HTMLMergeNode.Comment,
		diff.Config.Compare.HTMLMergeNode.Attribute,
		diff.Config.Compare.HTMLMergeNode.AttributeValue,
	}

	for i := 0; i < len(current); i++ {
		//Detect difference
		diffAppear, diffDisappear := diff.nodeDiff(current[i], known[i], diff.Payload)

		storage.appear = append(storage.appear, diffAppear)
		storage.appearHits += diffAppear.hit

		storage.disappear = append(storage.disappear, diffDisappear)
		storage.disappearHits += diffDisappear.hit

		totalHits += (diffAppear.hit + diffDisappear.hit)
	}

	return HTMLResult{
		// !Note : The order for this *MUST* follow the same as "known" and "current" above
		OK: (totalHits > 0),
		Appear: HTMLNodeDiff{
			TagStartHits:       storage.appear[0].hit,
			TagEndHits:         storage.appear[1].hit,
			TagSelfCloseHits:   storage.appear[2].hit,
			WordsHits:          storage.appear[3].hit,
			CommentHits:        storage.appear[4].hit,
			AttributeHits:      storage.appear[5].hit,
			AttributeValueHits: storage.appear[6].hit,
			HTMLNode: httpprepare.HTMLNode{
				TagStart:       storage.appear[0].data,
				TagEnd:         storage.appear[1].data,
				TagSelfClose:   storage.appear[2].data,
				Words:          storage.appear[3].data,
				Comment:        storage.appear[4].data,
				Attribute:      storage.appear[5].data,
				AttributeValue: storage.appear[6].data,
			},
		},
		Disappear: HTMLNodeDiff{
			TagStartHits:       storage.disappear[0].hit,
			TagEndHits:         storage.disappear[1].hit,
			TagSelfCloseHits:   storage.disappear[2].hit,
			WordsHits:          storage.disappear[3].hit,
			CommentHits:        storage.disappear[4].hit,
			AttributeHits:      storage.disappear[5].hit,
			AttributeValueHits: storage.disappear[6].hit,
			HTMLNode: httpprepare.HTMLNode{
				TagStart:       storage.disappear[0].data,
				TagEnd:         storage.disappear[1].data,
				TagSelfClose:   storage.disappear[2].data,
				Words:          storage.disappear[3].data,
				Comment:        storage.disappear[4].data,
				Attribute:      storage.disappear[5].data,
				AttributeValue: storage.disappear[6].data,
			},
		},
	}
}

// Take two prepared header structures and compare their differences
func (diff *Difference) GetHeadersDiff(HeaderNode httpprepare.Header) HeaderResult {
	var (
		appear      = httpprepare.NewHeader()
		disappear   = httpprepare.NewHeader()
		testedItems = make(map[string]struct{})
		totalHits   = 0
	)

	for currentHeader, currentHeaderData := range HeaderNode {
		// Diff Filter check
		if diff.FilterHeader(currentHeader) {
			continue
		}

		// Check if the information is unique in relation to the shared header names
		if knownHeaderData, ok := diff.Compare.HeaderMergeNode[currentHeader]; ok {

			// Add to tested header names
			testedItems[currentHeader] = struct{}{}

			// Search for unique values inside the current and known header data values
			if (!lstIntShareItem(knownHeaderData.Amount, currentHeaderData.Amount) ||
				!lstStringShareItem(knownHeaderData.Values, currentHeaderData.Values)) &&
				!slices.Contains(knownHeaderData.Values, diff.Payload) {

				appear[currentHeader] = currentHeaderData
			}
		} else {
			appear[currentHeader] = currentHeaderData
		}
	}

	// Add headers that disappear in the current response compare to the original
	for knownHeader, knownHeaderData := range diff.Compare.HeaderMergeNode {
		// Diff Filter check
		if diff.FilterHeader(knownHeader) {
			continue
		}

		if _, ok := testedItems[knownHeader]; !ok {
			// Extract the highest difference from the known values and add it
			disappear[knownHeader] = httpprepare.HeaderInfo{
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

func (diff *Difference) nodeDiff(current diffNode, known map[string][]int, payload string) (diffNode, diffNode) {
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
			if !current.checkRandomness || (current.checkRandomness && !diff.Config.Randomness.IsRandom(currentItem)) {
				appear.data[currentItem] = amountDiff
				appear.hit += amountDiff
			}
		}
	}

	// Check known item and see if any of them where not included in the current response, then add them as a valid diff
	for knownItem, knownValues := range known {
		if _, ok := testedItems[knownItem]; !ok {
			// Check randomness (false positive)
			if !current.checkRandomness || (current.checkRandomness && !diff.Config.Randomness.IsRandom(knownItem)) {
				value := highestLstIntValue(knownValues)
				disappear.data[knownItem] = value
				disappear.hit += value
			}
		}
	}
	return appear, disappear
}

// Filter HTTP header name
func (diff *Difference) FilterHeader(header string /*value string*/) bool {
	if len(diff.Filter.HeaderFilter.Header) == 0 {
		return false
	}
	_, ok := diff.Filter.HeaderFilter.Header[header]

	return ok
}

func highestLstIntValue(lst []int) int {
	v := 0
	for _, i := range lst {
		if i > v {
			v = i
		}
	}
	return v
}

// Return an int array of 2 that holds 0/1 (min/max) and length diff
func lengthMinMaxDiff(x, y int) int {
	if x < y {
		return (y - x)
	}
	return x - y
}

// Compare two int type lists and return the differences presented in lstCurrent
// Return a true if the lists has one item in common
func lstIntShareItem(lstCompare, lstCurrent []int) bool {
	// Convert list1 into a map for quick lookups
	m := make(map[int]struct{})
	for _, i := range lstCompare {
		m[i] = struct{}{}
	}

	// Iterate over list2 and check if any item exists in the map
	for _, i := range lstCurrent {
		// An item is shared between the two lists
		if _, ok := m[i]; ok {
			return true
		}
	}
	return false
}

// Compare two string type lists and return the differences presented in lstCurrent
// Return a true if the lists has one item in common
func lstStringShareItem(lstCompare, lstCurrent []string) bool {
	// Convert list1 into a map for quick lookups
	m := make(map[string]struct{})
	for _, i := range lstCompare {
		m[i] = struct{}{}
	}

	// Iterate over list2 and check if any item exists in the map
	for _, i := range lstCurrent {
		// An item is shared between the two lists
		if _, ok := m[i]; ok {
			return true
		}
	}
	return false
}
