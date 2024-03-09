package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/pkg/statistics"
	"github.com/charmbracelet/bubbles/spinner"
)

type ProgressBar struct {
	Counter         int
	time            time.Time
	delay           time.Duration
	Stats           *statistics.Statistic
	SegmentedDigits []string
	spinner         spinner.Model
}

func NewProgressBar(delayMS int, statistic *statistics.Statistic) ProgressBar {
	return ProgressBar{
		Counter: 0,
		delay:   time.Duration(delayMS) * time.Millisecond,
		time:    time.Now(),
		Stats:   statistic,
		spinner: spinner.New(spinner.WithSpinner(spinner.MiniDot)),
	}
}

// Display the progress the statistic structure
func (p *ProgressBar) Print() {
	t := p.Stats.GetTime()
	p.spinner, _ = p.spinner.Update(spinner.Tick())

	//fmt.Println(global.TERMINAL_CLEAR, "->", p.delay, x)
	fmt.Fprintf(os.Stderr, "%s%s Request:[%d], Scanned:[%d], Behavior:[%d], Filtered:[%d], Error:[%d], Time:[%d:%02d:%02d]",
		global.TERMINAL_CLEAR,
		p.spinner.View(),
		p.Stats.Request.GetCount(),
		p.Stats.Scanner.GetCount(),
		p.Stats.Behavior.GetCount(),
		p.Stats.Response.GetFilterCount(),
		p.Stats.Request.GetErrorCount(),
		t[0], t[1], t[2],
	)
}
