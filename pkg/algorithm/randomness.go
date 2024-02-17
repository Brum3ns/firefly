package algorithm

var vocals = map[rune]int{
	'a': 0,
	'e': 0,
	'i': 0,
	'o': 0,
	'u': 0,
	'A': 1,
	'E': 1,
	'I': 1,
	'O': 1,
	'U': 1,
}

// Return a map that contains the charFrequency and true/false if the value is random (*very likely to be random*):
func IsRandom(s string) (map[rune]int, bool) {
	var (
		hit            = 0
		consonantInRow = 3
		charFrequency  = make(map[rune]int)
	)

	//Calculate constant that comes in a row: (usually 3-4+ result in a random string)
	for _, r := range s {
		if _, ok := vocals[r]; ok {
			hit = 0
		}
		hit++
		//Likely to be random:
		if hit == consonantInRow {

			//Save random values and make an "average entropy value" adapted to the applicato so we can learn the randomness sturcture
			return charFrequency, true
		}
		charFrequency[r]++
	}
	return charFrequency, false
}
