package extract

import (
	"strings"

	"github.com/dlclark/regexp2"
)

type Extract struct {
	// The key of wordlist contains the prefix of the word.
	// It stores all words / sentences that has the same prefix
	wordlist map[string][]string
	regex    []*regexp2.Regexp
	config   Config
}

type Config struct {
	FilenameWordlist string
	FilenameRegexes  string
}

type Result struct {
	Regexes  map[string]Pattern
	Keywords map[string]Pattern
}

type Pattern struct {
	Count   int
	Indexes []int
}

// NewExtract creates an Extract instance compiling all patterns
func NewExtract(config Config) (Extract, error) {
	var (
		err      error
		wordlist = make(map[string][]string)
		regexes  []*regexp2.Regexp
	)
	// Make regexes
	if len(config.FilenameRegexes) > 0 {
		regexes, err = MakeRegexps(config.FilenameRegexes)
		if err != nil {
			return Extract{}, err
		}
	}
	// Make wordlist
	if len(config.FilenameWordlist) > 0 {
		wordlist, err = MakeWordlist(config.FilenameWordlist)
		if err != nil {
			return Extract{}, err
		}
	}
	return Extract{
		wordlist: wordlist,
		regex:    regexes,
		config:   config,
	}, nil
}

func (e *Extract) Run(source string) Result {
	return Result{
		Regexes:  e.FindRegexPatterns(source),
		Keywords: e.FindPatterns(source),
	}
}

// FindRegexPatterns scans the source using all regex patterns in e.regex,
// and returns the number of matches and their starting indexes for each pattern.
func (e *Extract) FindRegexPatterns(source string) map[string]Pattern {
	matches := make(map[string]Pattern)

	for _, re := range e.regex {
		start := 0
		for {
			// Match from current start index
			match, err := re.FindStringMatch(source[start:])
			if err != nil || match == nil {
				break
			}

			// Get absolute index of match
			regexKey := re.String()

			absIdx := start + match.Index
			p := matches[regexKey]
			p.Count++
			p.Indexes = append(p.Indexes, absIdx)
			matches[regexKey] = p

			// Move past this match
			start = absIdx + match.Length
		}
	}
	return matches
}

// FindPatterns looks up the given prefixes in the source,
// and returns all words from the wordlist (matched by prefix)
// along with how many times each word reflects and at which positions.
func (e *Extract) FindPatterns(source string) map[string]Pattern {
	matches := make(map[string]Pattern)

	for prefix, patterns := range e.wordlist {
		// Check if the prefix exists in the source
		if !strings.Contains(source, prefix) {
			continue
		}
		// Extract each pattern in the source
		for _, pattern := range patterns {
			// Search for all occurrences of word in source
			start := 0
			for {
				idx := strings.Index(source[start:], pattern)
				if idx == -1 {
					break
				}
				absIdx := start + idx
				p, exists := matches[pattern]
				if !exists {
					p = Pattern{}
				}
				p.Count++
				p.Indexes = append(p.Indexes, absIdx)
				matches[pattern] = p
				start = absIdx + len(pattern)
			}
		}
	}
	return matches
}
