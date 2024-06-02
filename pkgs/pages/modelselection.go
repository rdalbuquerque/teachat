package pages

import (
	"context"
	"teachat/pkgs/sections"
	"teachat/pkgs/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type ModelSelection struct {
	current         bool
	name            PageName
	sections        map[sections.SectionName]sections.Section
	orderedSections []sections.SectionName
}

func NewModelSelectionPage() PageInterface {
	p := &ModelSelection{}
	p.name = ModelSelectionPage
	p.AddSection(context.Background(), sections.ModelListSection)
	return p
}

func (p *ModelSelection) IsCurrentPage() bool {
	return p.current
}

func (p *ModelSelection) SetAsCurrentPage() {
	p.current = true
}

func (p *ModelSelection) UnsetCurrentPage() {
	p.current = false
}

func (p *ModelSelection) GetPageName() PageName {
	return p.name
}

func (p *ModelSelection) AddSection(ctx context.Context, section sections.SectionName) {
	if p.sections == nil {
		p.sections = make(map[sections.SectionName]sections.Section)
	}
	if len(p.orderedSections) > 0 {
		for _, sec := range p.orderedSections {
			p.sections[sec].Blur()
		}
	}
	newSection := sectionNewFuncs[section](ctx)
	newSection.SetDimensions(0, styles.Height)
	newSection.Show()
	newSection.Focus()
	p.orderedSections = append(p.orderedSections, section)
	p.sections[section] = newSection
}

func (p *ModelSelection) Update(msg tea.Msg) (PageInterface, tea.Cmd) {
	var cmds []tea.Cmd
	if p.current {
		// update all sections
		for i, s := range p.sections {
			var cmd tea.Cmd
			s, cmd = s.Update(msg)
			cmds = append(cmds, cmd)
			p.sections[i] = s
		}
		return p, tea.Batch(cmds...)
	}
	return p, nil
}

func (p *ModelSelection) View() string {
	var view string
	for _, section := range p.orderedSections {
		if !p.sections[section].IsHidden() {
			view = attachView(view, p.sections[section].View())
		}
	}
	return view
}

func (p *ModelSelection) SetDimensions(width, height int) {
	p.sections[sections.ModelListSection].SetDimensions(width, height)
}

func (p *ModelSelection) switchSection() {
	shownSections := []sections.SectionName{}
	for _, section := range p.orderedSections {
		if !p.sections[section].IsHidden() {
			shownSections = append(shownSections, section)
		}
	}
	for i, sec := range shownSections {
		section := p.sections[sec]
		if section.IsFocused() {
			section.Blur()
			nextKey := shownSections[0] // default to the first key
			if i+1 < len(shownSections) {
				nextKey = shownSections[i+1] // if there's a next key, use it
			}
			p.sections[nextKey].Focus()
			return
		}
	}
}
