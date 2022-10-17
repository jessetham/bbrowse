package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wrap"
)

type viewerModel struct {
	viewport   viewport.Model
	formatters []formatter
	content    string
}

func newViewerModel() viewerModel {
	return viewerModel{
		viewport:   viewport.New(0, 0),
		formatters: newFormatterList(),
	}
}

func (v viewerModel) Init() tea.Cmd {
	return nil
}

func (v viewerModel) Update(msg tea.Msg) (viewerModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.viewport.Width = msg.Width
		v.viewport.Height = msg.Height / 2

	case Pair:
		for _, formatter := range v.formatters {
			if c, ok := formatter(msg.Value); ok {
				v.content = c
				break
			}
		}

	case Bucket:
		v.content = fmt.Sprintf("# of nested buckets: %d | # of pairs: %d", len(msg.Buckets), len(msg.Pairs))

	default:
		v.viewport, cmd = v.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	v.viewport.SetContent(wrap.String(v.content, v.viewport.Width))
	return v, tea.Batch(cmds...)
}

func (v viewerModel) View() string {
	return v.viewport.View()
}
