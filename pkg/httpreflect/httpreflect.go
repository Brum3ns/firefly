package httpreflect

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

var (
	DEFAULT_INDEX_START_LENGTH = 11
	DEFAULT_INDEX_END_LENGTH   = 11
)

type Reflect struct {
	//Canary      string
	//prefix       string
	//suffix       string
	//regexBody    *regexp.Regexp
	//regexHeaders *regexp.Regexp

	Config
}

type Config struct {
	// Index length represents the character amount that could be used as a prefix and suffix
	// (canary surrounding) to be able to find the canary in a given data string
	IndexEndLength   int
	IndexStartLength int
	// The Canary to look for in the HTTP response
	Canary string
}

type Surrounding struct {
	Prefix      string
	Suffix      string
	PrefixIndex int
	SuffixIndex int
}

// Create a new reflect and adapt it to the given configurations
func NewReflect(config Config) Reflect {
	// Validate the prefix and suffix index length.
	// If the values are invalid, set the default value
	if config.IndexStartLength <= 0 {
		config.IndexStartLength = DEFAULT_INDEX_START_LENGTH
	}
	if config.IndexEndLength <= 0 {
		config.IndexEndLength = DEFAULT_INDEX_END_LENGTH
	}

	return Reflect{
		Config: config,
	}
}

// Extract all the surroundings for all discovered canaries found in the given data string
func (r Reflect) GetAllCanarySurroundings(data string) []Surrounding {
	var lst []Surrounding
	for {
		surr, ok := r.GetCanarySurrounding(data)
		if !ok {
			break
		}
		lst = append(lst, surr)

		// Update the data by removing anything behind the last discovered index
		data = data[surr.SuffixIndex:]

		// Fail safe to avoid getting stuck in an infinity loop
		if surr.SuffixIndex <= 0 {
			break
		}
	}
	return lst
}

// Extract the surrounding (prefix and suffix)
// for the *first discovered canary* wihtin the data string and return it
func (r Reflect) GetCanarySurrounding(data string) (Surrounding, bool) {
	// Make sure the canary is included in the given data
	// If not, return an empty Surrounding and false, to indicate that no surrounding was found
	startIndex := strings.Index(data, r.Canary)
	if startIndex == -1 {
		return Surrounding{}, false
	}

	// Set prefix details
	var (
		prefixIndex = startIndex
		prefix      = data[:prefixIndex]
	)
	// If the prefix's length is longer than the prefered index start length, trunicate it
	// Note : The prefix do not need to check if the canary exists within it. Since the first canary surrounding will only be found.
	if len(prefix) > r.Config.IndexStartLength {
		prefixIndex = (len(prefix) - r.Config.IndexStartLength)
		prefix = prefix[prefixIndex:]
	}

	// Set suffix details
	var (
		suffixIndex = startIndex + len(r.Canary)
		suffix      = data[suffixIndex:]
	)
	// Check if the suffix contains the known canary value.
	// If the suffix contains the canary, remove it and keep the pattern before the canary index hits.
	if canaryIndex := strings.Index(suffix, r.Canary); canaryIndex != -1 {
		suffix = suffix[:canaryIndex]
	}
	// If the suffix's length is longer than the prefered index end length, trunicate it
	if len(suffix) > r.Config.IndexEndLength {
		suffix = suffix[:r.Config.IndexEndLength]
	}

	return Surrounding{
		Prefix:      prefix,
		Suffix:      suffix,
		PrefixIndex: prefixIndex,
		SuffixIndex: suffixIndex,
	}, true
}

// Extract all values in relation to the given surroundings
// Note : Surrounding should already have been extracted with a known canary.
func (r Reflect) ExtractAll(data string, surroundings []Surrounding) []string {
	var canaries []string

	for _, surrounding := range surroundings {
		canaryReflect, arryIdx, ok := r.Extract(data, surrounding)
		canaries = append(canaries, canaryReflect)
		if !ok {
			break
		}
		data = data[arryIdx[1]:]
	}

	return canaries
}

func (r Reflect) Extract(data string, surrounding Surrounding) (string, [2]int, bool) {
	if r.Canary == "" || data == "" {
		return "", [2]int{}, false
	}

	prefixIndex := strings.Index(data, surrounding.Prefix)
	if prefixIndex == -1 {
		return "", [2]int{}, false
	}

	startIndex := prefixIndex + len(surrounding.Prefix)

	endIndex := strings.Index(data[startIndex:], surrounding.Suffix)
	if endIndex == -1 {
		return "", [2]int{}, false
	}
	endIndex += startIndex

	// If the full data is the reflected value:
	if startIndex == 0 && endIndex == 0 {
		return data, [2]int{prefixIndex, endIndex}, true

		// The prefix is having a value but the suffix is empty (end of data)
	} else if (endIndex - startIndex) == 0 {
		return data[startIndex:], [2]int{prefixIndex, endIndex}, true
	}

	return data[startIndex:endIndex], [2]int{prefixIndex, endIndex}, true
}

// Set a new Canary
func (r *Reflect) SetCanary(value string) {
	r.Canary = value
}

// Return the statistic canary value
func (r Reflect) GetCanary() string {
	return r.Canary
}

func MakeHash(prefix, suffix string) string {
	hash := md5.Sum([]byte(prefix + suffix))
	return hex.EncodeToString(hash[:])
}

func MergeUniqueSurroundings(knownSurroundings, newSurroundings []Surrounding) []Surrounding {
	var surr = knownSurroundings
	for _, newSurr := range newSurroundings {
		ok := true
		for _, knownSurr := range knownSurroundings {
			if newSurr == knownSurr {
				ok = false
				break
			}
		}
		if ok {
			surr = append(surr, newSurr)
		}

	}
	return surr
}
