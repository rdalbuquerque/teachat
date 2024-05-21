package pages

import (
	"context"
	"teachat/pkgs/sections"

	tea "github.com/charmbracelet/bubbletea"
)

type SectionConstructor func(context.Context) sections.Section

var (
	sectionNewFuncs = map[sections.SectionName]SectionConstructor{
		sections.HelpSection:   sections.NewHelp,
		sections.PromptSection: sections.NewPrompt,
		sections.ConvoSection:  sections.NewConvo,
	}
)

type PageInterface interface {
	GetPageName() PageName
	AddSection(context.Context, sections.SectionName)
	SetDimensions(width, height int)
	Update(tea.Msg) (PageInterface, tea.Cmd)
	View() string
	IsCurrentPage() bool
	SetAsCurrentPage()
	UnsetCurrentPage()
}
