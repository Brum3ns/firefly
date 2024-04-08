package output

import (
	"fmt"
	"strconv"

	"github.com/Brum3ns/firefly/pkg/design"
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
	difference bool
	design     *design.Design

	style
}

type style struct {
	detail lipgloss.Style
}

func NewDisplay(detailed bool, design *design.Design) *Display {
	return &Display{
		difference: detailed,
		design:     design,
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

	stout := fmt.Sprintf("%s╭ \033[33m%s\033[0m Status:%s, Words:%s, Lines:%s, CL:%s, CT:%s, Time:%sms\n"+
		"╰╴Errors:[Body:%s, Header:%s] Diff:[Tag:%s, Attr:%s, AttrVal:%s, Words:%s, Comments:%s, Header:%s, HeaderVal:%s] %s\n",
		TERMINAL_CLEAR,
		strconv.Quote(d.Payload),
		// Response information
		d.design.StatusCode(d.Response.StatusCode),
		d.design.WordCount(d.Response.WordCount),
		d.design.LineCount(d.Response.LineCount),
		d.design.ContentLength(d.Response.ContentLength),
		d.design.ContentType(d.Response.ContentType),
		d.design.ResponseTime(d.Response.Time),
		//Extract:
		//d.design.Diff(d.Scanner.Extract.TotalHits),
		d.design.Diff(len(d.Scanner.Extract.PatternBody)),
		d.design.Diff(len(d.Scanner.Extract.PatternHeaders)),
		//Difference - HTMLNode:
		d.design.Diff(
			d.Scanner.Diff.HTMLNodeDiff.TagStartHits+
				d.Scanner.Diff.HTMLNodeDiff.TagEndHits+
				d.Scanner.Diff.HTMLNodeDiff.TagSelfClose,
		),
		d.design.Diff(d.Scanner.Diff.HTMLNodeDiff.AttributeHits),
		d.design.Diff(d.Scanner.Diff.HTMLNodeDiff.AttributeValueHits),
		d.design.Diff(d.Scanner.Diff.HTMLNodeDiff.WordsHits),
		d.design.Diff(d.Scanner.Diff.HTMLNodeDiff.CommentHits),
		//Difference - Headers:
		"000", //d.design.Diff(d.Scanner.Diff.Headers),
		"000", //d.design.Diff(d.Scanner.Diff.Headers),
		d.transformation(),
	)

	if d.difference {
		//Code...
	}
	fmt.Println(stout)

}

// Display payload transformation:
func (d Display) transformation() string {
	if len(d.Scanner.Transformation.Format) > 0 {
		return fmt.Sprintf(" Transformation: [\033[1;32m%s\033[0m => \033[1;32m%s\033[0m]", strconv.Quote(d.Scanner.Transformation.Payload), strconv.Quote(d.Scanner.Transformation.Format))
	}
	return ""
}
