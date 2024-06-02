package sections

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Section interface {
	IsHidden() bool
	IsFocused() bool
	Hide()
	Show()
	Focus()
	Blur()
	View() string
	Update(msg tea.Msg) (Section, tea.Cmd)
	SetDimensions(width, height int)
}

type SectionName string

const (
	HelpSection      SectionName = "help"
	PromptSection    SectionName = "prompt"
	ConvoSection     SectionName = "convo"
	ModelListSection SectionName = "modellist"
)
