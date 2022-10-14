package main

import (
	bbrowseGob "bbrowse/gob"
	"bytes"
	"encoding/hex"
)

type formatter func([]byte) (string, bool)

func newFormatterList() []formatter {
	return []formatter{
		gobFormatter,
		// Fallback if none of the other formatters work.
		hexStringFormatter,
	}
}

func gobFormatter(b []byte) (string, bool) {
	var buf bytes.Buffer
	if err := bbrowseGob.Debug(bytes.NewReader(b), &buf); err != nil {
		return "", false
	}
	return buf.String(), true
}

func hexStringFormatter(b []byte) (string, bool) {
	return hex.EncodeToString(b), true
}
