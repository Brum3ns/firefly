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
	duration := durationFormat(time.Since(p.Stats.Timer))

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
		p.Stats.Request.Count,
		p.Stats.Scanner.Count,
		p.Stats.Behavior,
		p.Stats.Request.Filtered,
		p.Stats.Error,
		duration[0], duration[1], duration[2],
	)

}

// Take "time.Duration" value and return list of containing the hour, minute and secounds
func durationFormat(d time.Duration) [3]time.Duration {
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return [3]time.Duration{h, m, s}
}
