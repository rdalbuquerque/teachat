package llmclients

import (
	"teachat/pkgs/llminterface"
	"teachat/pkgs/ollama"
	"teachat/pkgs/openai"
	"teachat/pkgs/types"
)

type InitFunc func(bool) llminterface.Client

var PlatformInitialization = map[types.LLMPlatform]InitFunc{
	types.OpenAI: openai.New,
	types.Ollama: ollama.New,
}
