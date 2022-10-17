package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type navModel struct {
	list     list.Model
	keys     *navKeyMap
	stack    []Bucket
	stackLen int
}

func newNavModel() navModel {
	navKeys := newNavKeyMap()

	navDelegate := list.NewDefaultDelegate()
	navKeyBindings := []key.Binding{navKeys.forward, navKeys.back, navKeys.toggleFocus}
	navDelegate.ShortHelpFunc = func() []key.Binding { return navKeyBindings }
	navDelegate.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{navKeyBindings} }
	navDelegate.ShowDescription = false

	navList := list.New([]list.Item{}, navDelegate, 0, 0)

	return navModel{
		list:  navList,
		keys:  navKeys,
		stack: []Bucket{},
	}
}

func (n navModel) Init() tea.Cmd {
	return nil
}

func (n navModel) Update(msg tea.Msg) (navModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		n.list.SetSize(msg.Width-2, msg.Height/2-2)

	case tea.KeyMsg:
		// Don't navigate when filtering list items
		if n.list.SettingFilter() {
			break
		}
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
		cmd = n.list.SetItems(newItems)
		cmds = append(cmds, cmd)

		n.stackLen = len(n.stack)
	} else {
		n.list, cmd = n.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return n, tea.Batch(cmds...)
}

func (n navModel) View() string {
	return n.list.View()
}

func (n navModel) Height() int {
	return n.list.Height()
}

func (n navModel) Width() int {
	return n.list.Width()
}

type navKeyMap struct {
	forward     key.Binding
	back        key.Binding
	toggleFocus key.Binding
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
		toggleFocus: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle focus"),
		),
	}
}
