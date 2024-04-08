package httpdiff

import (
	"net/http"

	"github.com/Brum3ns/firefly/pkg/httpprepare"
	"github.com/Brum3ns/firefly/pkg/randomness"
)

type Difference struct {
	*Properties
}

type Properties struct {
	Payload       string
	PayloadVerify string
	*ResponseBody
	*ResponseHeaders
	CompareHTMLNodes httpprepare.HTMLNodeCombine
	//CompareHeaders   httpprepare.Header
}

type ResponseBody struct {
	Body     string
	HtmlNode httpprepare.HTMLNode
}

type ResponseHeaders struct {
	HeaderString string
	Headers      http.Header
	HeadersKnown http.Header
}

type HeaderResult struct {
	OK     bool
	header http.Header
}

type HTMLResult struct {
	OK bool
	HTMLNodeDiff
}

type HTMLNodeDiff struct {
	TagStartHits       int
	TagEndHits         int
	TagSelfClose       int
	WordsHits          int
	CommentHits        int
	AttributeHits      int
	AttributeValueHits int
	httpprepare.HTMLNode
}

type diffNode struct {
	hit             int
	checkRandomness bool
	data            map[string]int
}

func NewDifference(p Properties) *Difference {
	return &Difference{
		Properties: &p,
	}
}

func newDiffNode() diffNode {
	return diffNode{
		data: make(map[string]int),
	}
}

// Run the [diff]erence enumiration process for the HTML node
func (diff *Difference) GetHTMLNodeDiff() HTMLResult {
	totalHits := 0
	storage := []diffNode{}

	//!Note : Order for "current" and "known" togther with the list length *MUST* be the same:
	current := [7]diffNode{
		{data: diff.Properties.HtmlNode.TagStart},
		{data: diff.Properties.HtmlNode.TagEnd},
		{data: diff.Properties.HtmlNode.TagSelfClose},
		{data: diff.Properties.HtmlNode.Words, checkRandomness: true},
		{data: diff.Properties.HtmlNode.Comment, checkRandomness: true},
		{data: diff.Properties.HtmlNode.Attribute},
		{data: diff.Properties.HtmlNode.AttributeValue, checkRandomness: true},
	}
	known := [7]map[string][]int{
		diff.Properties.CompareHTMLNodes.TagStart,
		diff.Properties.CompareHTMLNodes.TagEnd,
		diff.Properties.CompareHTMLNodes.TagSelfClose,
		diff.Properties.CompareHTMLNodes.Words,
		diff.Properties.CompareHTMLNodes.Comment,
		diff.Properties.CompareHTMLNodes.Attribute,
		diff.Properties.CompareHTMLNodes.AttributeValue,
	}

	for i := 0; i < len(current); i++ {
		//Detect difference
		diffCurrent := nodeDiff(current[i], known[i], diff.Payload)

		storage = append(storage, diffCurrent)
		totalHits += diffCurrent.hit
	}

	return HTMLResult{
		// !Note : The order for this *MUST* follow the same as "known" and "current" above
		OK: (totalHits > 0),
		HTMLNodeDiff: HTMLNodeDiff{
			TagStartHits:       storage[0].hit,
			TagEndHits:         storage[1].hit,
			TagSelfClose:       storage[2].hit,
			WordsHits:          storage[3].hit,
			CommentHits:        storage[4].hit,
			AttributeHits:      storage[5].hit,
			AttributeValueHits: storage[6].hit,
			HTMLNode: httpprepare.HTMLNode{
				TagStart:       storage[0].data,
				TagEnd:         storage[1].data,
				TagSelfClose:   storage[2].data,
				Words:          storage[3].data,
				Comment:        storage[4].data,
				Attribute:      storage[5].data,
				AttributeValue: storage[6].data,
			},
		},
	}
}

func (diff *Difference) GetHeadersDiff() HeaderResult {
	return HeaderResult{}
}

/*
	func lstDiff(lst1, lst2 []string) []string {
		diff := []string{}
		set := make(map[string]bool)

		// Populate set with elements from list2
		for _, elem := range lst2 {
			set[elem] = true
		}

		// Check each element in list1 against the set
		for _, elem := range lst1 {
			if !set[elem] {
				diff = append(diff, elem)
			}
		}
		return diff
	}
*/
func nodeDiff(current diffNode, known map[string][]int, payload string) diffNode {
	diffhtmlNode := newDiffNode()

	for currentItem, currentValue := range current.data {
		// Set the diff as true by default
		diff := true
		amountDiff := 0

		// Check if the current item exists in the known
		if knownValues, ok := known[currentItem]; ok {
			// Compare with the known values
			for _, knownValue := range knownValues {
				if currentValue == knownValue {
					diff = false
					break

					// If the current value isen't in the list check the amount diff and update to the highest amount diff
				} else if v := lengthMinMaxDiff(currentValue, knownValue); v > amountDiff {
					amountDiff = v
				}
			}
		} else if amountDiff == 0 {
			amountDiff = currentValue
		}

		// FIX THIS
		if diff && currentItem != payload {
			// Check false positive
			if !current.checkRandomness || (current.checkRandomness && !randomness.IsRandom(currentItem, true)) {
				diffhtmlNode.data[currentItem] = amountDiff
				diffhtmlNode.hit += amountDiff
			}
		}
	}

	//fmt.Println("DIFF", payload, diffhtmlNode) //DEBUG - REMOVE
	return diffhtmlNode
}

// Return an int array of 2 that holds 0/1 (min/max) and length diff
func lengthMinMaxDiff(x, y int) int {
	if x < y {
		return (y - x)
	}
	return x - y
}
