package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type viewerModel struct {
	viewport   viewport.Model
	adapters   []adapter
	adapterIdx int
}

func newViewerModel() viewerModel {
	return viewerModel{
		viewport: viewport.New(0, 0),
		adapters: newAdapterList(),
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
		s, err := v.adapters[v.adapterIdx].toString(msg.Value)
		if err != nil {
			v.viewport.SetContent("error converting pair value")
		} else {
			v.viewport.SetContent(s)
		}

	case Bucket:
		v.viewport.SetContent("Bucket")
	}
	return v, nil
}

func (v viewerModel) View() string {
	return v.viewport.View()
}
