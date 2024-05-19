package randomness

import (
	"regexp"
	"strings"
)

var (
	// Represent all default values
	DEFAULT_INROW      = 3
	DEFAULT_VOCAL      = false
	DEFAULT_DIGIT      = true
	DEFAULT_CONSONANT  = true
	DEFAULT_WHITELISTS = []string{}
	DEFAULT_BLACKLISTS = []string{}
	DEFAULT_SPACES     = []string{}

	// All vocals
	VOCALS = map[rune]struct{}{
		'a': {},
		'e': {},
		'i': {},
		'o': {},
		'u': {},
		'y': {},
		'A': {},
		'E': {},
		'I': {},
		'O': {},
		'U': {},
		'Y': {},
	}
)

type Randomness struct {
	// Config represents the configuration that all methods within the Randomness structure will use
	Conifg Config
}

type Config struct {
	// When Vocals, Consonant, Digits etc...
	// have been found to reapet x amount of time in a row, it will be a random value.
	InRow int
	// Decide whether triggers should be triggered based on Case sensitive or not
	//CaseSensitive bool
	// If set to true, if a vocal character is discovered (togehte with any other "triggered" set to true).
	// It will trigger the "InRow" and if it reaches "x" of these triggers in a row. The full value will be threated as a random value.
	Vocal     bool
	Digit     bool
	Consonant bool
	// BlackRegex is used to find patterns in a string that will be triggered as blacklisted
	BlackRegex string
	// BlackRegex is used to find patterns in a string that will be triggered as whitelisted
	WhiteRegex string
	// Whitelist represents a list that if a certain subvalue (keyword) is found in the full value being tested,
	// the subvalue in the whitelist will not be affected by triggers
	// The whitelist can contain either a single [char]acter or a keyword.
	Whitelist []string
	// Blacklist represents a list that if a certain subvalue (keyword) is found in the full value being tested,
	// the subvalue in the blacklist will make the full value threated as a random value.
	// The whitelist can contain either a single [char]acter or a keyword.
	Blacklist []string
	// Spaces repesents a list of valid space characters (Ex: <space>, <newline>, _, - etc...)
	Spaces []rune
	// This maps are local and all public lists that are set will be transformed into maps for fast key lookup
	//whitelists map[string]struct{}
	//blacklists map[string]struct{}
	spaces map[rune]struct{}
}

func NewRandomness(config Config) Randomness {
	config.spaces = runelistToMap(config.Spaces)

	return Randomness{
		Conifg: config,
	}
}

// Set the needed values within the config
func SetNeededValues(config Config) Config {
	if config.InRow == 0 {
		config.InRow = DEFAULT_INROW
	}
	return config
}

// Set the config to its default values
func DefaultConfig() Config {
	return Config{
		InRow:     DEFAULT_INROW,
		Digit:     DEFAULT_DIGIT,
		Vocal:     DEFAULT_VOCAL,
		Consonant: DEFAULT_CONSONANT,
	}
}

// Check if the string given is likely to be random
// Whitelist and regex will be prioritized if a match is found
func (r Randomness) IsRandom(s string) bool {
	hit := 0

	if r.IsWhite(s) {
		return false // => Not random
	}
	if r.IsBlack(s) {
		return true // => Random
	}

	// Analyze the source string and check it it's random
	for _, char := range s {
		if !r.IsTrigger(char) {
			hit = 0
		} else {
			hit++
		}
		//The value is random
		if hit == r.Conifg.InRow {
			return true
		}
	}
	return false
}

// Check white-regex/list

func (r Randomness) IsWhite(s string) bool {
	if (len(r.Conifg.Whitelist) > 0 && r.IsWhitelist(s)) ||
		(len(r.Conifg.WhiteRegex) > 0 && r.IsRegex(r.Conifg.WhiteRegex, s)) {
		return true
	}
	return false
}

// Check black-regex/list
func (r Randomness) IsBlack(s string) bool {
	// Check black-regex/list
	if (len(r.Conifg.Blacklist) > 0 && r.IsBlacklist(s)) ||
		(len(r.Conifg.BlackRegex) > 0 && r.IsRegex(r.Conifg.BlackRegex, s)) {
		return true
	}
	return false
}

// Check if a given char will pass or trigger in relation to the configuration of the Randomness structure
func (r Randomness) IsTrigger(char rune) bool {
	switch {
	case r.Conifg.Consonant && r.IsConsonant(char):
		return true
	case r.Conifg.Digit && r.IsDigit(char):
		return true
	case r.Conifg.Vocal && r.IsVocal(char):
		return true
	default:
		return false
	}
}

func (r Randomness) IsRegex(pattern, str string) bool {
	if ok, err := regexp.MatchString(pattern, str); err != nil {
		return false
	} else {
		return ok
	}
}

// Check if a string is within the whitelist
func (r Randomness) IsWhitelist(s string) bool {
	return listItemContains(r.Conifg.Whitelist, s)
}

// Check if a string is within the blacklist
func (r Randomness) IsBlacklist(s string) bool {
	return listItemContains(r.Conifg.Blacklist, s)
}

// Check if the given string is a space in relation to the configuration of the Randomness structure
func (r Randomness) IsSpace(char rune) bool {
	_, ok := r.Conifg.spaces[char]
	return ok
}

// Check if the given rune is a vocal (including y and Y)
func (r Randomness) IsVocal(char rune) bool {
	_, ok := VOCALS[char]
	return ok
}

// Check if the given rune is a number ([0-9])
func (r Randomness) IsDigit(char rune) bool {
	return (char >= '0' && char <= '9')
}

// Check if the char is a consonant
// Note : (Alias of : !r.IsVocal(char))
func (r Randomness) IsConsonant(char rune) bool {
	return !r.IsVocal(char)
}

// Check if any item in thte list are a substring of the given string
func listItemContains(lst []string, s string) bool {
	for _, i := range lst {
		if strings.Contains(s, i) {
			return true
		}
	}
	return false
}

func runelistToMap(lst []rune) map[rune]struct{} {
	m := make(map[rune]struct{})
	for _, i := range lst {
		m[i] = struct{}{}
	}
	return m
}
