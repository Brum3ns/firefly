package prepare

import (
	"fmt"

	"regexp"
	"strings"

	"github.com/Brum3ns/FireFly/pkg/firefly/types"
	"github.com/Brum3ns/FireFly/pkg/functions"
	fc "github.com/Brum3ns/FireFly/pkg/functions"
	"golang.org/x/net/html"
)

type BodyData struct {
	t *types.VerifyInfoData
}

func Body(body, pyld string) types.VerifyInfoData {
	data := &BodyData{
		t: types.Setup_VerifyInfoData(),
	}
	data.t.PayloadMark = "{FIREFLY_PAYLOAD}"
	data.t.PayloadRegex = functions.PayloadRegexMark(1)

	//Prepare body (filter out payload reflectations):
	data.t.BodyPrepared = data.Body_PayloadMark(body)

	//Extract data from the new body that have marked the reflected payload positions and use it to gather more data:
	data.t.PayloadReflectCount = data.Count_PayloadReflect()
	data.Words() //<-- Improve with net/html
	data.Lines()
	data.HTMLAttributes()
	data.HTMLCounter()

	return *data.t
}

/**Count all the data amount*/
func (bd *BodyData) HTMLCounter() {
	for attr, _ := range bd.t.HTMLAttr {
		bd.t.HTMLAttrCount += bd.t.HTMLAttr[attr]
	}
	for attrVal, _ := range bd.t.HTMLAttrValue {
		bd.t.HTMLAttrValueCount += bd.t.HTMLAttrValue[attrVal]
	}
	for tag, _ := range bd.t.HTMLTags {
		bd.t.HTMLTagCount += bd.t.HTMLTags[tag]
	}
}

/**Return a new body where all the payload reflectation are replaced with the 'payloadMark' value*/
func (bd *BodyData) Body_PayloadMark(BODY string) string {
	regx := regexp.MustCompile(bd.t.PayloadRegex)
	return regx.ReplaceAllString(BODY, bd.t.PayloadMark)
}

/**Return the amount of times the payload reflected in the body*/
func (bd *BodyData) Count_PayloadReflect() int {
	return strings.Count(bd.t.BodyPrepared, bd.t.PayloadMark)
}

/**Return all words that are included in the body and count the amount for each unique*/
func (bd *BodyData) Words() {
	for _, word := range strings.Fields(bd.t.BodyPrepared) {
		bd.t.WordCount += 1
		bd.t.Words[word] += 1
	}
}

/**Return all lines that are included in the body and count the amount for each unique*/
func (bd *BodyData) Lines() {
	for _, line := range strings.Split(bd.t.BodyPrepared, "\n") {
		bd.t.LineCount += 1
		bd.t.Lines[line] += 1
	}
}

/**Return all HTML attributes that are included in the body and count the amount for each*/
func (bd *BodyData) HTMLAttributes() {
	HTMLDocument, err := html.Parse(strings.NewReader(bd.t.BodyPrepared))
	fc.IFError("f", err)

	//Run a isolated function to extract the HTML nodes:
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode {
			//Add HTML tag to map:
			bd.t.HTMLTags[fmt.Sprintf("%v", node.DataAtom)] += 1

			//Extract all attribute(s) and their value from the current HTML tag:
			for _, attr := range node.Attr {
				//Add attribute name to map:
				if len(attr.Key) > 0 {
					bd.t.HTMLAttr[attr.Key] += 1

					//Add attribute value to map:
					if len(attr.Val) > 0 {
						bd.t.HTMLAttrValue[attr.Val] += 1
					}
					/* if !fc.InLst(bd.t.HTML[attr.Key], attr.Val) {// -> Think adding duplicate in this case is better
						bd.t.HTML[attr.Key] = append(bd.t.HTML[attr.Key], attr.Val)
					} */
					bd.t.HTML[attr.Key] = append(bd.t.HTML[attr.Key], attr.Val)
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(HTMLDocument)
}
