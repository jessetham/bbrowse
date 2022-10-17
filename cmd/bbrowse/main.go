package main

import (
	"flag"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var styles = lipgloss.NewStyle()

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
		if msg.String() == "tab" {
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

	return lipgloss.JoinVertical(lipgloss.Left, styles.Render(m.viewer.View()), styles.Render(m.nav.View()))
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
