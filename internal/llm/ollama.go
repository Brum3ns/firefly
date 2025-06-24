package llm

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type LLM struct {
	config ConfigOllama
}

type ConfigOllama struct {
	OllamaURI  string
	OllamaRole string
	Model      string
	Stream     bool
}

type Generate struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type Chat struct {
	Model   string  `json:"model"`
	Message Message `json:"message"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	response string `json:"response"`
}

func NewLLM(config ConfigOllama) (LLM, error) {
	return LLM{
		config: config,
	}, nil
}

func (llm *LLM) Generate(prompt string) (Response, error) {
	var response Response

	jsonData, err := json.Marshal(Generate{
		Model:  llm.config.Model,
		Prompt: prompt,
		Stream: llm.config.Stream,
	})

	if err != nil {
		return response, err
	}
	respOllama, err := http.Post(
		llm.config.OllamaURI,
		"application/json",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return response, err
	}

	// Decode JSON from io.Reader into struct
	if err := json.NewDecoder(respOllama.Body).Decode(&response); err != nil && err != io.EOF {
		return response, err
	}

	return response, nil
}

func (llm *LLM) Chat() (string, error) {
	return "", nil
}
