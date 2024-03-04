package statistics

import (
	"fmt"
	"os"
	"time"

	"github.com/Brum3ns/firefly/pkg/firefly/global"
)

type ProgressBar struct {
	Counter         int
	time            time.Time
	delay           time.Duration
	Stats           *Statistic
	Classic         []string
	SegmentedDigits []string
}

func NewProgressBar(delayMS int, statistic *Statistic) ProgressBar {
	return ProgressBar{
		Counter: 0,
		delay:   time.Duration(delayMS) * time.Millisecond,
		time:    time.Now(),
		Stats:   statistic,
		Classic: []string{"⠙", "⠸", "⠴", "⠦", "⠇", "⠋"},
	}
}

// Display the progress the statistic structure
func (p *ProgressBar) Print() {
	t := p.Stats.GetTime()

	if time.Now().After(p.time) {
		p.Counter++
		p.time = time.Now().Add(p.delay)

	}
	if p.Counter >= len(p.Classic) {
		p.Counter = 0
	}

	//fmt.Println(global.TERMINAL_CLEAR, "->", p.delay, x)
	fmt.Fprintf(os.Stderr, "%s%s Request:[%d], Scanned:[%d], Behavior:[%d], Filtered:[%d], Error:[%d], Time:[%d:%02d:%02d]",
		global.TERMINAL_CLEAR,
		p.Classic[p.Counter],
		p.Stats.Request.GetCount(),
		p.Stats.Scanner.GetCount(),
		p.Stats.Behavior.GetCount(),
		p.Stats.Response.GetFilterCount(),
		p.Stats.Request.GetErrorCount(),
		t[0], t[1], t[2],
	)
}
