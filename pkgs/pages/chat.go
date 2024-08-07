package pages

import (
	"teachat/pkgs/sections"
	"teachat/pkgs/styles"
	"teachat/pkgs/teamsg"

	tea "github.com/charmbracelet/bubbletea"
)

type Chat struct {
	current         bool
	name            PageName
	sections        map[sections.SectionName]sections.Section
	orderedSections []sections.SectionName
}

func NewChatPage() PageInterface {
	p := &Chat{}
	p.name = ChatPage
	p.AddSection(sections.NewPrompt())
	p.AddSection(sections.NewConvo())
	p.switchSection()
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

func (p *Chat) AddSection(section sections.Section) {
	if p.sections == nil {
		p.sections = make(map[sections.SectionName]sections.Section)
	}
	if len(p.orderedSections) > 0 {
		for _, sec := range p.orderedSections {
			p.sections[sec].Blur()
		}
	}
	section.SetDimensions(0, styles.Height)
	section.Show()
	section.Focus()
	p.orderedSections = append(p.orderedSections, section.GetSectionName())
	p.sections[section.GetSectionName()] = section
}

func (p *Chat) Update(msg tea.Msg) (PageInterface, tea.Cmd) {
	var cmds []tea.Cmd
	if !p.current {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg, teamsg.ModelSelectedMsg, teamsg.GetSupportedModelsMsg:
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyDown, tea.KeyEnd:
			sec, cmd := p.sections[sections.ConvoSection].Update(msg)
			p.sections[sections.ConvoSection] = sec
			return p, cmd
		case tea.KeyTab:
			p.switchSection()
			return p, nil
		}
	}
	// update all sections
	for i, s := range p.sections {
		var cmd tea.Cmd
		s, cmd = s.Update(msg)
		cmds = append(cmds, cmd)
		p.sections[i] = s
	}
	return p, tea.Batch(cmds...)
}

func (p *Chat) View() string {
	var view string
	for _, section := range p.orderedSections {
		if !p.sections[section].IsHidden() {
			view = attachView(view, p.sections[section].View())
		}
	}
	return view
}

func (p *Chat) SetDimensions(width, height int) {
	p.sections[sections.ConvoSection].SetDimensions(int(float64(width)*0.7), height)
	p.sections[sections.PromptSection].SetDimensions(int(float64(width)*0.2), height)
}

func (p *Chat) switchSection() {
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
