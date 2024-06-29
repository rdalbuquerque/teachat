package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"teachat/pkgs/llminterface"
	"teachat/pkgs/types"
	"teachat/pkgs/utils"

	openai "github.com/sashabaranov/go-openai"
)

func (c *Client) GetSupportedModels() []string {
	resp, err := c.ListModels(context.Background())
	if err != nil {
		panic(err)
	}
	var gptModels []string
	for i := range resp.Models {
		if strings.Contains(resp.Models[i].ID, "gpt") {
			gptModels = append(gptModels, resp.Models[i].ID)
		}
	}
	return gptModels
}

func New(stream bool) llminterface.Client {
	openai_api_key := os.Getenv("OPENAI_API_KEY")
	c := openai.NewClient(openai_api_key)
	messages := make([]openai.ChatCompletionMessage, 0)
	// default type
	return &Client{
		Client:   c,
		stream:   stream,
		messages: messages,
	}
}

func (c *Client) SetModel(model string) {
	c.model = model
}

type streamReader struct {
	done     bool
	stream   *openai.ChatCompletionStream
	response openai.ChatCompletionStreamResponse
}

// The idea is to mimic the behavior of bufio.Scanner
// So, if there are no more tokens to scan, we return false
// If there is an error, we panic
func (s *streamReader) Scan() bool {
	response, err := s.stream.Recv()
	if errors.Is(err, io.EOF) {
		return false
	}
	if err != nil {
		s.stream.Close()
		panic(err)
	}
	s.response = response
	return true
}

func (s *streamReader) Close() {
	s.stream.Close()
}

func (s streamReader) Bytes() []byte {
	return []byte(s.response.Choices[0].Delta.Content)
}

type Client struct {
	messages        []openai.ChatCompletionMessage
	currentResponse string
	*openai.Client
	model  string
	stream bool
}

func (c *Client) Prompt(ctx context.Context, prompt string) (types.StreamReader, error) {
	c.messages = append(c.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
	for _, message := range c.messages {
		utils.LogToFile("openai.log", "info", fmt.Sprintf("role: %s, message: %s", message.Role, message.Content))
	}
	req := openai.ChatCompletionRequest{
		Model:    string(c.model),
		Messages: c.messages,
		Stream:   c.stream,
	}
	chatStream, err := c.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return &streamReader{
		stream: chatStream,
	}, nil
}

func (c *Client) GetDelta(ctx context.Context, stream types.StreamReader) (*types.ChatResponse, types.StreamReader, error) {
	openaiStream := stream.(*streamReader)
	done := !stream.Scan()
	c.currentResponse = c.currentResponse + openaiStream.response.Choices[0].Delta.Content
	if done {
		c.messages = append(c.messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: c.currentResponse,
		})
		c.currentResponse = ""
	}
	return &types.ChatResponse{
		Done: done,
		Text: openaiStream.response.Choices[0].Delta.Content,
	}, stream, nil
}
