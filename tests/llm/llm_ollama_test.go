package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Brum3ns/firefly/internal/llm"
)

func Test_ollama(t *testing.T) {

	const prompt = "Are you running?"

	l, err := llm.NewLLM(llm.ConfigOllama{
		OllamaURI:  "http://127.0.0.1:11434/",
		OllamaRole: "guest",
		Model:      "gemma3:4b",
		Stream:     false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	respLLM, err := l.Generate(prompt)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(respLLM)
}
