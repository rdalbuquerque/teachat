package sections

import (
	"teachat/pkgs/styles"
	"teachat/pkgs/teamsg"
	"teachat/pkgs/types"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelList struct {
	hidden  bool
	focused bool
	list    list.Model
}

func NewModelList() Section {
	list := list.New([]list.Item{}, types.ModelItemDelegate{}, 0, 0)

	return &ModelList{
		list: list,
	}
}

func (s *ModelList) GetSectionName() SectionName {
	return ModelListSection
}

func (s *ModelList) SetDimensions(width, height int) {
	s.list.SetWidth(width)
	s.list.SetHeight(height)
}

func (s *ModelList) IsHidden() bool {
	return s.hidden
}

func (s *ModelList) IsFocused() bool {
	return s.focused
}

func (s *ModelList) Update(msg tea.Msg) (Section, tea.Cmd) {
	if s.focused {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				selectedModel := s.list.SelectedItem().(types.Model)
				return s, func() tea.Msg { return teamsg.ModelSelectedMsg(selectedModel) }
			}
			vp, cmd := s.list.Update(msg)
			s.list = vp
			return s, cmd
		case teamsg.ModelsMsg:
			items := make([]list.Item, len(msg))
			for i := range msg {
				items[i] = msg[i]
			}
			return s, s.list.SetItems(items)
		case teamsg.GetSupportedModelsMsg:
			models := []types.Model{
				{Name: types.GPT35, Platform: types.OpenAI},
				{Name: types.GPT4, Platform: types.OpenAI},
				{Name: types.GPT4o, Platform: types.OpenAI},
				{Name: types.Llama3, Platform: types.Ollama},
			}
			return s, func() tea.Msg { return teamsg.ModelsMsg(models) }
		}
	}
	return s, nil
}

func (s *ModelList) View() string {
	if !s.hidden {
		if s.focused {
			return styles.ActiveStyle.Render(s.list.View())
		}
		return styles.InactiveStyle.Render(s.list.View())
	}
	return ""
}

func (s *ModelList) Hide() {
	s.hidden = true
}

func (s *ModelList) Show() {
	s.hidden = false
}

func (s *ModelList) Focus() {
	s.Show()
	s.focused = true
}

func (s *ModelList) Blur() {
	s.focused = false
}
