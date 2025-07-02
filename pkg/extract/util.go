package extract

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/dlclark/regexp2"
)

func MakeRegexps(filenameRegex string) ([]*regexp2.Regexp, error) {
	var regexes []*regexp2.Regexp

	res, err := ReadFile(filenameRegex)
	if err != nil {
		return regexes, nil
	}

	// Compile regex patterns from config.Regex
	for _, pattern := range res {
		re, err := regexp2.Compile(pattern, regexp2.RE2)
		if err != nil {
			return regexes, fmt.Errorf("could not make pattern, pattern:[%v], error : %v", pattern, err)
		}
		regexes = append(regexes, re)
	}
	return regexes, nil
}

// Make a wordlist that store the prefix of all words as the map key and all relevant words/sentences
// that share the same prefix as the key's value
func MakeWordlist(FilenameWordlist string) (map[string][]string, error) {
	wordlist, err := ReadFile(FilenameWordlist)
	if err != nil {
		return make(map[string][]string), err
	}
	return buildPrefixMap(wordlist), nil
}

func ReadFile(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("coult not read the given file, file:[%s], error : %v", filename, err)
	}
	return lines, nil
}

// buildPrefixMap takes a slice of words and returns a map where each key is
// the longest common prefix shared by at least two words, and its value is
// the slice of all words sharing that prefix. Any words not in a multi‑word
// prefix get mapped to themselves.
func buildPrefixMap(words []string) map[string][]string {
	// 1) Build prefix → set of words
	prefixSets := make(map[string]map[string]struct{})
	for _, w := range words {
		for i := 1; i <= len(w); i++ {
			p := w[:i]
			if prefixSets[p] == nil {
				prefixSets[p] = make(map[string]struct{})
			}
			prefixSets[p][w] = struct{}{}
		}
	}

	// 2) Keep only prefixes covering ≥2 words, turn sets into sorted slices
	type group struct {
		prefix string
		words  []string
	}
	var groups []group
	for p, set := range prefixSets {
		if len(set) < 2 {
			continue
		}
		slice := make([]string, 0, len(set))
		for w := range set {
			slice = append(slice, w)
		}
		sort.Strings(slice)
		groups = append(groups, group{prefix: p, words: slice})
	}

	// 3) Sort groups by descending prefix‐length (longer prefixes first)
	sort.Slice(groups, func(i, j int) bool {
		if len(groups[i].prefix) != len(groups[j].prefix) {
			return len(groups[i].prefix) > len(groups[j].prefix)
		}
		// tie‐breaker: lex order (optional)
		return groups[i].prefix < groups[j].prefix
	})

	// 4) Greedily assign words to the longest prefix‑group they belong to
	result := make(map[string][]string, len(groups))
	assigned := make(map[string]bool, len(words))
	for _, g := range groups {
		var unassigned []string
		for _, w := range g.words {
			if !assigned[w] {
				unassigned = append(unassigned, w)
			}
		}
		// only keep it if ≥2 words remain unassigned
		if len(unassigned) >= 2 {
			result[g.prefix] = unassigned
			for _, w := range unassigned {
				assigned[w] = true
			}
		}
	}

	// 5) Any words still unassigned get mapped to themselves
	for _, w := range words {
		if !assigned[w] {
			result[w] = []string{w}
		}
	}
	return result
}
