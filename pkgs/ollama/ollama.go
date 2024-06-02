package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"teachat/pkgs/llminterface"
	"teachat/pkgs/types"
	"teachat/pkgs/utils"

	"github.com/ollama/ollama/version"
)

func GetSupportedModels() []types.Model {
	return []types.Model{
		types.Llama3,
	}
}

type Client struct {
	base     *url.URL
	http     *http.Client
	stream   bool
	model    types.Model
	messages []Message
}

type OllamaHost struct {
	Scheme string
	Host   string
	Port   string
}

type streamReader struct {
	*bufio.Scanner
	io.ReadCloser
}

func (s *streamReader) Scan() bool {
	return s.Scanner.Scan()
}

func (s streamReader) Bytes() []byte {
	return s.Scanner.Bytes()
}

func (s *streamReader) Close() {
	s.ReadCloser.Close()
}

func GetOllamaHost() (OllamaHost, error) {
	defaultPort := "11434"

	hostVar := os.Getenv("OLLAMA_HOST")
	hostVar = strings.TrimSpace(strings.Trim(strings.TrimSpace(hostVar), "\"'"))

	scheme, hostport, ok := strings.Cut(hostVar, "://")
	switch {
	case !ok:
		scheme, hostport = "http", hostVar
	case scheme == "http":
		defaultPort = "80"
	case scheme == "https":
		defaultPort = "443"
	}

	hostport = strings.TrimRight(hostport, "/")

	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host, port = "127.0.0.1", defaultPort
		if ip := net.ParseIP(strings.Trim(hostport, "[]")); ip != nil {
			host = ip.String()
		} else if hostport != "" {
			host = hostport
		}
	}

	if portNum, err := strconv.ParseInt(port, 10, 32); err != nil || portNum > 65535 || portNum < 0 {
		return OllamaHost{}, ErrInvalidHostPort
	}

	return OllamaHost{
		Scheme: scheme,
		Host:   host,
		Port:   port,
	}, nil
}

func New(stream bool) llminterface.Client {
	ollamaHost, err := GetOllamaHost()
	if err != nil {
		panic(err)
	}
	model := types.Llama3
	return &Client{
		model:  model,
		stream: stream,
		base: &url.URL{
			Scheme: ollamaHost.Scheme,
			Host:   net.JoinHostPort(ollamaHost.Host, ollamaHost.Port),
		},
		http: http.DefaultClient,
	}
}

func (c *Client) SetModel(model types.Model) {
	c.model = model
}

func (c *Client) Prompt(ctx context.Context, prompt string) (types.StreamReader, error) {
	c.messages = append(c.messages, Message{
		Role:    "user",
		Content: prompt,
	})
	req := ChatRequest{
		Model:    string(c.model),
		Messages: c.messages,
		Stream:   utils.Ptr(c.stream),
	}
	return c.getStream(ctx, http.MethodPost, "/api/chat", req)
}

func (c Client) GetDelta(ctx context.Context, stream types.StreamReader) (*types.ChatResponse, types.StreamReader, error) {
	var resp ChatResponse
	stream.Scan()
	if err := json.Unmarshal(stream.Bytes(), &resp); err != nil {
		return nil, stream, err
	}
	return &types.ChatResponse{
		Done: resp.Done,
		Text: resp.Message.Content,
	}, stream, nil
}

func (c Client) getStream(ctx context.Context, method, path string, data any) (types.StreamReader, error) {
	var buf *bytes.Buffer
	if data != nil {
		bts, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		buf = bytes.NewBuffer(bts)
	}

	requestURL := c.base.JoinPath(path)
	utils.LogToFile("ollama.log", "info", fmt.Sprintf("requestURL: %v", requestURL.String()))
	request, err := http.NewRequestWithContext(ctx, method, requestURL.String(), buf)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/x-ndjson")
	request.Header.Set("User-Agent", fmt.Sprintf("ollama/%s (%s %s) Go/%s", version.Version, runtime.GOARCH, runtime.GOOS, runtime.Version()))

	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(response.Body)
	return &streamReader{
		Scanner:    scanner,
		ReadCloser: response.Body,
	}, nil
}
