package pages

import (
	"context"
	"teachat/pkgs/sections"
	"teachat/pkgs/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type Chat struct {
	current  bool
	name     PageName
	sections map[sections.SectionName]sections.Section
}

func NewChatPage() PageInterface {
	p := &Chat{}
	p.name = ChatPage
	p.AddSection(context.Background(), sections.PromptSection)
	p.AddSection(context.Background(), sections.ConvoSection)
	return p
}

func (p *Chat) IsCurrentPage() bool {
	return p.current
}

func (p *Chat) SetAsCurrentPage() {
	p.current = true
}

func (p *Chat) UnsetCurrentPage() {
	p.current = false
}

func (p *Chat) GetPageName() PageName {
	return p.name
}

func (p *Chat) AddSection(ctx context.Context, section sections.SectionName) {
	if p.sections == nil {
		p.sections = make(map[sections.SectionName]sections.Section)
	}
	newSection := sectionNewFuncs[section](ctx)
	newSection.SetDimensions(0, styles.Height)
	newSection.Show()
	newSection.Focus()
	p.sections[section] = newSection
}

func (p *Chat) View() string {
	return p.sections[sections.HelpSection].View()
}

func (p *Chat) Update(msg tea.Msg) (PageInterface, tea.Cmd) {
	if p.current {
		sec, cmd := p.sections[sections.HelpSection].Update(msg)
		p.sections[sections.HelpSection] = sec
		return p, cmd
	}
	return p, nil
}

func (p *Chat) SetDimensions(width, height int) {
	for s := range p.sections {
		p.sections[s].SetDimensions(width, height)
	}
}
