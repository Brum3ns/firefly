package difference

import (
	"net/http"

	"github.com/Brum3ns/firefly/pkg/algorithm"
	"github.com/Brum3ns/firefly/pkg/prepare"
)

type Difference struct {
	*Properties
}

type Properties struct {
	Payload       string
	PayloadVerify string
	*ResponseBody
	*ResponseHeaders
	CompareHTMLNodes prepare.HTMLNodeCombine
}

type ResponseBody struct {
	Body     string
	HtmlNode prepare.HTMLNode
}

type ResponseHeaders struct {
	HeaderString string
	Headers      http.Header
	KnownHeaders []http.Header
}

type Result struct {
	OK                 bool
	TagHits            int
	WordsHits          int
	CommentHits        int
	AttributeHits      int
	AttributeValueHits int
	prepare.HTMLNode
}

type htmlnode struct {
	diff       map[string]int
	randomness bool
	hit        int
}

func NewDifference(p Properties) *Difference {
	return &Difference{
		Properties: &p,
	}
}

// Run the [diff]erence enumiration process
func (diff *Difference) Run() Result {
	totalHits := 0
	storage := []htmlnode{}

	//Order for "current" and "known" togther with the list length *MUST* be the same:
	current := [5]htmlnode{
		{diff: diff.Properties.HtmlNode.Tag},
		{diff: diff.Properties.HtmlNode.Words, randomness: true},
		{diff: diff.Properties.HtmlNode.Comment, randomness: true},
		{diff: diff.Properties.HtmlNode.Attribute},
		{diff: diff.Properties.HtmlNode.AttributeValue, randomness: true},
	}
	known := [5]map[string][]int{
		diff.Properties.CompareHTMLNodes.Tag,
		diff.Properties.CompareHTMLNodes.Words,
		diff.Properties.CompareHTMLNodes.Comment,
		diff.Properties.CompareHTMLNodes.Attribute,
		diff.Properties.CompareHTMLNodes.AttributeValue,
	}

	for i := 0; i < len(current); i++ {
		//Detect difference
		diffCurrent, hit := GetDiff(current[i], known[i], diff.Payload)

		storage = append(storage, diffCurrent)
		totalHits += hit
	}

	return Result{
		OK:                 (totalHits > 0),
		TagHits:            storage[0].hit,
		WordsHits:          storage[1].hit,
		CommentHits:        storage[2].hit,
		AttributeHits:      storage[3].hit,
		AttributeValueHits: storage[4].hit,
		HTMLNode: prepare.HTMLNode{
			Tag:            storage[0].diff,
			Words:          storage[1].diff,
			Comment:        storage[2].diff,
			Attribute:      storage[3].diff,
			AttributeValue: storage[4].diff,
		},
	}
}

func GetDiff(current htmlnode, known map[string][]int, payload string) (htmlnode, int) {
	hit := 0
	for item, valueCurrent := range current.diff {
		diff := true
		valueHighestDiff := 0

		//If the key exists and they share the same value, delete the key from the "current" map:
		if lstValue, ok := known[item]; ok {
			if len(lstValue) == 1 && valueCurrent == lstValue[0] {
				diff = false

			} else {
				for _, value := range lstValue {
					if valueCurrent == value {
						diff = false
						break

						//Calculate the value diff and update the highest value diff:
					} else if valueDiff := lengthMinMaxDiff(valueCurrent, value); valueDiff > valueHighestDiff {
						valueHighestDiff = valueDiff
					}
				}
			}
		}
		//Confirm that the item is a diff and not a false positive / random dynamic value:
		if diff && item != payload {
			if _, ok := algorithm.IsRandom(item); !ok {
				current.diff[item] = valueHighestDiff
				current.hit += valueHighestDiff
				continue
			}
		}
		//Not a diff, delete the item from the core htmlnode "current map":
		delete(current.diff, item)
	}

	return current, hit
}

// Return an int array of 2 that holds 0/1 (min/max) and length diff
func lengthMinMaxDiff(x, y int) int {
	if x < y {
		return (y - x)
	}
	return x - y
}
