package sections

import (
	"context"
	"strings"
	"teachat/pkgs/llmclients"
	"teachat/pkgs/llminterface"
	"teachat/pkgs/styles"
	"teachat/pkgs/teamsg"
	"teachat/pkgs/types"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Convo struct {
	hidden     bool
	focused    bool
	messages   []string
	viewport   viewport.Model
	style      lipgloss.Style
	chatClient llminterface.Client
}

func NewConvo() Section {

	vp := viewport.New(0, 0)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	convo := &Convo{
		viewport: vp,
		style:    styles.ActiveStyle.Copy(),
	}

	return convo
}

func (c *Convo) GetSectionName() SectionName {
	return ConvoSection
}

func (c *Convo) SetDimensions(width, height int) {
	c.viewport.Width = width
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
	case teamsg.ChatPromptMsg:
		prompt := string(msg)
		c.messages = append(c.messages, styles.SenderStyle.Render("\nYou: ")+prompt, styles.AiStyle.Render("\nAI: "))
		c.viewport.SetContent(strings.Join(c.messages, "\n"))
		c.viewport.GotoBottom()
		return c, func() tea.Msg { return c.chat(string(msg)) }
	case teamsg.ChatStreamMsg:
		cs := types.ChatStream{Stream: msg}
		return c, func() tea.Msg { return c.receiveChatStream(cs) }
	case teamsg.ChatStreamDeltaMsg:
		curMessage := c.messages[len(c.messages)-1]
		curMessage = curMessage + msg.Response.Text
		c.messages[len(c.messages)-1] = wordwrap.String(curMessage, c.viewport.Width)
		c.viewport.SetContent(strings.Join(c.messages, "\n"))
		c.viewport.GotoBottom()
		return c, func() tea.Msg { return c.receiveChatStream(types.ChatStream(msg)) }
	case teamsg.ChatStreamCloseMsg:
		msg.Stream.Close()
		return c, nil
	case teamsg.ModelSelectedMsg:
		client := llmclients.PlatformInitialization[msg.Platform](true)
		client.SetModel(msg.Name)
		c.chatClient = client
		return c, nil
	}
	return c, nil
}

func (c Convo) View() string {
	if !c.hidden {
		if c.focused {
			return styles.ActiveStyle.Render(c.viewport.View())
		}
		return styles.InactiveStyle.Render(c.viewport.View())
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
	return teamsg.ChatStreamMsg(streamreader)
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
		return teamsg.ChatStreamCloseMsg(chatStream)
	}
	return teamsg.ChatStreamDeltaMsg(chatStream)
}
