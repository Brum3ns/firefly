package payloads

type Mutade struct {
	payloads       map[string]bool
	patterns       map[string]int
	characterPairs map[rune]string
	characters     map[rune]int
	properties     Properties
}

type Properties struct {
	// Represet if the character haven't had any success 'x' amount of times
	GiveUpCharacter int

	MaxRepeat int
	MinRepeat int
}

func NewMutade(properties Properties) *Mutade {
	return &Mutade{
		properties: properties,
	}

}

func (mutade *Mutade) GetPayload(payload string, success bool) {
	/* payloadChars := inputChars(payload)
	if success {

	} else {

	} */
}

func (mutade *Mutade) Feed(payload string) {

}

func (mutade *Mutade) GetPattern() string {
	return ""
}

// Return all the chars the input contains and how many of each
func inputChars(input string) map[rune]int {
	m := make(map[rune]int)
	for _, r := range input {
		m[r]++
	}
	return m
}
