package difference

import (
	"math"
	"net/http"
	"sync"

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
}

type ResponseBody struct {
	Body          string
	HtmlNode      prepare.HTMLNode
	KnownHTMLNode []prepare.HTMLNode
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

func NewDifference(p Properties) *Difference {
	return &Difference{
		Properties: &p,
	}
}

// Run the [diff]erence enumiration process
func (diff *Difference) Run() Result {
	result := Result{
		HTMLNode: prepare.NewHTMLNode(),
	}

	var wg sync.WaitGroup

	for _, node := range diff.ResponseBody.KnownHTMLNode {
		wg.Add(1)
		go func(node prepare.HTMLNode) {
			result.Words, _ = GetDifference(diff.HtmlNode.Words, node.Words)
			result.AttributeValue, _ = GetDifference(diff.HtmlNode.AttributeValue, node.AttributeValue)
			result.Comment, result.CommentHits = GetDifference(diff.HtmlNode.Comment, node.Comment)
			result.Attribute, result.AttributeHits = GetDifference(diff.HtmlNode.Attribute, node.Attribute)
			result.Tag, result.TagHits = GetDifference(diff.HtmlNode.Tag, node.Tag)
			wg.Done()
		}(node)
	}

	//Wait until all [diff]erence analyses have been completed
	wg.Wait()

	result.Words, result.WordsHits = DetectRandomness(result.Words)
	result.AttributeValue, result.AttributeValueHits = DetectRandomness(result.AttributeValue)

	//Provide success status
	if (result.TagHits + result.WordsHits + result.CommentHits + result.AttributeHits + result.AttributeValueHits) > 0 {
		result.OK = true
	}

	return result
}

func (diff *Difference) AppendKnownHTMLNode(htmlNode prepare.HTMLNode) {
	diff.Properties.ResponseBody.KnownHTMLNode = append(diff.Properties.ResponseBody.KnownHTMLNode, htmlNode)
}
func (diff *Difference) AppendKnownHeaders(headers http.Header) {
	diff.Properties.ResponseHeaders.KnownHeaders = append(diff.Properties.ResponseHeaders.KnownHeaders, headers)
}

// Compare and return the difference between two maps
// This function is highly adapted for speed and preformance. It's capable of comparing all items within a map using any form of nested loops.
func GetDifference(current, known map[string]int) (map[string]int, int) {
	var (
		ItemHits int
		m_diff   = make(map[string]int)
	)

	if current != nil || known != nil {
		// [Variable DESC]
		// a2_ids   : (Array) - Holds the *order* of the original and fuzzed map (Verified = 1, Fuzzed = 0) - [Note] : Only *the order* of value 1/0 do matter
		// m_holder : (Map)   - Combined two maps (0/1) that holds all the items and amount into one

		var (
			holder    = make(map[int]map[string]int)
			a2_ids, _ = lengthMinMax(len(current), len(known))
		)
		holder[a2_ids[0]] = current
		holder[a2_ids[1]] = known

		//Deep copy of the min length map:
		m_copy := make(map[string]int)
		for k, v := range holder[1] {
			m_copy[k] = v
		}

		for item, amount := range holder[0] {
			if _, ok := m_diff[item]; ok {
				continue
			}

			hit := 0

			//Both map share same item, if they have different lengths = diff. ("if statements" mush be clustered)
			if _, ok := holder[1][item]; ok {
				if holder[1][item] != amount {
					hit++
				}
				delete(m_copy, item)

				//Maps do not share the same item = diff
			} else {
				hit++
			}
			//Check if there was a diff:
			if hit == 1 {
				ItemHits++
				m_diff[item] = (holder[0][item] - holder[1][item]) //Add item | lenDiff:
			}
		}
		//Add the rest of the items that appeared in the min map ('m_copy') but not max map:
		if len(m_copy) != 0 {
			for item, value := range m_copy {
				m_diff[item] = value
			}
		}
	}

	return m_diff, ItemHits
}

var vocals = map[rune]int{
	'a': 0,
	'e': 0,
	'i': 0,
	'o': 0,
	'u': 0,
	'A': 1,
	'E': 1,
	'I': 1,
	'O': 1,
	'U': 1,
	//The digits 0, 3 and 4 can sometimes refer to as the letters o, e and a:
	//'0': 1,
	//'3': 1,
	//'4': 1,
}

func DetectRandomness(m map[string]int) (map[string]int, int) {
	totalHits := 0
	for s, hits := range m {
		//Quick randomness check:
		if _, ok := IsRandom(s); ok {
			delete(m, s)

			//Math algoritm (*shannon entropy*) to verify:
		} else {
			totalHits += hits
		}
	}
	return m, totalHits
}

// Return a map that contains the charFrequency and true/false if the value is random (*very likely to be random*):
func IsRandom(s string) (map[rune]int, bool) {
	var (
		hit            = 0
		consonantInRow = 3
		charFrequency  = make(map[rune]int)
	)

	//Calculate constant that comes in a row: (usually 3-4+ result in a random string)
	for _, r := range s {
		if _, ok := vocals[r]; ok {
			hit = 0
		}
		hit++
		//Likely to be random:
		if hit == consonantInRow {

			//Save random values and make an "average entropy value" adapted to the applicato so we can learn the randomness sturcture
			return charFrequency, true
		}
		charFrequency[r]++
	}
	return charFrequency, false
}

func Entropy(s string, charFrequency map[rune]int) float64 {
	var (
		value float64
		runes = []rune(s)
	)
	//In case we do not have the char frequency set:
	if charFrequency == nil {
		charFrequency = make(map[rune]int)
		for _, r := range runes {
			charFrequency[r]++
		}
	}
	characters := float64(len(runes))
	for _, count := range charFrequency {
		probability := float64(count) / characters
		value -= probability * math.Log2(probability)
	}
	return value
}

// Return an int array of 2 that holds 0/1 (min/max) and length diff
func lengthMinMax(x, y int) ([2]int, int) {
	if x < y {
		return [2]int{0, 1}, (y - x)
	}
	return [2]int{1, 0}, (x - y)
}
