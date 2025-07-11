package randomness

import "math"

func LogicalRandomnessScore(s string) float64 {
	if len(s) < 4 {
		return 0.0 // too short to judge
	}

	freq := map[rune]float64{}
	for _, ch := range s {
		freq[ch]++
	}

	var entropy float64
	for _, count := range freq {
		p := count / float64(len(s))
		entropy -= p * math.Log2(p)
	}

	maxEntropy := math.Log2(float64(len(freq)))
	if maxEntropy == 0 {
		return 0.0
	}

	normalized := entropy / maxEntropy
	adjusted := normalized * math.Log2(float64(len(s)))

	// Return a score between 0.0 - ~6.0 (can clamp to 5.0)
	return adjusted
}

func IsRandomByEntropy(value string, entropy float64) bool {
	return CalcEntropy(value) > entropy
}

func CalcEntropy(value string) float64 {
	if len(value) == 0 {
		return 0
	}
	freq := make(map[rune]float64)
	for _, c := range value {
		freq[c]++
	}
	var entropy float64
	for _, count := range freq {
		p := count / float64(len(value))
		entropy -= p * math.Log2(p)
	}
	return entropy
}
