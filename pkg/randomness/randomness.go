package randomness

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// Represent all default values
	DEFAULT_INROW      = 3
	DEFAULT_VOCAL      = false
	DEFAULT_DIGIT      = true
	DEFAULT_CONSONANT  = true
	DEFAULT_SPACES     = []rune{' ', '-', '_'}
	DEFAULT_WHITELISTS = []string{}
	DEFAULT_BLACKLISTS = []string{}
	DEFAULT_WHITEREGEX = ""
	DEFAULT_BLACKREGEX = "(" +
		"[Jj]x|[Jj]z|[Kk]q|[Kk]x|[Qq]g|[Qq]j|[Qq]k|[Qq]x|[Qq]z" +
		"[Vv]x|[Vv]z|[Ww]x|[Ww]z|[Xx]j|[Xx]k|[Xx]q|[Zz]j|[Zz]q" +
		"[Zz]x|[Tt]z|[Ss]x|[Qq]d|[Dd]h|[Hh]c|[Cc]q|[Ll]k|[Uu]y" +
		")"

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
	Config
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
	// Internal variables that are modified based on the given [Config]uration
	blackRegex *regexp.Regexp
	whiteRegex *regexp.Regexp
	spaces     map[rune]struct{}
}

// Return the Randomness structure adjusted to the configurations given.
// Note : if the config variable InRow is set to zero, the default value will be set
func NewRandomness(config Config) (Randomness, error) {
	rand := Randomness{
		Config: config,
	}

	if rand.Config.InRow == 0 {
		rand.Config.InRow = DEFAULT_INROW
	}

	// Generate lookup map for the given spaces
	rand.spaces = runeListToLookupMap(rand.Config.Spaces)

	// Set black/white-regex

	if err := rand.SetBlackRegex(rand.Config.BlackRegex); err != nil {
		return rand, errors.New("invalid regex set")
	}
	if err := rand.SetWhiteRegex(rand.Config.WhiteRegex); err != nil {
		return rand, errors.New("invalid regex set")
	}
	// Set black/white-list
	rand.SetBlacklist(rand.Config.Blacklist)
	rand.SetWhitelist(rand.Config.Whitelist)

	return rand, nil
}

// Set the config to its default values
func DefaultConfig() Config {
	return Config{
		InRow:      DEFAULT_INROW,
		Vocal:      DEFAULT_VOCAL,
		Digit:      DEFAULT_DIGIT,
		Consonant:  DEFAULT_CONSONANT,
		BlackRegex: DEFAULT_BLACKREGEX,
		WhiteRegex: DEFAULT_WHITEREGEX,
		Spaces:     DEFAULT_SPACES,
		Whitelist:  DEFAULT_WHITELISTS,
		Blacklist:  DEFAULT_BLACKLISTS,
	}
}

// Set the whitelist. When a keyword in the whitelist is a sub-string of the tested value, the value will be treated as a non-random value
func (r *Randomness) AppendWhitelist(lst []string) {
	r.Config.Blacklist = append(r.Config.Blacklist, lst...)
}

// Set the blacklist. When a keyword in the blacklist is a sub-string of the tested value, the value will be treated as a random value
func (r *Randomness) AppendBlacklist(lst []string) {
	r.Config.Whitelist = append(r.Config.Whitelist, lst...)
}

// Set the whitelist. When a keyword in the whitelist is a sub-string of the tested value, the value will be treated as a non-random value
func (r *Randomness) SetWhitelist(lst []string) {
	r.Config.Blacklist = lst
}

// Set the blacklist. When a keyword in the blacklist is a sub-string of the tested value, the value will be treated as a random value
func (r *Randomness) SetBlacklist(lst []string) {
	r.Config.Whitelist = lst
}

// Set the blackregex that makes the string containg the keyword be treated as a *non-random value*
// Note : This method must be used when setting a regex for the Randomness structure
func (r *Randomness) SetWhiteRegex(regex string) error {
	if re, err := regexp.Compile(regex); err != nil {
		return err
	} else {
		r.whiteRegex = re
		return nil
	}
}

// Set the blackregex that makes the string containg the keyword be treated as a *random value*
// Note : This method must be used when setting a regex for the Randomness structure
func (r *Randomness) SetBlackRegex(regex string) error {
	if re, err := regexp.Compile(regex); err != nil {
		return err
	} else {
		r.Config.blackRegex = re
		return nil
	}
}

// Set trigger, this parameter will be seen as a possible start of a random value.
// If the trigger(s) that are set meet the InRow parameter value, it is then treated as a random value
// Note : Valid trigger(s) are: vocal, digit or consonant
func (r *Randomness) SetTrigger(s string) error {
	switch strings.ToLower(s) {
	case "vocal":
		r.Config.Vocal = true
	case "digit":
		r.Config.Digit = true
	case "consonant":
		r.Config.Consonant = true
	default:
		return errors.New("invalid trigger set, the valid triggers are: \"vocal\", \"digit\" or \"consonant\"")
	}
	return nil
}

// unSet trigger, this parameter not be seen as a possible start of a random value.
// Note : Valid trigger(s) are: vocal, digit or consonant
func (r *Randomness) UnsetTrigger(s string) error {
	switch strings.ToLower(s) {
	case "vocal":
		r.Config.Vocal = false
	case "digit":
		r.Config.Digit = false
	case "consonant":
		r.Config.Consonant = false
	default:
		return errors.New("invalid trigger set, the valid triggers are: \"vocal\", \"digit\" or \"consonant\"")
	}
	return nil
}

// Check if the string given is likely to be random
// Whitelist and regex will be prioritized if a match is found
func (r *Randomness) IsRandom(s string) bool {
	hit := 0
	for _, char := range s {
		if !r.IsTrigger(char) {
			hit = 0
		} else {
			hit++
		}
		//The value is random
		if hit == r.Config.InRow {
			return true
		}
	}

	if r.ContainsValidValue(s) {
		return false // => Not random
	}
	if r.ContainsInvalidValue(s) {
		return true // => Random
	}
	return false
}

// Check white-regex/list
func (r *Randomness) ContainsValidValue(s string) bool {
	if s == "" {
		return false
	}
	if (len(r.Config.Whitelist) > 0 && r.IsWhitelist(s)) ||
		(r.IsWhiteRegex(s)) {
		return true
	}
	return false
}

// Check black-regex/list
func (r *Randomness) ContainsInvalidValue(s string) bool {
	if s == "" {
		return false
	}
	if (len(r.Config.Blacklist) > 0 && r.IsBlacklist(s)) ||
		(r.IsBlackRegex(s)) {
		return true
	}
	return false
}

// Check if a given char will pass or trigger in relation to the configuration of the Randomness structure
func (r *Randomness) IsTrigger(char rune) bool {
	switch {
	case r.Config.Consonant && r.IsConsonant(char):
		return true
	case r.Config.Digit && r.IsDigit(char):
		return true
	case r.Config.Vocal && r.IsVocal(char):
		return true
	default:
		return false
	}
}

// Check if the string match the "whitelisted" regex
func (r *Randomness) IsWhiteRegex(s string) bool {
	if s == "" {
		return false
	}
	return issetRegex(r.whiteRegex) && r.whiteRegex.MatchString(s)
}

// Check if the string match the "blacklisted" regex
func (r *Randomness) IsBlackRegex(s string) bool {
	if s == "" {
		return false
	}
	return issetRegex(r.blackRegex) && r.Config.blackRegex.MatchString(s)
}

// Check if a string is within the whitelist
func (r *Randomness) IsWhitelist(s string) bool {
	return listItemContains(r.Config.Whitelist, s)
}

// Check if a string is within the blacklist
func (r *Randomness) IsBlacklist(s string) bool {
	return listItemContains(r.Config.Blacklist, s)
}

// Check if the given string is a space in relation to the configuration of the Randomness structure
func (r *Randomness) IsSpace(char rune) bool {
	_, ok := r.spaces[char]
	return ok
}

// Check if the given rune is a vocal (including y and Y)
func (r *Randomness) IsVocal(char rune) bool {
	_, ok := VOCALS[char]
	return ok
}

// Check if the given rune is a number ([0-9])
func (r *Randomness) IsDigit(char rune) bool {
	return (char >= '0' && char <= '9')
}

// Check if the char is a consonant
// Note : (Alias of : !r.IsVocal(char))
func (r *Randomness) IsConsonant(char rune) bool {
	return !r.IsVocal(char)
}

// Check if any item in thte list are a substring of the given string
func listItemContains(lst []string, s string) bool {
	if s == "" {
		return false
	}
	for _, i := range lst {
		if strings.Contains(s, i) {
			return true
		}
	}
	return false
}

// Convert a rune list into a lookup map
func runeListToLookupMap(lst []rune) map[rune]struct{} {
	m := make(map[rune]struct{})
	for _, i := range lst {
		m[i] = struct{}{}
	}
	return m
}

// Make sure a compiled regexp is not empty
func issetRegex(re *regexp.Regexp) bool {
	return re != nil && re.String() != ""
}
