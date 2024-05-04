package output

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/httpprepare"
	"github.com/charmbracelet/lipgloss"
)

var (
	TERMINAL_CLEAR = "\r\x1b[2K"

	COLOR_BLACK  = lipgloss.Color("#000000")
	COLOR_WHITE  = lipgloss.Color("#D9DCCF")
	COLOR_GREY   = lipgloss.Color("#383838")
	COLOR_GREEN  = lipgloss.Color("#3AF191")
	COLOR_ORANGE = lipgloss.Color("#D98D00")
	COLOR_YELLOW = lipgloss.Color("#FFDF00")
	COLOR_RED    = lipgloss.Color("#EB2D3A")
)

type Display struct {
	ResultFinal
	detailed bool
	design   *design.Design

	style
}

type style struct {
	detail lipgloss.Style
}

func NewDisplay(detailed bool, design *design.Design) *Display {
	return &Display{
		detailed: detailed,
		design:   design,
	}
}

func (d Display) MakeStyles() {
	d.style = style{
		detail: lipgloss.NewStyle().Foreground(COLOR_GREY),
	}
}

// Display the information to the screen from a given structure (result data) to the command line interface (CLI) [show: on/off]) and any struct that *include JSON supported tags*.
// The function use color highlighting in the CLI by using a mixture of stderr and stout. The output values will be in stout
// version which makes it possible to support pipelining without including garbage in the values.
func (d *Display) ToScreen(result ResultFinal) {
	d.ResultFinal = result
	//print("\033[?25l")

	boxChar := "╰╴"
	if d.detailed {
		boxChar = "├╴"
	}

	diff := d.Scanner.Diff

	stout := fmt.Sprintf("%s╭ \033[33m%s\033[0m Status:%s, Words:%s, Lines:%s, CL:%s, CT:%s, Time:%sms\n"+
		"%sErrors:[Body:%s, Header:%s] Diff:[Tag:%s, Attr:%s, AttrVal:%s, Words:%s, Comments:%s, Header:%s] %s\n",
		TERMINAL_CLEAR,
		d.Payload,
		// Response information
		d.design.StatusCode(d.Response.StatusCode),
		d.design.WordCount(d.Response.WordCount),
		d.design.LineCount(d.Response.LineCount),
		d.design.ContentLength(d.Response.ContentLength),
		d.design.ContentType(d.Response.ContentType),
		d.design.ResponseTime(d.Response.Time),
		// Box Draw character
		boxChar,
		//Extract:
		//d.design.Diff(d.Scanner.Extract.TotalHits),
		d.design.Highlight(len(d.Scanner.Extract.PatternBody)),
		d.design.Highlight(len(d.Scanner.Extract.PatternHeaders)),
		//Difference - HTMLNode:
		d.design.Highlight(
			diff.HTMLResult.Appear.TagStartHits+
				diff.HTMLResult.Appear.TagEndHits+
				diff.HTMLResult.Appear.TagSelfCloseHits,
		),
		d.design.Highlight(diff.HTMLResult.Appear.AttributeHits),
		d.design.Highlight(diff.HTMLResult.Appear.AttributeValueHits),
		d.design.Highlight(diff.HTMLResult.Appear.WordsHits),
		d.design.Highlight(diff.HTMLResult.Appear.CommentHits),
		//Difference - Headers:
		d.design.Highlight(diff.HeaderResult.HeaderHits),
		d.transformation(),
	)

	if d.detailed {
		prefix := "|"
		stout += "\n├╴[Header]\n" +
			d.getDetailDiff("Appear", strings.Join(headerNodeToLst(prefix, diff.HeaderResult.Appear), "\n")) +
			d.getDetailDiff("Disapear", strings.Join(headerNodeToLst(prefix, diff.HeaderResult.Disappear), "\n")) +
			"\n├╴[HTML]\n" +
			d.getDetailDiff("Appear", strings.Join(htmlNodeToLst(prefix, diff.HTMLResult.Appear.HTMLNode), "\n")) +
			d.getDetailDiff("Disappear", strings.Join(htmlNodeToLst(prefix, diff.HTMLResult.Disappear.HTMLNode), "\n"))
	}
	fmt.Println(stout)
}

func (d *Display) getDetailDiff(title, s string) string {
	if len(s) > 0 {
		return fmt.Sprintf("├╴%s\n%s\n", title, (d.design.Color.GREY + s + d.design.Color.WHITE))
	}
	return ""
}

func htmlNodeLst(htmlNode string) {

}

func headerNodeToLst(prefix string, headerNode httpprepare.Header) []string {
	var lst []string
	// Todo ...
	for header, headerInfo := range headerNode {
		s := ""
		// Check the amount of items in the lists of header info:
		if len(headerInfo.Amount) == 1 && len(headerInfo.Values) == 1 {
			s = fmt.Sprintf("%s(%d) : %s: %s", prefix, headerInfo.Amount[0], header, headerInfo.Values[0])

			// This should not be possible, but just in case, output it if it happen to be
		} else {
			s = fmt.Sprintf("%s(%d) : %s:\n\t%s", prefix, headerInfo.Amount, header, strings.Join(headerInfo.Values, "\n"))
		}
		lst = append(lst, s)
	}

	return lst
}

func htmlNodeToLst(prefix string, htmlnode httpprepare.HTMLNode) []string {
	var lst []string

	if len(htmlnode.TagStart) > 0 {
		lst = append(lst,
			fmt.Sprintf("%sTag-start: %v", prefix, htmlnode.TagStart),
		)
	}
	if len(htmlnode.TagEnd) > 0 {
		lst = append(lst,
			fmt.Sprintf("%sTag-end: %v", prefix, htmlnode.TagEnd),
		)
	}
	if len(htmlnode.TagSelfClose) > 0 {
		lst = append(lst,
			fmt.Sprintf("%sTag-selfclose: %v", prefix, htmlnode.TagSelfClose),
		)
	}
	if len(htmlnode.Attribute) > 0 {
		lst = append(lst,
			fmt.Sprintf("%sAttribute: %v", prefix, htmlnode.Attribute),
		)
	}
	if len(htmlnode.AttributeValue) > 0 {
		lst = append(lst,
			fmt.Sprintf("%sAttributeValue: %v", prefix, htmlnode.AttributeValue),
		)
	}
	if len(htmlnode.Comment) > 0 {
		lst = append(lst,
			fmt.Sprintf("%sComment: %v", prefix, htmlnode.Comment),
		)
	}
	if len(htmlnode.Words) > 0 {
		lst = append(lst,
			fmt.Sprintf("%sWords: %v", prefix, htmlnode.Words),
		)
	}

	return lst
}

// Display payload transformation:
func (d Display) transformation() string {
	if len(d.Scanner.Transformation.Format) > 0 {
		return fmt.Sprintf(" Transformation: [\033[1;32m%s\033[0m => \033[1;32m%s\033[0m]", strconv.Quote(d.Scanner.Transformation.Payload), strconv.Quote(d.Scanner.Transformation.Format))
	}
	return ""
}
