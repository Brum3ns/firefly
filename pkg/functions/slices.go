package functions

// Split a string by comma but ignore escaped comma characters (\,) to be splitted.
// Return a string based list of all the items.
func SplitEscape(s string, sep rune) []string {
	var (
		l   []string
		str string
	)
	for idx, r := range s {
		if len(s) >= 2 && r == sep {
			if s[idx-1] != '\\' {
				l = append(l, str)
				str = ""
				continue

			} else if s[idx-1] == '\\' {
				str = str[:len(str)-1]
			}
		}
		str += string(r)

		if idx == len(s)-1 {
			l = append(l, str)
		}
	}
	return l
}
