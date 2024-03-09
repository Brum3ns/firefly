package ui

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"time"

	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/pkg/statistics"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Result struct {
	ResultFinal output.ResultFinal
	stats       statistics.Statistic
}

type ProgramModel struct {
	quit   bool
	Cancel bool
	// Holds the terminal user-interface (UI) design that is presented in the terminal stdin
	terminalUI TerminalUI
	// How aften to check if the currecnt terminal width changes in milliseconds
	TerminalWidthCheckDelay time.Duration
	// Listen if the terminal width change during the running process
	channelScreenWidth chan int
	// Listen for new result to be deisplayed in the terminal user interface
	channelResult chan Result
	// The result that contains all data that will be displayed during the process
	result Result
	// Represent the visual progress bar
	spinner spinner.Model

	list List
}

type List struct {
	list   list.Model
	choice string
	quit   bool
}

func NewProgram() *tea.Program {
	return tea.NewProgram(ProgramModel{
		spinner:                 spinner.New(spinner.WithSpinner(spinner.MiniDot)),
		channelResult:           make(chan Result),
		channelScreenWidth:      make(chan int),
		TerminalWidthCheckDelay: 1, //Todo
		terminalUI:              NewTerminalUI(),
	})
}

// Processes that are running in the background during the program core process
func (m ProgramModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}

// Listen for changes during by intercepting commands from other processes.
// Then change the needed data (Ex: the result of the runner)
func (m ProgramModel) Update(data tea.Msg) (tea.Model, tea.Cmd) {
	switch d := data.(type) {
	case tea.KeyMsg:
		if d.Type == tea.KeyCtrlC {
			m.quit = true
			return m, tea.Quit
		}
		return m, nil

	case statistics.Statistic:
		m.result.stats = d
		return m, func() tea.Msg { return m.result }

	case output.ResultFinal:
		m.result.ResultFinal = d
		return m, func() tea.Msg { return m.result }

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(d)
		return m, cmd

	// Updated the spinner (loadingbar)
	default:
		return m, nil
	}
}

// Output the terminal user interface (UI) to the terminal
func (m ProgramModel) View() string {
	m.terminalUI.spinner = m.spinner.View()
	return m.terminalUI.Render(m.result)
}
