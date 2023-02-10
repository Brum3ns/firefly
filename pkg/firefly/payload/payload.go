package payload

import "strings"

type payloadData struct {
	s_given         string
	a_charPairs1    [4]string
	a_charPairs2    [4]string
	m_chars         map[string][]string
	m_verifiedChars map[rune][]string

	//Return
	l_patterns []string
}

/**Detect payload from a known pattern given in a [m]ap. Return the possible payload modification
that the target should contain in the response (default behavior)*/
func Detect(m map[rune][]string, p string) ([]string, bool) {
	pyld := &payloadData{
		s_given:         "\"x'", //TEMP - replace with "p"
		m_verifiedChars: m,
		a_charPairs1:    [4]string{")", "]", "}", ">"},
		a_charPairs2:    [4]string{"(", "[", "{", "<"},
		l_patterns:      []string{},
	}
	//Check what chars that are included in the payload:
	pyld.m_chars = pyld.Chars()

	//exit if the payload do not contains any chars:
	if len(pyld.m_chars) <= 0 {
		return nil, false
	}

	//Add combination that the payload can appear in related to the given `m_char` result:
	for s_rootChar, l_chars := range pyld.m_chars {
		s_newPayload := pyld.s_given
		for _, s_char := range l_chars {
			s_newPayload = strings.ReplaceAll(s_newPayload, s_rootChar, s_char)
			pyld.l_patterns = append(pyld.l_patterns, s_newPayload)
		}
	}

	return pyld.l_patterns, true
}

/**Extract chars from the payload and add the known pattern/encoding the target modifies the chars into. Return the list of known char modification*/
func (payload *payloadData) Chars() map[string][]string {
	m_chars := make(map[string][]string)

	//Extract chars that are included in the payload:
	for _, r_character := range payload.s_given {
		s_character := string(r_character)

		//Char exist, add the known char modification (list):
		if len(payload.m_verifiedChars[r_character]) > 0 {
			m_chars[s_character] = append(m_chars[s_character], payload.m_verifiedChars[r_character]...)
		}
	}
	return m_chars
}
