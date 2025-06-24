package payload

import (
	"errors"
	"strings"
)

// Wordlist structure stores the wordlist and tags
type Payload struct {
	wordlist []string
	//wordlistVerify []string
	config Config
}

type Config struct {
	WordlistFile  string
	PayloadVerify string
}

// Create a new payload instance
func NewPayload(config Config) (Payload, error) {
	return Payload{
		config: config,
	}, nil
}

// Get the payload wordlist
func (p *Payload) GetWordlist() []string {
	return p.wordlist
}

// Set the payload wordlist to be empty
func (p *Payload) SetEmptyWordlist() {
	p.wordlist = []string{}
}

// Make the payload wordlist based on file in config
func (p *Payload) MakeWordlist() error {
	if p.config.WordlistFile == "" {
		return errors.New("wordlist is not set")
	}
	var err error
	p.wordlist, err = readFile(p.config.WordlistFile)
	return err
}

func Insert(source, placeholder, payload string) string {
	return strings.ReplaceAll(source, placeholder, payload)
}
