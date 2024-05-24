package sections

import (
	"context"
	"teachat/pkgs/styles"
	"teachat/pkgs/teamsgs"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Prompt struct {
	hidden   bool
	focused  bool
	textarea textarea.Model
	style    lipgloss.Style
}

func NewPrompt(_ context.Context) Section {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return &Prompt{
		textarea: ta,
		style:    styles.ActiveStyle.Copy(),
	}
}

func (p *Prompt) SetDimensions(width, height int) {
	p.textarea.SetWidth(width)
	p.textarea.SetHeight(height)
}

func (p *Prompt) IsHidden() bool {
	return p.hidden
}

func (p *Prompt) IsFocused() bool {
	return p.focused
}

func (p *Prompt) Update(msg tea.Msg) (Section, tea.Cmd) {
	if p.focused {

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				prompt := p.textarea.Value()
				p.textarea.Reset()

				return p, func() tea.Msg { return teamsgs.ChatPromptMsg(prompt) }
			}
		}

		vp, cmd := p.textarea.Update(msg)
		p.textarea = vp
		return p, cmd
	}
	return p, nil
}

func (p *Prompt) View() string {
	if !p.hidden {
		if p.focused {
			return styles.ActiveStyle.Render(p.textarea.View())
		}
		return styles.InactiveStyle.Render(p.textarea.View())
	}
	return ""
}

func (p *Prompt) Hide() {
	p.hidden = true
}

func (p *Prompt) Show() {
	p.hidden = false
}

func (p *Prompt) Focus() {
	p.Show()
	p.focused = true
}

func (p *Prompt) Blur() {
	p.focused = false
}
