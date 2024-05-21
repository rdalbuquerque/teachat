package pages

import (
	"context"
	"teachat/pkgs/sections"
	"teachat/pkgs/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type Help struct {
	current  bool
	name     PageName
	sections map[sections.SectionName]sections.Section
}

func NewHelpPage() PageInterface {
	p := &Help{}
	p.name = HelpPage
	p.AddSection(context.Background(), sections.HelpSection)
	return p
}

func (p *Help) IsCurrentPage() bool {
	return p.current
}

func (p *Help) SetAsCurrentPage() {
	p.current = true
}

func (p *Help) UnsetCurrentPage() {
	p.current = false
}

func (p *Help) GetPageName() PageName {
	return p.name
}

func (p *Help) AddSection(ctx context.Context, section sections.SectionName) {
	if p.sections == nil {
		p.sections = make(map[sections.SectionName]sections.Section)
	}
	newSection := sectionNewFuncs[section](ctx)
	newSection.SetDimensions(0, styles.Height)
	newSection.Show()
	newSection.Focus()
	p.sections[section] = newSection
}

func (p *Help) View() string {
	return p.sections[sections.HelpSection].View()
}

func (p *Help) Update(msg tea.Msg) (PageInterface, tea.Cmd) {
	if p.current {
		sec, cmd := p.sections[sections.HelpSection].Update(msg)
		p.sections[sections.HelpSection] = sec
		return p, cmd
	}
	return p, nil
}

func (p *Help) SetDimensions(width, height int) {
	p.sections[sections.HelpSection].SetDimensions(width, height)
}
