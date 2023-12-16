package output

import (
	"fmt"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/firefly/global"
)

type Display struct {
	ResultFinal
	design *design.Design
}

func NewDisplay(design *design.Design) *Display {
	return &Display{
		design: design,
	}
}

// Display the information to the screen from a given structure (result data) to the command line interface (CLI) [show: on/off]) and any struct that *include JSON supported tags*.
// The function use color highlighting in the CLI by using a mixture of stderr and stout. The output values will be in stout
// version which makes it possible to support pipelining without including garbage in the values.
func (d *Display) ToScreen(result ResultFinal) {
	d.ResultFinal = result

	stout := fmt.Sprintf("╭ \033[33m%s\033[0m Status:%s, Words:%s, Lines:%s, CL:%s, CT:%s, Time:%sms\n"+
		"╰╴Errors:[Hits:%s, Body:%s, Header:%s] Diff:[Tag:%s, Attr:%s, AttrValues:%s, Words:%s, Comments:%s]",
		d.Payload,
		d.design.StatusCode(d.Response.StatusCode),
		d.design.WordCount(d.Response.WordCount),
		d.design.LineCount(d.Response.LineCount),
		d.design.ContentLength(d.Response.ContentLength),
		d.design.ContentType(d.Response.ContentType),
		d.design.ResponseTime(d.Response.Time),
		//Extract:
		d.design.IsDiff(d.Scanner.Extract.TotalHits),
		d.design.IsDiff(len(d.Scanner.Extract.PatternBody)),
		d.design.IsDiff(len(d.Scanner.Extract.PatternHeaders)),
		//Difference:
		d.design.IsDiff(d.Scanner.Diff.TagHits),
		d.design.IsDiff(d.Scanner.Diff.AttributeHits),
		d.design.IsDiff(d.Scanner.Diff.AttributeValueHits),
		d.design.IsDiff(d.Scanner.Diff.WordsHits),
		d.design.IsDiff(d.Scanner.Diff.CommentHits),
	)
	if d.Scanner.Transformation.OK {
		stout += d.transformation()
	}

	fmt.Println(global.TERMINAL_CLEAR + stout + "\n")
}

// Display payload transformation:
func (d Display) transformation() string {
	return fmt.Sprintf(" Transformation: [\033[1;32m%s\033[0m => \033[1;32m%s\033[0m]", d.Scanner.Transformation.Payload, d.Scanner.Transformation.Format)
}
