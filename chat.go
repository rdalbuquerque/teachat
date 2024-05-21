package main

import (
	"context"

	"teachat/pkgs/llminterface"

	tea "github.com/charmbracelet/bubbletea"
)

type chatStream struct {
	response *llminterface.ChatResponse
	stream   llminterface.StreamReader
}

type chatPromptMsg string
type chatStreamMsg llminterface.StreamReader
type chatStreamDeltaMsg chatStream
type chatStreamCloseMsg chatStream

func (m model) receiveChatStream(stream chatStream) tea.Msg {
	resp, respstream, err := m.chatClient.GetDelta(context.Background(), stream.stream)
	if err != nil {
		panic(err)
	}
	chatStream := chatStream{
		response: resp,
		stream:   respstream,
	}
	if resp.Done {
		return chatStreamCloseMsg(chatStream)
	}
	return chatStreamDeltaMsg(chatStream)
}

func (m model) chat(msg string) tea.Msg {
	streamreader, err := m.chatClient.Prompt(context.Background(), msg)
	if err != nil {
		panic(err)
	}
	return chatStreamMsg(streamreader)
}
