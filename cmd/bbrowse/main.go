package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var styles = lipgloss.NewStyle()

type model struct {
	filename string
	nav      navModel
	viewer   viewerModel
	err      error
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		updatedNavModel, cmd := m.nav.Update(msg)
		m.nav = updatedNavModel
		cmds = append(cmds, cmd)

		updatedViewerModel, cmd := m.viewer.Update(updatedNavModel.list.SelectedItem())
		m.viewer = updatedViewerModel
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case error:
		m.err = msg
	}

	updatedNavModel, cmd := m.nav.Update(msg)
	m.nav = updatedNavModel
	cmds = append(cmds, cmd)

	updatedViewerModel, cmd := m.viewer.Update(msg)
	m.viewer = updatedViewerModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nThere's been an error: %v\n\n", m.err)
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, styles.Render(m.nav.View()), styles.Render(m.viewer.View()))
}

func main() {
	filename := "572009747.db"

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := tea.NewProgram(newModel(filename)).Start(); err != nil {
		log.Fatal(err)
	}
}
