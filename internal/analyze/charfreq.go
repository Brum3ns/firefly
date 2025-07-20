package analyze

// CharFrequency holds character counts for a group of payloads.
type CharFrequency map[rune]int

// AnalyzeCharFrequencyByGroup computes character frequency per result group.
// Input is the output from GroupByScan.
func AnalyzeCharFrequencyByGroup(groups []ResultGroup) []CharFrequency {
	var result []CharFrequency
	for _, group := range groups {
		freq := make(CharFrequency)

		for _, payload := range group.Payloads {
			for _, char := range payload {
				freq[char]++
			}
		}
		result = append(result, freq)
	}
	return result
}
