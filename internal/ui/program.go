package ui

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"fmt"

	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/pkg/statistics"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	WINDOW_PAYLOAD        = 1
	WINDOW_TRANSFORMATION = 0
)

type ProgramModel struct {
	keys   keyMap
	quit   bool
	Cancel bool
	// Holds the terminal user-interface (UI) design that is presented in the terminal stdin
	terminalUI *TerminalUI
	// How aften to check if the currecnt terminal width changes in milliseconds
	//TerminalWidthCheckDelay time.Duration
	// Listen if the terminal width change during the running process
	//channelScreenWidth chan int
	// Listen for new result to be deisplayed in the terminal user interface
	//channelResult chan Data
	// The data that contains all data that will be displayed during the process
	data Data
	// Represent the visual progress bar
	spinner spinner.Model

	// Contains the index of the current window
	window int

	help helpmenu
}

type Data struct {
	ResultFinal output.ResultFinal
	stats       statistics.Statistic
}

func NewProgram() *tea.Program {
	return tea.NewProgram(ProgramModel{
		terminalUI: NewTerminalUI(),
		spinner:    spinner.New(spinner.WithSpinner(spinner.MiniDot)),
		//channelResult:           make(chan Data),
		//channelScreenWidth:      make(chan int),
		//TerminalWidthCheckDelay: 1, //Todo
		help: helpmenu{
			menu:  help.New(),
			style: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		},

		keys: keyMap{
			Tab: key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("«»", "switch window"),
			),
			Up: key.NewBinding(
				key.WithKeys("up"),
				key.WithHelp("↑", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down"),
				key.WithHelp("↓", "move down"),
			),
			Left: key.NewBinding(
				key.WithKeys("left"),
				key.WithHelp("←", "move left"),
			),
			Right: key.NewBinding(
				key.WithKeys("right"),
				key.WithHelp("→", "move right"),
			),
			Exit: key.NewBinding(
				key.WithKeys("ctrl+c"),
				key.WithHelp("ctrl+c", "exit"),
			),
		},
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
func (m ProgramModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Tab):
			m.changeWindow()
		case key.Matches(msg, m.keys.Up):
			fmt.Println("↑")
		case key.Matches(msg, m.keys.Down):
			fmt.Println("↓")
		case key.Matches(msg, m.keys.Left):
			fmt.Println("←")
		case key.Matches(msg, m.keys.Right):
			fmt.Println("→")

		case key.Matches(msg, m.keys.Exit):
			m.quit = true
			return m, tea.Quit
		}
		return m, nil

	case statistics.Statistic:
		m.data.stats = msg
		return m, func() tea.Msg { return m.data }

	case output.ResultFinal:
		m.data.ResultFinal = msg
		return m, func() tea.Msg { return m.data }

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// Updated the spinner (loadingbar)
	default:
		return m, nil
	}
}

// Output the terminal user interface (UI) to the terminal
func (m ProgramModel) View() string {
	var view string

	m.prepareTerminalUI()
	view = m.terminalUI.Render(m.data)

	/* model := m.currentFocusedModel()
	if m.state == timerView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), modelStyle.Render(m.spinner.View()))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), focusedModelStyle.Render(m.spinner.View()))
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	*/
	view += "\n" + m.help.menu.View(m.keys)

	return view
}

func (m ProgramModel) prepareTerminalUI() {
	m.terminalUI.SetWindow(m.window)
	m.terminalUI.SetSpinner(m.spinner.View())
}

func (m ProgramModel) changeWindow() {
	if m.window == WINDOW_PAYLOAD {
		m.window = WINDOW_TRANSFORMATION
	} else {
		m.window = WINDOW_PAYLOAD
	}
}
