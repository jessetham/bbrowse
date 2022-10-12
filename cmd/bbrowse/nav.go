package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type navModel struct {
	list     list.Model
	keys     *navKeyMap
	styles   lipgloss.Style
	stack    []Bucket
	stackLen int
}

func newNavModel(styles lipgloss.Style) navModel {
	navKeys := newNavKeyMap()

	navDelegate := list.NewDefaultDelegate()
	navKeyBindings := []key.Binding{navKeys.forward, navKeys.back}
	navDelegate.ShortHelpFunc = func() []key.Binding { return navKeyBindings }
	navDelegate.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{navKeyBindings} }

	navList := list.New([]list.Item{}, navDelegate, 0, 0)

	return navModel{
		list:   navList,
		keys:   navKeys,
		styles: styles,
		stack:  []Bucket{},
	}
}

func (n navModel) Init() tea.Cmd {
	return nil
}

func (n navModel) Update(msg tea.Msg) (navModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		n.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, n.keys.forward):
			if bucket, ok := n.list.SelectedItem().(Bucket); ok {
				n.stack = append(n.stack, bucket)
			}
		case key.Matches(msg, n.keys.back):
			if len(n.stack) > 1 {
				n.stack = n.stack[:len(n.stack)-1]
			}
		}

	case Bucket:
		n.stack = append(n.stack, msg)
	}

	// An item has been pushed or popped off the stack and we need to update our list to reflect that
	if n.stackLen != len(n.stack) {
		var path []string
		for _, b := range n.stack {
			path = append(path, string(b.Name))
		}
		n.list.Title = strings.Join(path, "/")

		newItems := []list.Item{}
		displayedBucket := n.stack[len(n.stack)-1]
		for _, b := range displayedBucket.Buckets {
			newItems = append(newItems, Bucket(*b))
		}
		for _, p := range displayedBucket.Pairs {
			newItems = append(newItems, Pair(*p))
		}
		cmd := n.list.SetItems(newItems)

		n.stackLen = len(n.stack)

		return n, cmd
	}

	updatedModel, cmd := n.list.Update(msg)
	n.list = updatedModel

	return n, cmd
}

func (n navModel) View() string {
	return n.styles.Render(n.list.View())
}

type navKeyMap struct {
	forward key.Binding
	back    key.Binding
}

func newNavKeyMap() *navKeyMap {
	return &navKeyMap{
		forward: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "details"),
		),
		back: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "back"),
		),
	}
}
