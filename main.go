package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	"teachat/pkgs/llminterface"
	"teachat/pkgs/openai"
	"teachat/pkgs/utils"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	textarea    textarea.Model
	senderStyle lipgloss.Style
	aiStyle     lipgloss.Style
	chatClient  llminterface.Client
	messages    []string
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(0, 0)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	chatClient := openai.New(llminterface.GPT35, true)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		aiStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		chatClient:  chatClient,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds  []tea.Cmd
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		height := int(float64(msg.Height) * 0.9)
		m.viewport.Width = int(float64(msg.Width) * 0.7)
		m.textarea.SetWidth(int(float64(msg.Width) * 0.2))
		m.viewport.Height = height
		m.textarea.SetHeight(height)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp, tea.KeyDown, tea.KeyEnd:
			m.viewport, vpCmd = m.viewport.Update(msg)
			return m, vpCmd
		case tea.KeyEnter:
			prompt := m.textarea.Value()
			chatPromptCmd := func() tea.Msg { return chatPromptMsg(prompt) }
			cmds = append(cmds, chatPromptCmd)
			m.messages = append(m.messages, m.senderStyle.Render("\nYou: ")+prompt)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			m.viewport, vpCmd = m.viewport.Update(msg)
			cmds = append(cmds, vpCmd)
			return m, tea.Batch(cmds...)
		}
	case chatStreamMsg:
		cs := chatStream{stream: msg}
		return m, func() tea.Msg { return m.receiveChatStream(cs) }
	case chatStreamDeltaMsg:
		curMessage := m.messages[len(m.messages)-1]
		curMessage = fmt.Sprintf("%s%s", curMessage, msg.response.Text)
		m.messages[len(m.messages)-1] = wordwrap.String(curMessage, m.viewport.Width)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		return m, func() tea.Msg { return m.receiveChatStream(chatStream(msg)) }
	case chatStreamCloseMsg:
		msg.stream.Close()
		return m, nil
	case chatPromptMsg:
		utils.LogToFile("chat.log", "info", "received chatPromptMsg with prompt "+string(msg))
		m.messages = append(m.messages, m.aiStyle.Render("\nAI: "))
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.viewport.GotoBottom()
		cmds = append(cmds, func() tea.Msg { return m.chat(string(msg)) })
		return m, tea.Batch(cmds...)

	}

	m.textarea, tiCmd = m.textarea.Update(msg)
	cmds = append(cmds, vpCmd, tiCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Left, m.textarea.View(), m.viewport.View())
}
