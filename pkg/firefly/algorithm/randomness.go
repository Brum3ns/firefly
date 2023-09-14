package algorithm

//TODO
/*
func Entropy() {

} */

// Use the 'Shannon entropy' algorithm to detect randomness combined with patterns and whitelist of characters.
// Return: Entropy value | likely to be random(-1=no, 1=yes, 0=unsure) | checked status(true/false)
/*func Detect(s string) (float64, int, bool) {
	if len(s) == 0 {
		return 0, -1, false
	}
	var (
		m      = map[rune]float64{}
		hm     float64
		likely = 0
	)

		var (


			length = len(s)
			amount = 5
			InRow  = 0 // (InRow >= 'amount') constant (*or* mixed with digits) in a row is likely to be a random generated string
		)

		if length >= 8 {
			amount -= 1
		} else if length >= 10 {
			amount -= 2
		}

		for _, rn := range s {
			if _, ok := global.RANDOMNESS_WHITELIST[rn]; !ok {
				return 0, 0, false
			}

			if likely == 0 {
				_, ok := global.CONSONANTS_DIGITS[rn]
				switch {
				case InRow >= amount:
					likely = 1
				case ok:
					InRow++
				default:
					InRow = 0
				}
			}
			m[rn]++
		}

	for _, c := range m {
		hm += (c * math.Log2(c))
	}
	l := float64(len(s))
	entropyValue := math.Log2(l) - hm/l

	return entropyValue, likely, true
}
*/
