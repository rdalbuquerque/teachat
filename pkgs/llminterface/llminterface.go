package llminterface

import (
	"context"
	"teachat/pkgs/types"
)

type Client interface {
	Prompt(context.Context, string) (types.StreamReader, error)
	GetDelta(context.Context, types.StreamReader) (*types.ChatResponse, types.StreamReader, error)
	SetModel(types.Model)
}
