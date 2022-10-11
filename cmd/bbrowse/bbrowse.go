package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	filename  string
	stack     []Bucket
	err       error
	list      list.Model
	stackSize int
}

func newModel(filename string) model {
	return model{
		list:     list.New([]list.Item{Pair{Key: []byte("test"), Value: []byte("test")}}, list.NewDefaultDelegate(), 50, 25),
		filename: filename,
	}
}

func (m model) Init() tea.Cmd {
	return openAndReadBoltDB(m.filename)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case Bucket:
		m.stack = append(m.stack, msg)

	case tea.KeyMsg:
		switch {
		case msg.String() == "l":
			if b, ok := m.list.SelectedItem().(Bucket); ok {
				m.stack = append(m.stack, b)
			}
		case msg.String() == "h":
			if len(m.stack) > 1 {
				m.stack = m.stack[:len(m.stack)-1]
			}
		}

	case error:
		m.err = msg
	}

	if m.stackSize != len(m.stack) {
		var title string
		for _, b := range m.stack {
			title += string(b.Name) + "/"
		}
		m.list.Title = title
		
		m.stackSize = len(m.stack)
		newItems := []list.Item{}
		bucket := m.stack[len(m.stack)-1]
		for _, b := range bucket.Buckets {
			newItems = append(newItems, Bucket(*b))
		}
		for _, p := range bucket.Pairs {
			newItems = append(newItems, Pair(*p))
		}
		cmd := m.list.SetItems(newItems)

		return m, cmd
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nThere's been an error: %v\n\n", m.err)
	} else if len(m.stack) == 0 {
		return "loading db"
	}
	return m.list.View()
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
