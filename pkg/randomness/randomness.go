package randomness

var (
	VOCALS = map[rune]int{
		'a': 0,
		'e': 0,
		'i': 0,
		'o': 0,
		'u': 0,
		'y': 0,
		'A': 1,
		'E': 1,
		'I': 1,
		'O': 1,
		'U': 1,
		'Y': 1,
	}
)

// Return a map that contains the charFrequency and true/false if the value is random (*likely to be random*):
// If "ignoreNumbers" is set to "true" numbers will be ignored to check in the string.
func IsRandom(s string, ignoreNumbers bool) bool {
	//Calculate constant that comes in a row: (usually 3-4+ result in a random string)
	hit := 0
	consonantInRow := 3
	for _, r := range s {

		// Check if the rune is a valid digit
		if IsVocal(r) || r == '_' || ignoreNumbers || IsNumber(r) {
			hit = 0
		} else {
			hit++
		}

		//Likely to be random:
		if hit == consonantInRow {

			//Save random values and make an "average entropy value" adapted to the applicato so we can learn the randomness sturcture
			return true
		}
	}
	return false
}

// Check if the given rune is a vocal (including y and Y)
func IsVocal(r rune) bool {
	_, ok := VOCALS[r]
	return ok
}

// Check if the given rune is a number ([0-9])
func IsNumber(r rune) bool {
	return (r >= '0' && r <= '9')
}
