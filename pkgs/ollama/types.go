package ollama

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"time"
)

// StatusError is an error with and HTTP status code.
type StatusError struct {
	StatusCode   int
	Status       string
	ErrorMessage string `json:"error"`
}

func (e StatusError) Error() string {
	switch {
	case e.Status != "" && e.ErrorMessage != "":
		return fmt.Sprintf("%s: %s", e.Status, e.ErrorMessage)
	case e.Status != "":
		return e.Status
	case e.ErrorMessage != "":
		return e.ErrorMessage
	default:
		// this should not happen
		return "something went wrong, please see the ollama server logs for details"
	}
}

// ImageData represents the raw binary data of an image file.
type ImageData []byte

// GenerateRequest describes a request sent by [Client.Generate]. While you
// have to specify the Model and Prompt fields, all the other fields have
// reasonable defaults for basic uses.
type GenerateRequest struct {
	// Model is the model name; it should be a name familiar to Ollama from
	// the library at https://ollama.com/library
	Model string `json:"model"`

	// Prompt is the textual prompt to send to the model.
	Prompt string `json:"prompt"`

	// System overrides the model's default system message/prompt.
	System string `json:"system"`

	// Template overrides the model's default prompt template.
	Template string `json:"template"`

	// Context is the context parameter returned from a previous call to
	// Generate call. It can be used to keep a short conversational memory.
	Context []int `json:"context,omitempty"`

	// Stream specifies whether the response is streaming; it is true by default.
	Stream *bool `json:"stream,omitempty"`

	// Raw set to true means that no formatting will be applied to the prompt.
	Raw bool `json:"raw,omitempty"`

	// Format specifies the format to return a response in.
	Format string `json:"format"`

	// KeepAlive controls how long the model will stay loaded in memory following
	// this request.
	KeepAlive *Duration `json:"keep_alive,omitempty"`

	// Images is an optional list of base64-encoded images accompanying this
	// request, for multimodal models.
	Images []ImageData `json:"images,omitempty"`

	// Options lists model-specific options. For example, temperature can be
	// set through this field, if the model supports it.
	Options map[string]interface{} `json:"options"`
}

// ChatRequest describes a request sent by [Client.Chat].
type ChatRequest struct {
	// Model is the model name, as in [GenerateRequest].
	Model string `json:"model"`

	// Messages is the messages of the chat - can be used to keep a chat memory.
	Messages []Message `json:"messages"`

	// Stream enable streaming of returned response; true by default.
	Stream *bool `json:"stream,omitempty"`

	// Format is the format to return the response in (e.g. "json").
	Format string `json:"format"`

	// KeepAlive controls how long the model will stay loaded into memory
	// followin the request.
	KeepAlive *Duration `json:"keep_alive,omitempty"`

	// Options lists model-specific options.
	Options map[string]interface{} `json:"options"`
}

// Message is a single message in a chat sequence. The message contains the
// role ("system", "user", or "assistant"), the content and an optional list
// of images.
type Message struct {
	Role    string      `json:"role"`
	Content string      `json:"content"`
	Images  []ImageData `json:"images,omitempty"`
}

// ChatResponse is the response returned by [Client.Chat]. Its fields are
// similar to [GenerateResponse].
type ChatResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   Message   `json:"message"`

	Done bool `json:"done"`

	Metrics
}

type Metrics struct {
	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
}

// Options specified in [GenerateRequest], if you add a new option here add it
// to the API docs also.
type Options struct {
	Runner

	// Predict options used at runtime
	NumKeep          int      `json:"num_keep,omitempty"`
	Seed             int      `json:"seed,omitempty"`
	NumPredict       int      `json:"num_predict,omitempty"`
	TopK             int      `json:"top_k,omitempty"`
	TopP             float32  `json:"top_p,omitempty"`
	TFSZ             float32  `json:"tfs_z,omitempty"`
	TypicalP         float32  `json:"typical_p,omitempty"`
	RepeatLastN      int      `json:"repeat_last_n,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	RepeatPenalty    float32  `json:"repeat_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	Mirostat         int      `json:"mirostat,omitempty"`
	MirostatTau      float32  `json:"mirostat_tau,omitempty"`
	MirostatEta      float32  `json:"mirostat_eta,omitempty"`
	PenalizeNewline  bool     `json:"penalize_newline,omitempty"`
	Stop             []string `json:"stop,omitempty"`
}

// Runner options which must be set when the model is loaded into memory
type Runner struct {
	UseNUMA   bool `json:"numa,omitempty"`
	NumCtx    int  `json:"num_ctx,omitempty"`
	NumBatch  int  `json:"num_batch,omitempty"`
	NumGQA    int  `json:"num_gqa,omitempty"`
	NumGPU    int  `json:"num_gpu,omitempty"`
	MainGPU   int  `json:"main_gpu,omitempty"`
	LowVRAM   bool `json:"low_vram,omitempty"`
	F16KV     bool `json:"f16_kv,omitempty"`
	LogitsAll bool `json:"logits_all,omitempty"`
	VocabOnly bool `json:"vocab_only,omitempty"`
	UseMMap   bool `json:"use_mmap,omitempty"`
	UseMLock  bool `json:"use_mlock,omitempty"`
	NumThread int  `json:"num_thread,omitempty"`

	// Unused: RopeFrequencyBase is ignored. Instead the value in the model will be used
	RopeFrequencyBase float32 `json:"rope_frequency_base,omitempty"`
	// Unused: RopeFrequencyScale is ignored. Instead the value in the model will be used
	RopeFrequencyScale float32 `json:"rope_frequency_scale,omitempty"`
}

// EmbeddingRequest is the request passed to [Client.Embeddings].
type EmbeddingRequest struct {
	// Model is the model name.
	Model string `json:"model"`

	// Prompt is the textual prompt to embed.
	Prompt string `json:"prompt"`

	// KeepAlive controls how long the model will stay loaded in memory following
	// this request.
	KeepAlive *Duration `json:"keep_alive,omitempty"`

	// Options lists model-specific options.
	Options map[string]interface{} `json:"options"`
}

// EmbeddingResponse is the response from [Client.Embeddings].
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// CreateRequest is the request passed to [Client.Create].
type CreateRequest struct {
	Model        string `json:"model"`
	Path         string `json:"path"`
	Modelfile    string `json:"modelfile"`
	Stream       *bool  `json:"stream,omitempty"`
	Quantization string `json:"quantization,omitempty"`

	// Name is deprecated, see Model
	Name string `json:"name"`
}

// DeleteRequest is the request passed to [Client.Delete].
type DeleteRequest struct {
	Model string `json:"model"`

	// Name is deprecated, see Model
	Name string `json:"name"`
}

// ShowRequest is the request passed to [Client.Show].
type ShowRequest struct {
	Model    string `json:"model"`
	System   string `json:"system"`
	Template string `json:"template"`

	Options map[string]interface{} `json:"options"`

	// Name is deprecated, see Model
	Name string `json:"name"`
}

// ShowResponse is the response returned from [Client.Show].
type ShowResponse struct {
	License    string       `json:"license,omitempty"`
	Modelfile  string       `json:"modelfile,omitempty"`
	Parameters string       `json:"parameters,omitempty"`
	Template   string       `json:"template,omitempty"`
	System     string       `json:"system,omitempty"`
	Details    ModelDetails `json:"details,omitempty"`
	Messages   []Message    `json:"messages,omitempty"`
}

// CopyRequest is the request passed to [Client.Copy].
type CopyRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// PullRequest is the request passed to [Client.Pull].
type PullRequest struct {
	Model    string `json:"model"`
	Insecure bool   `json:"insecure,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Stream   *bool  `json:"stream,omitempty"`

	// Name is deprecated, see Model
	Name string `json:"name"`
}

// ProgressResponse is the response passed to progress functions like
// [PullProgressFunc] and [PushProgressFunc].
type ProgressResponse struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// PushRequest is the request passed to [Client.Push].
type PushRequest struct {
	Model    string `json:"model"`
	Insecure bool   `json:"insecure,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Stream   *bool  `json:"stream,omitempty"`

	// Name is deprecated, see Model
	Name string `json:"name"`
}

// ListResponse is the response from [Client.List].
type ListResponse struct {
	Models []ModelResponse `json:"models"`
}

// ModelResponse is a single model description in [ListResponse].
type ModelResponse struct {
	Name       string       `json:"name"`
	Model      string       `json:"model"`
	ModifiedAt time.Time    `json:"modified_at"`
	Size       int64        `json:"size"`
	Digest     string       `json:"digest"`
	Details    ModelDetails `json:"details,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// GenerateResponse is the response passed into [GenerateResponseFunc].
type GenerateResponse struct {
	// Model is the model name that generated the response.
	Model string `json:"model"`

	//CreatedAt is the timestamp of the response.
	CreatedAt time.Time `json:"created_at"`

	// Response is the textual response itself.
	Response string `json:"response"`

	// Done specifies if the response is complete.
	Done bool `json:"done"`

	// Context is an encoding of the conversation used in this response; this
	// can be sent in the next request to keep a conversational memory.
	Context []int `json:"context,omitempty"`

	Metrics
}

// ModelDetails provides details about a model.
type ModelDetails struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

// ErrInvalidOpts is returned when invalid options are passed to the client.
var ErrInvalidHostPort = errors.New("invalid port specified in OLLAMA_HOST")

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	if d.Duration < 0 {
		return []byte("-1"), nil
	}
	return []byte("\"" + d.Duration.String() + "\""), nil
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	d.Duration = 5 * time.Minute

	switch t := v.(type) {
	case float64:
		if t < 0 {
			d.Duration = time.Duration(math.MaxInt64)
		} else {
			d.Duration = time.Duration(int(t) * int(time.Second))
		}
	case string:
		d.Duration, err = time.ParseDuration(t)
		if err != nil {
			return err
		}
		if d.Duration < 0 {
			d.Duration = time.Duration(math.MaxInt64)
		}
	default:
		return fmt.Errorf("Unsupported type: '%s'", reflect.TypeOf(v))
	}

	return nil
}
