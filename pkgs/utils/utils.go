package utils

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Ptr[T any](input T) *T {
	return &input
}

func LogToFile(path, prefix, content string) {
	f, err := tea.LogToFile(path, prefix)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(content + "\n")
}
