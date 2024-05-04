package httpprepare

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type HTMLNode struct {
	TagStart       map[string]int
	TagEnd         map[string]int
	TagSelfClose   map[string]int
	Words          map[string]int
	Comment        map[string]int
	Attribute      map[string]int
	AttributeValue map[string]int
	//Text           map[string]int
}

// !Note : (MUST be the same name as the "HTMLNode")
type HTMLNodeCombine struct {
	TagStart       map[string][]int `json:"TagStart"`
	TagEnd         map[string][]int `json:"TagEnd"`
	TagSelfClose   map[string][]int `json:"TagSelfClose"`
	Words          map[string][]int `json:"Words"`
	Comment        map[string][]int `json:"Comment"`
	Attribute      map[string][]int `json:"Attribute"`
	AttributeValue map[string][]int `json:"AttributeValue"`
	//Text           map[string][]int `json:"Text"`
}

func NewHTMLNode() HTMLNode {
	return HTMLNode{
		TagStart:       make(map[string]int),
		TagEnd:         make(map[string]int),
		TagSelfClose:   make(map[string]int),
		Comment:        make(map[string]int),
		Attribute:      make(map[string]int),
		AttributeValue: make(map[string]int),
		Words:          make(map[string]int),
		//Text:           make(map[string]int),
	}
}

func NewCombineHTMLNode() HTMLNodeCombine {
	return HTMLNodeCombine{
		TagStart:       make(map[string][]int),
		TagEnd:         make(map[string][]int),
		TagSelfClose:   make(map[string][]int),
		Words:          make(map[string][]int),
		Comment:        make(map[string][]int),
		Attribute:      make(map[string][]int),
		AttributeValue: make(map[string][]int),
		//Text:           make(map[string][]int),
	}
}

// Collect HTML elements
// Note : (This function is used within the difference technique process)
func GetHTMLNode(body string) HTMLNode {
	htmlNode := NewHTMLNode()

	// Create a tokenizer and parse the HTML content
	tokenizer := html.NewTokenizer(strings.NewReader(body))

	for {
		typ := tokenizer.Next()

		// Reached the end of the response body:
		if typ == html.ErrorToken {
			if err := tokenizer.Err(); err == io.EOF {
				break
			}
		}

		//Check what type of token and sort it into the HTML [struct]ure:
		t := tokenizer.Token()
		switch typ {
		case html.StartTagToken:
			htmlNode.TagStart[t.Data]++
			for _, attr := range t.Attr {
				htmlNode.Attribute[attr.Key]++
				htmlNode.AttributeValue[attr.Val]++
			}

		case html.EndTagToken:
			htmlNode.TagEnd[t.Data]++
			for _, attr := range t.Attr {
				htmlNode.Attribute[attr.Key]++
				htmlNode.AttributeValue[attr.Val]++
			}

		case html.TextToken:
			for _, w := range strings.Fields(t.Data) {
				htmlNode.Words[w]++
			}
			//htmlNode.Text[t.Data]++

		case html.CommentToken:
			for _, w := range strings.Fields(t.Data) {
				htmlNode.Words[w]++
			}
			htmlNode.Comment[t.Data]++

		case html.SelfClosingTagToken:
			htmlNode.TagSelfClose[t.Data]++
			for _, attr := range t.Attr {
				htmlNode.Attribute[attr.Key]++
				htmlNode.AttributeValue[attr.Val]++
			}
		}
	}
	return htmlNode
}
