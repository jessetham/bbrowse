package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wrap"

	"golang.design/x/clipboard"
)

type model struct {
	filename               string
	list                   list.Model
	keys                   *keybindingMap
	stack                  []Bucket
	stackLen               int
	viewport               viewport.Model
	currentViewportContent string
	viewportFocused        bool
	err                    error
}

func newModel(filename string) model {
	keybindingsMap := newKeybindingMap()

	delegate := list.NewDefaultDelegate()
	keybindings := []key.Binding{keybindingsMap.forward, keybindingsMap.back, keybindingsMap.toggleFocus, keybindingsMap.copy}
	delegate.ShortHelpFunc = func() []key.Binding { return keybindings }
	delegate.FullHelpFunc = func() [][]key.Binding { return [][]key.Binding{keybindings} }
	delegate.ShowDescription = false

	l := list.New([]list.Item{}, delegate, 0, 0)
	// We want to use left and right to navigate buckets and not change list pages.
	l.KeyMap.NextPage.Unbind()
	l.KeyMap.PrevPage.Unbind()

	vp := viewport.New(0, 0)

	return model{
		filename: filename,
		list:     l,
		keys:     keybindingsMap,
		stack:    []Bucket{},
		viewport: vp,
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
	case error:
		m.err = msg

	case tea.WindowSizeMsg:
		width := msg.Width / 2
		height := msg.Height
		m.list.SetSize(width, height)
		m.viewport.Height = height
		m.viewport.Width = width

	// Handle custom keybindings that aren't implemented in UI components
	case tea.KeyMsg:
		// Don't navigate when filtering list items
		if m.list.SettingFilter() {
			break
		}
		switch {
		case key.Matches(msg, m.keys.forward) && !m.viewportFocused:
			if bucket, ok := m.list.SelectedItem().(Bucket); ok {
				m.stack = append(m.stack, bucket)
			}
		case key.Matches(msg, m.keys.back) && !m.viewportFocused:
			if len(m.stack) > 1 {
				m.stack = m.stack[:len(m.stack)-1]
			}
		case key.Matches(msg, m.keys.copy):
			clipboard.Write(clipboard.FmtText, []byte(m.currentViewportContent))
			cmds = append(cmds, m.list.NewStatusMessage("Copied!"))

		case key.Matches(msg, m.keys.toggleFocus):
			m.viewportFocused = !m.viewportFocused
		}

	case initialBucket:
		m.stack = append(m.stack, Bucket(msg))
	}

	// An item has been pushed or popped off the stack and we need to update our list to reflect that
	if m.stackLen != len(m.stack) {
		var path []string
		for _, b := range m.stack {
			path = append(path, string(b.Name))
		}
		m.list.Title = strings.Join(path, "/")

		newItems := []list.Item{}
		displayedBucket := m.stack[len(m.stack)-1]
		for _, b := range displayedBucket.Buckets {
			newItems = append(newItems, Bucket(*b))
		}
		for _, p := range displayedBucket.Pairs {
			newItems = append(newItems, Pair(*p))
		}
		cmd = m.list.SetItems(newItems)
		cmds = append(cmds, cmd)

		m.stackLen = len(m.stack)
	}

	if m.viewportFocused {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.currentViewportContent = getViewportContent(m.list.SelectedItem())
	m.viewport.SetContent(wrap.String(m.currentViewportContent, m.viewport.Width))

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("got an error: %s", m.err)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), m.viewport.View())
}

func getViewportContent(item list.Item) string {
	var content string
	switch item := item.(type) {
	case Pair:
		for _, formatter := range formatters {
			if c, ok := formatter(item.Value); ok {
				content = c
				break
			}
		}

	case Bucket:
		content = fmt.Sprintf("# of nested buckets: %d | # of pairs: %d", len(item.Buckets), len(item.Pairs))

	default:
		content = "got an unknown item type"
	}

	return content
}

type keybindingMap struct {
	forward     key.Binding
	back        key.Binding
	toggleFocus key.Binding
	copy        key.Binding
}

func newKeybindingMap() *keybindingMap {
	return &keybindingMap{
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
		copy: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "copy current content"),
		),
	}
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		log.Fatal("no filename given")
	}

	if err := clipboard.Init(); err != nil {
		log.Fatal(err)
	}

	if err := tea.NewProgram(newModel(filename), tea.WithAltScreen()).Start(); err != nil {
		log.Fatal(err)
	}
}
