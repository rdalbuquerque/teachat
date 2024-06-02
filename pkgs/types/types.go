package types

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model string

// implement list.Item interface
func (m Model) FilterValue() string { return "" }

const (
	Llama3 Model = "llama3"
	GPT35  Model = "gpt-3.5-turbo-0125"
	GPT4   Model = "gpt-4"
	GPT4o  Model = "gpt-4o"
)

type StreamReader interface {
	Bytes() []byte
	Scan() bool
	Close()
}

type ChatResponse struct {
	Done bool
	Text string
}

type ChatStream struct {
	Response *ChatResponse
	Stream   StreamReader
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type ModelItemDelegate struct{}

func (d ModelItemDelegate) Height() int                             { return 1 }
func (d ModelItemDelegate) Spacing() int                            { return 0 }
func (d ModelItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d ModelItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Model)
	if !ok {
		return
	}
	modelStr := string(i)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("| " + modelStr)
		}
	}

	fmt.Fprint(w, fn(modelStr))
}
