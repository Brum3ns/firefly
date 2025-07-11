package payload

import (
	"errors"
	"fmt"
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

func ReplaceInOrder(payload string, values []string) (string, error) {
	for _, i := range values {
		l := strings.Split(i, " => ")
		if len(l) != 2 {
			return payload, fmt.Errorf("could not replace payload, invalid value:[%s]", i)
		}
		payload = strings.ReplaceAll(payload, l[0], l[1])
	}
	return payload, nil
}
