package teamsg

import "teachat/pkgs/types"

type ChatPromptMsg string
type ChatStreamMsg types.StreamReader
type ChatStreamDeltaMsg types.ChatStream
type ChatStreamCloseMsg types.ChatStream
type ModelSelectedMsg string
type GetSupportedModelsMsg bool
type ModelsMsg struct {
	Models []string
}
