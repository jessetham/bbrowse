package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type viewerModel struct {
	viewport   viewport.Model
	formatters   []formatter
}

func newViewerModel() viewerModel {
	return viewerModel{
		viewport: viewport.New(0, 0),
		formatters: newFormatterList(),
	}
}

func (v viewerModel) Init() tea.Cmd {
	return nil
}

func (v viewerModel) Update(msg tea.Msg) (viewerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.viewport.Width = msg.Width
		v.viewport.Height = msg.Height

	case Pair:
		for _, formatter := range v.formatters {
			if c, ok := formatter(msg.Value); ok {
				v.viewport.SetContent(c)
				break
			}
		}

	case Bucket:
		v.viewport.SetContent("Bucket")
	}
	return v, nil
}

func (v viewerModel) View() string {
	return v.viewport.View()
}
