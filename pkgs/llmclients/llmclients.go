package llmclients

import (
	"teachat/pkgs/llminterface"
	"teachat/pkgs/ollama"
	"teachat/pkgs/openai"
)

type Platform string

const (
	OpenAI Platform = "openai"
	Ollama Platform = "ollama"
)

type InitFunc func(bool) llminterface.Client

var PlatformInitialization = map[Platform]InitFunc{
	OpenAI: openai.New,
	Ollama: ollama.New,
}
