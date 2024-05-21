package sections

import (
	"context"
	"strings"
	"teachat/pkgs/llminterface"
	"teachat/pkgs/openai"
	"teachat/pkgs/styles"
	"teachat/pkgs/teamsgs"
	"teachat/pkgs/types"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Convo struct {
	hidden     bool
	focused    bool
	messages   []string
	viewport   viewport.Model
	style      lipgloss.Style
	chatClient llminterface.Client
}

func NewConvo(_ context.Context) Section {
	vp := viewport.New(0, 0)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	chatClient := openai.New(types.GPT35, true)

	return &Convo{
		chatClient: chatClient,
		viewport:   vp,
		style:      styles.ActiveStyle.Copy(),
	}
}

func (c *Convo) SetDimensions(width, height int) {
	c.viewport.Height = height
}

func (c Convo) IsHidden() bool {
	return c.hidden
}

func (c Convo) IsFocused() bool {
	return c.focused
}

func (c *Convo) Update(msg tea.Msg) (Section, tea.Cmd) {
	if c.focused {
		vp, cmd := c.viewport.Update(msg)
		c.viewport = vp
		return c, cmd
	}
	switch msg := msg.(type) {
	case teamsgs.ChatPromptMsg:
		prompt := string(msg)
		c.messages = append(c.messages, styles.SenderStyle.Render("\nYou: ")+prompt, styles.AiStyle.Render("\nAI: "))
		c.viewport.SetContent(strings.Join(c.messages, "\n"))
		c.viewport.GotoBottom()
		return c, func() tea.Msg { return c.chat(string(msg)) }
	case teamsgs.ChatStreamMsg:
		cs := types.ChatStream{Stream: msg}
		return c, func() tea.Msg { return c.receiveChatStream(cs) }
	}
	return c, nil
}

func (c Convo) View() string {
	if c.focused {
		return c.style.Width(styles.Width).Render(c.viewport.View())
	}
	return ""
}

func (c *Convo) Hide() {
	c.hidden = true
}

func (c *Convo) Show() {
	c.hidden = false
}

func (c *Convo) Focus() {
	c.Show()
	c.focused = true
}

func (c *Convo) Blur() {
	c.focused = false
}

func (c Convo) chat(msg string) tea.Msg {
	streamreader, err := c.chatClient.Prompt(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	return teamsgs.ChatStreamMsg(streamreader)
}

func (c Convo) receiveChatStream(stream types.ChatStream) tea.Msg {
	resp, respstream, err := c.chatClient.GetDelta(context.Background(), stream.Stream)
	if err != nil {
		panic(err)
	}
	chatStream := types.ChatStream{
		Response: resp,
		Stream:   respstream,
	}
	if resp.Done {
		return teamsgs.ChatStreamCloseMsg(chatStream)
	}
	return teamsgs.ChatStreamDeltaMsg(chatStream)
}
