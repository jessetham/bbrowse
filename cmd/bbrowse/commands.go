package main

import (
	"bbrowse"

	tea "github.com/charmbracelet/bubbletea"
)

type initMsg Bucket

func openAndReadBoltDB(filename string) tea.Cmd {
	return func() tea.Msg {
		db, err := bbrowse.OpenAndCopyBoltDB(filename)
		if err != nil {
			return err
		}
		return Bucket(*db)
	}
}
