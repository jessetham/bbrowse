package main

import (
	"flag"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func getUnfocusedStyle(height, width int) lipgloss.Style {
	return lipgloss.NewStyle().Height(height).Width(width).BorderStyle(lipgloss.NormalBorder())
}

func getFocusedStyle(height, width int) lipgloss.Style {
	return getUnfocusedStyle(height, width).BorderForeground(lipgloss.Color("42"))
}

type model struct {
	filename      string
	nav           navModel
	viewer        viewerModel
	err           error
	viewerFocused bool
}

func newModel(filename string) model {
	return model{
		filename: filename,
		nav:      newNavModel(),
		viewer:   newViewerModel(),
	}
}

func (m model) Init() tea.Cmd {
	return openAndReadBoltDB(m.filename)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		msgString := msg.String()

		if msgString == "ctrl+c" {
			return m, tea.Quit
		}

		if msgString == "tab" {
			m.viewerFocused = !m.viewerFocused
		} else if m.viewerFocused {
			m.viewer, cmd = m.viewer.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.nav, cmd = m.nav.Update(msg)
			cmds = append(cmds, cmd)

			m.viewer, cmd = m.viewer.Update(m.nav.list.SelectedItem())
			cmds = append(cmds, cmd)
		}

	case error:
		m.err = msg

	default:
		m.nav, cmd = m.nav.Update(msg)
		cmds = append(cmds, cmd)

		m.viewer, cmd = m.viewer.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nThere's been an error: %v\n\n", m.err)
	}

	var (
		viewerRender string
		navRender    string
	)
	if m.viewerFocused {
		viewerRender = getFocusedStyle(m.viewer.Height(), m.viewer.Width()).Render(m.viewer.View())
		navRender = getUnfocusedStyle(m.nav.Height(), m.nav.Width()).Render(m.nav.View())
	} else {
		viewerRender = getUnfocusedStyle(m.viewer.Height(), m.viewer.Width()).Render(m.viewer.View())
		navRender = getFocusedStyle(m.nav.Height(), m.nav.Width()).Render(m.nav.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, viewerRender, navRender)
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		log.Fatal("no filename given")
	}

	if err := tea.NewProgram(newModel(filename), tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
