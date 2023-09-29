package prepare

import (
	"fmt"
	"io"
	"net/http"
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
	Tag            map[string]int
	Words          map[string]int
	Comment        map[string]int
	Attribute      map[string]int
	AttributeValue map[string]int
	//Text           map[string]int
}

// !Note : (MUST be the same name as the "HTMLNode")
type HTMLNodeCombine struct {
	Tag            map[string][]int `json:"Tag"`
	Words          map[string][]int `json:"Words"`
	Comment        map[string][]int `json:"Comment"`
	Attribute      map[string][]int `json:"Attribute"`
	AttributeValue map[string][]int `json:"AttributeValue"`
	//Text           map[string][]int `json:"Text"`
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

func NewCombineHTMLNode() HTMLNodeCombine {
	return HTMLNodeCombine{
		Tag:            make(map[string][]int),
		Words:          make(map[string][]int),
		Comment:        make(map[string][]int),
		Attribute:      make(map[string][]int),
		AttributeValue: make(map[string][]int),
		//Text:           make(map[string][]int),
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

// Append a string to a list Works similar as append but do not append duplicates or empty strings
/* func appendUniqueInt(l []int, i int) []int {
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
*/
