package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	AUTHOR             = "By: @yeswehack / Brumens"
	TOOL               = "FireFly (v1.3.1)"
	MODE_DARK          = "#383838"
	MODE_LIGHT         = "#D9DCCF"
	BACKGROUND_PATTERN = "萤火虫"

	COLOR_BLACK  = lipgloss.Color("#000000")
	COLOR_WHITE  = lipgloss.Color("#D9DCCF")
	COLOR_GREY   = lipgloss.Color("#383838")
	COLOR_GREEN  = lipgloss.Color("#3AF191")
	COLOR_ORANGE = lipgloss.Color("#D98D00")
	COLOR_YELLOW = lipgloss.Color("#FFDF00")
	COLOR_RED    = lipgloss.Color("#EB2D3A")
)

type TerminalUI struct {
	Data
	Style              *style
	listPayload        list.Model
	listTransformation list.Model
	spinner            string
	window             int
}

type style struct {
	item         lipgloss.Style
	core         lipgloss.Style
	banner       lipgloss.Style
	author       lipgloss.Style
	column       lipgloss.Style
	payload      lipgloss.Style
	infobox      lipgloss.Style
	bar          lipgloss.Style
	done         lipgloss.Style
	table        lipgloss.Style
	spinner      lipgloss.Style
	window       lipgloss.Style
	windowActive lipgloss.Style
	adColor      lipgloss.AdaptiveColor
	oversize     func(s ...string) string
	header       func(s ...string) string
	width
}

type width struct {
	columnRight int
	columnleft  int
	columnMid   int
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func NewTerminalUI() *TerminalUI {
	demoList := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
	}

	return &TerminalUI{
		Style:              NewStyle(),
		listPayload:        list.New(demoList, list.NewDefaultDelegate(), 0, 0),
		listTransformation: list.New(demoList, list.NewDefaultDelegate(), 0, 0),
	}
}

// Make style definitions for the terminal user-interface (UI)
func NewStyle() *style {
	border := lipgloss.NormalBorder()
	adaptiveColor := lipgloss.AdaptiveColor{Light: MODE_LIGHT, Dark: MODE_DARK}
	s := &style{
		adColor: adaptiveColor,
		bar:     lipgloss.NewStyle(),
		spinner: lipgloss.NewStyle(),
		core:    lipgloss.NewStyle().Padding(1, 2, 1, 2),
		item:    lipgloss.NewStyle().Padding(0, 1, 0, 1),
		banner:  lipgloss.NewStyle().Foreground(COLOR_YELLOW),
		payload: lipgloss.NewStyle().Foreground(COLOR_ORANGE),
		author:  lipgloss.NewStyle().Foreground(COLOR_WHITE),

		done: lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			SetString("✓").
			Foreground(COLOR_GREEN),

		infobox: lipgloss.NewStyle().
			Padding(0, 1, 0, 1).
			Background(COLOR_WHITE).
			Foreground(COLOR_BLACK),

		oversize: lipgloss.NewStyle().
			Foreground(COLOR_GREY).
			Render,

		header: lipgloss.NewStyle().
			Foreground(COLOR_GREY).
			Bold(true).
			Render,

		table: lipgloss.NewStyle().
			Foreground(COLOR_GREY).
			Height(20),

		window: lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder()),

		windowActive: lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69")),

		column: lipgloss.NewStyle().
			Align(lipgloss.Left).
			Border(border, false, true, true, false).
			BorderForeground(adaptiveColor).
			Height(20),
	}

	return s
}

func (t *TerminalUI) Render(r Data) string {
	s := t.Style
	screenWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	s.width.setColumns(screenWidth)

	// Prepare banner
	banner := t.banner(screenWidth)

	//Get the process time
	timerDuration := r.stats.GetTime()
	timer := fmt.Sprintf("%d:%02d:%02d", timerDuration[0], timerDuration[1], timerDuration[2])

	//	columnWidth := screenWidth / 2
	stringUI := strings.Builder{}

	// Add the banner at the top
	stringUI.WriteString(banner + "\n")

	//Make all the columns
	{
		columnLeft := s.column.Render(
			s.HeaderRender(lipgloss.Center, "Analyzer", s.width.columnleft) +
				lipgloss.JoinVertical(lipgloss.Left,
					s.viewItem("Hits........", 0),
					s.viewItem("Chars.......", 0),
					s.viewItem("Patterns....", 0),
					s.viewItem("Mutations...", 0),
					s.viewItem("Reflections.", 0),
				) + "\n" +
				s.HeaderRender(lipgloss.Center, "Scanner", s.width.columnleft) +
				lipgloss.JoinVertical(lipgloss.Left,
					s.viewItem("Behavior...", r.stats.Behavior.GetCount()),
					s.viewItem("Difference.", r.stats.Difference.GetCount()),
					s.viewItem("Transform..", r.stats.Transformation.GetCount()),
				) + "\n" +
				s.HeaderRender(lipgloss.Center, "Request", s.width.columnleft) +
				lipgloss.JoinVertical(lipgloss.Left,
					s.viewItem("Requests.....", r.stats.Request.GetCount()),
					s.viewItem("Responses....", r.stats.Response.GetCount()),
					s.viewItem("Filtered.....", r.stats.Request.GetFilterCount()),
					s.viewItem("Forbidden....", r.stats.Request.GetCountForbidden()),
					s.viewItem("Request Err..", r.stats.Request.GetErrorCount()),
					s.viewItem("Response Err.", r.stats.Response.GetErrorCount()),
					s.viewItem("Average Time.", r.stats.Response.GetAverageTime()),
				),
		)
		columMid := s.column.Render(
			s.HeaderRender(lipgloss.Center, "Payload", s.width.columnMid) +
				lipgloss.JoinVertical(lipgloss.Left),
			//t.listPayload.View(),
		)

		columnRight := s.column.Render(
			s.HeaderRender(lipgloss.Center, "Transformation", s.width.columnMid) +
				lipgloss.JoinVertical(lipgloss.Left), //t.listTransformation.View(),

		)

		//Write all columns
		stringUI.WriteString(
			lipgloss.JoinHorizontal(lipgloss.Top,
				lipgloss.JoinHorizontal(
					lipgloss.Top,
					columnLeft,
					columMid,
					columnRight,
				),
			) + "\n",
		)
	}

	// Status bar
	{
		statusBar := s.bar.Render(
			lipgloss.JoinHorizontal(lipgloss.Left,
				s.infobox.Background(COLOR_GREEN).Render(timer),
				s.viewItem("Current", r.ResultFinal.Payload, s.payload),
			),
		)
		stringUI.WriteString(statusBar + "\n")
	}

	if screenWidth > 0 {
		t.Style.core = t.Style.core.MaxWidth(screenWidth)
	}

	return t.Style.core.Render(stringUI.String())
}

func (s *style) HeaderRender(position lipgloss.Position, header string, width int) string {
	return lipgloss.JoinVertical(position,
		lipgloss.Place(width, 1,
			lipgloss.Center, lipgloss.Center,
			s.header(header),
			lipgloss.WithWhitespaceChars("─"),
			lipgloss.WithWhitespaceForeground(s.adColor),
		),
	) + "\n"
}

/* func (s *style) TableRender(position lipgloss.Position, table string, width int) string {
	return lipgloss.JoinVertical(position,
		lipgloss.Place(width, 1,
			lipgloss.Center, lipgloss.Center,
			table,
			lipgloss.WithWhitespaceForeground(s.adColor),
		),
	) + "\n"
}
*/
// Set the window that will be in focus
func (m *TerminalUI) SetSpinner(spinner string) {
	m.spinner = spinner
}

func (m *TerminalUI) SetWindow(wid int) {
	m.window = wid
}

// Set column width and return the width used from the screen width
func (w *width) setColumns(screenWidth int) int {
	sw := screenWidth / 3
	widthUsed := 0

	//Left column
	if sw >= 24 {
		w.columnleft = 24
		widthUsed += 24
		sw = (screenWidth - 24)
	} else {
		w.columnleft = sw
		widthUsed += sw
	}
	sw = (screenWidth - w.columnleft) / 2
	w.columnRight = sw
	w.columnMid = sw

	return widthUsed
}

// Show the banner in the terminal UI
func (t *TerminalUI) banner(width int) string {
	banner := lipgloss.Place(width, 1,
		lipgloss.Center, lipgloss.Center,
		t.Style.banner.Render(TOOL),
		lipgloss.WithWhitespaceChars(BACKGROUND_PATTERN),
		lipgloss.WithWhitespaceForeground(t.Style.adColor),
	)
	author := lipgloss.Place(width, 1,
		lipgloss.Center, lipgloss.Center,
		t.Style.author.Render(AUTHOR),
		lipgloss.WithWhitespaceChars(BACKGROUND_PATTERN),
		lipgloss.WithWhitespaceForeground(t.Style.adColor),
	)
	return banner + "\n" + author
}

// Take a argument name and a value that represent an unkown value.
// The value vill be escaped with the 'strconv.Quote()' function to avoid ASNI injections.
// If the value given is longer than 20 characters it will be cutted.
// !WARNING! The "name" argument MUST BE TRUSTED!
func (s *style) viewItem(trustedNameValue string, value any, lipglosStyle ...lipgloss.Style) string {
	var (
		maxLength = 24
		v         string
		oversize  string
		sep       = ":"
	)

	/** The reason why we not use the method: fmt.Sprintf("%v", value) first is because
	* we will have a better performance since we need to ASNI escape the string.
	* If we know that it is an int type, it can be handled faster.
	* Major of the items will be int based values
	 */
	switch val := value.(type) {
	case string:
		v = strconv.Quote(val)
	case int:
		v = strconv.Itoa(val)
	case float64, float32:
		if v = "0"; strings.Index(v, ".") != -1 {
			v = fmt.Sprintf("%3.f", val)
		}
	case nil:
		v = "-"
	default:
		v = strconv.Quote(fmt.Sprintf("%v", val))
	}

	// Value are longer than expected then cut it:
	if len(v) > maxLength {
		oversize = "..."
		v = v[0:maxLength]
	}

	if len(lipglosStyle) > 0 {
		return s.item.Render(trustedNameValue+sep) + lipglosStyle[0].Render(v) + s.oversize(oversize)
	}
	return s.item.Render(trustedNameValue+sep, v, s.oversize(oversize))
}

func (m *style) payloadView() string {
	return "test" //spin + info + gap + prog + pkgCount
}

// Update the process bar
func (s *style) progressBar(header string, current, max int) string {
	v := header + strconv.Itoa(current/max)
	if current == max {
		return s.done.Render(v)
	}
	return s.spinner.Render() + v
}
