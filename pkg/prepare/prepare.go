package prepare

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

type Prepare struct {
	properties Properties
	HTMLNode   HTMLNode
}

type Properties struct {
	//Fragment string
	Paylaod string
	Body    string
	Headers http.Header
}

type HTMLNode struct {
	Tag map[string]int
	//Text           map[string]int
	Words          map[string]int
	Comment        map[string]int
	Attribute      map[string]int
	AttributeValue map[string]int
}

type HTMLNodeCombine struct {
	Tag            map[string][]int
	Text           map[string][]int
	Words          map[string][]int
	Comment        map[string][]int
	Attribute      map[string][]int
	AttributeValue map[string][]int
}

func NewHTMLNode() HTMLNode {
	return HTMLNode{
		Tag: make(map[string]int),
		//Text:           make(map[string]int),
		Comment:        make(map[string]int),
		Attribute:      make(map[string]int),
		AttributeValue: make(map[string]int),
	}
}

func NewHTMLNodeCombine() HTMLNodeCombine {
	return HTMLNodeCombine{
		Tag:            make(map[string][]int),
		Text:           make(map[string][]int),
		Words:          make(map[string][]int),
		Comment:        make(map[string][]int),
		Attribute:      make(map[string][]int),
		AttributeValue: make(map[string][]int),
	}
}

func NewPrepare(p Properties) *Prepare {
	return &Prepare{
		properties: p,
		HTMLNode:   NewHTMLNode(),
	}
}

// Collect HTML elements
// Note : (This function is used within the difference technique process)
func GetHTMLNode(body string) HTMLNode {
	htmlNode := HTMLNode{
		Tag: make(map[string]int),
		//Text:           make(map[string]int),
		Comment:        make(map[string]int),
		Attribute:      make(map[string]int),
		AttributeValue: make(map[string]int),
		Words:          make(map[string]int),
	}

	// Create a tokenizer and parse the HTML content
	tokenizer := html.NewTokenizer(strings.NewReader(body))

	for {
		typ := tokenizer.Next()

		// Reached the end of the response body:
		if typ == html.ErrorToken {
			if err := tokenizer.Err(); err == io.EOF {
				break
			} else {
				fmt.Println("[\033[1;31mCRITICAL\033[0m] could not properly get the HTML node (\"content\") from the target HTTP response body! Error:", err)
			}
		}

		//Check what type of token and sort it into the HTML [struct]ure:
		switch typ {
		case html.StartTagToken, html.SelfClosingTagToken:
			t := tokenizer.Token()
			htmlNode.Tag[t.Data]++
			for _, attr := range t.Attr {
				htmlNode.Attribute[attr.Key]++
				htmlNode.AttributeValue[attr.Val]++
			}
		case html.TextToken:
			t := tokenizer.Token()
			for _, w := range strings.Fields(t.Data) {
				htmlNode.Words[w]++
			}
			//htmlNode.Text[t.Data]++

		case html.CommentToken:
			t := tokenizer.Token()
			for _, w := range strings.Fields(t.Data) {
				htmlNode.Words[w]++
			}
			htmlNode.Comment[t.Data]++
		}
	}
	return htmlNode
}

// Combine a list of HTMLNodes and combine its node into a new structure containing the map for each node and a list of int referring to its number of repetitions.func HTMLNodesCombined(htmlNode HTMLNode) HTMLNodeCombined {
func GetHTMLNodesCombined(htmlNodes []HTMLNode) HTMLNodeCombine {
	htmlNodeCombined := NewHTMLNodeCombine()
	for _, hNode := range htmlNodes {
		v := reflect.ValueOf(hNode)
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			for k, v := range v.Field(i).Interface().(map[string]int) {
				switch t.Field(i).Name {
				case "Tag":
					htmlNodeCombined.Tag[k] = appendUniqueInt(htmlNodeCombined.Tag[k], v)
				case "Text":
					htmlNodeCombined.Text[k] = appendUniqueInt(htmlNodeCombined.Text[k], v)
				case "Words":
					htmlNodeCombined.Words[k] = appendUniqueInt(htmlNodeCombined.Words[k], v)
				case "Comment":
					htmlNodeCombined.Comment[k] = appendUniqueInt(htmlNodeCombined.Comment[k], v)
				case "Attribute":
					htmlNodeCombined.Attribute[k] = appendUniqueInt(htmlNodeCombined.Attribute[k], v)
				case "AttributeValue":
					htmlNodeCombined.AttributeValue[k] = appendUniqueInt(htmlNodeCombined.AttributeValue[k], v)
				}
			}
		}
	}
	return htmlNodeCombined
}

// Combine one HTMLNode with another:
func CombineHTMLNode(combine HTMLNodeCombine, htmlNodes HTMLNode) HTMLNodeCombine {
	if reflect.DeepEqual(combine, HTMLNodeCombine{}) {
		combine = NewHTMLNodeCombine()
	}
	v := reflect.ValueOf(htmlNodes)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		for k, v := range v.Field(i).Interface().(map[string]int) {
			switch t.Field(i).Name {
			case "Tag":
				combine.Tag[k] = appendUniqueInt(combine.Tag[k], v)
			case "Text":
				combine.Text[k] = appendUniqueInt(combine.Text[k], v)
			case "Words":
				combine.Words[k] = appendUniqueInt(combine.Words[k], v)
			case "Comment":
				combine.Comment[k] = appendUniqueInt(combine.Comment[k], v)
			case "Attribute":
				combine.Attribute[k] = appendUniqueInt(combine.Attribute[k], v)
			case "AttributeValue":
				combine.AttributeValue[k] = appendUniqueInt(combine.AttributeValue[k], v)
			}
		}
	}
	return combine
}

// Append a string to a list Works similar as append but do not append duplicates or empty strings
func appendUniqueInt(l []int, i int) []int {
	if len(l) == 0 {
		return append(l, i)
	}
	for _, item := range l {
		if item == i {
			return l
		}
	}
	return append(l, i)
}
