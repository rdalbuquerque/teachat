package main

import (
	"context"
	"fmt"

	"teachat/pkgs/pages"
	"teachat/pkgs/styles"
	"teachat/pkgs/teamsg"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	ctx       context.Context
	cancel    context.CancelFunc
	pages     map[pages.PageName]pages.PageInterface
	pageStack pages.Stack
	height    int
	width     int
}

func initialModel() model {
	ctx, cancel := context.WithCancel(context.Background())
	helpPage := pages.NewHelpPage()
	chatPage := pages.NewChatPage()
	modelSelectionPage := pages.NewModelSelectionPage()
	pagesMap := map[pages.PageName]pages.PageInterface{
		pages.ModelSelectionPage: modelSelectionPage,
		pages.ChatPage:           chatPage,
		pages.HelpPage:           helpPage,
	}
	pageStack := pages.Stack{}
	m := model{
		ctx:       ctx,
		cancel:    cancel,
		pages:     pagesMap,
		pageStack: pageStack,
	}
	m.addPage(pages.ModelSelectionPage)
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancel()
			return m, tea.Quit
		case "ctrl+h":
			if m.pageStack.Peek().GetPageName() != pages.HelpPage {
				m.addPage(pages.HelpPage)
			}
			return m, nil
		case "ctrl+b":
			m.removeCurrentPage()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		for _, p := range m.pages {
			styles.SetDimensions(m.width, msg.Height-3)
			p.SetDimensions(m.width, msg.Height-3)
		}
		return m, nil
	case teamsg.ModelSelectedMsg:
		m.addPage(pages.ChatPage)
	}
	// update all pages
	updatedPages := make(map[pages.PageName]pages.PageInterface)
	var cmds []tea.Cmd
	for _, p := range m.pages {
		updatedPage, cmd := p.Update(msg)
		updatedPages[updatedPage.GetPageName()] = updatedPage
		cmds = append(cmds, cmd)
	}
	m.pages = updatedPages
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return m.pageStack.Peek().View()
}

func (m *model) addPage(pageName pages.PageName) {
	if len(m.pageStack) > 0 {
		m.pageStack.Peek().UnsetCurrentPage()
	}
	p := m.pages[pageName]
	if p == nil {
		panic(fmt.Sprintf("page %s not found", pageName))
	}
	p.SetAsCurrentPage()
	m.pageStack.Push(p)
}

func (m *model) removeCurrentPage() {
	m.pageStack.Peek().UnsetCurrentPage()
	m.pageStack.Pop()
	m.pageStack.Peek().SetAsCurrentPage()
}

func main() {
	initialModel := initialModel()
	if _, err := tea.NewProgram(&initialModel).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}
}
