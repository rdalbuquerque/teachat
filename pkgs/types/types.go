package types

type Model string

const (
	Llama3 Model = "llama3"
	GPT35  Model = "gpt-3.5-turbo-0125"
	GPT4   Model = "gpt-4"
	GPT4o  Model = "gpt-4o"
)

type StreamReader interface {
	Bytes() []byte
	Scan() bool
	Close()
}

type ChatResponse struct {
	Done bool
	Text string
}

type ChatStream struct {
	Response *ChatResponse
	Stream   StreamReader
}
