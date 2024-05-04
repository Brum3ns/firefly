package payloads

var DEFAULT_CHARS = []rune{
	'{',
	'}',
	'[',
	']',
	'(',
	')',
	'+',
	'-',
	'*',
	'/',
	'%',
	'=',
	'!',
	'>',
	'<',
	'&',
	'|',
	'^',
	'~',
	'.',
	',',
	':',
	';',
	'"',
	'\'',
	'#',
	'@',
	'?',
	'_',
	'$',
	'\\',
}

type Mutation struct {
	Chars map[rune]payloadInfo
	Seed  string
	cwe   CWE
	feedback
}

type payloadInfo struct {
	level int
	// Relation contains the char that the current char have a relation with
	// and how high the relation is given by a int value.
	relation map[rune]int
}

type feedback struct {
	payloadsSuccess map[string]int
	payloadFail     map[string]int
	cache           []string
}

func NewMutation() Mutation {
	return Mutation{}
}

func (m Mutation) Run() feedback {
	return feedback{}
}

// Set custom chars to use for the mutation process.
// True/False to also include the default characters (no duplicates)
func (m Mutation) makeChars(chars []rune, includeDefault bool) {

}

// Set the seed, this will be the root payload for the mutation
func (m Mutation) SetSeed(seed string) {
	m.Seed = seed
}

func (f Mutation) SetCWEFocus() {

}

// Set the chars to be used within the mutation process and the one to focus on (if any)
func (m Mutation) SetChars(chars []rune, charsFocus ...rune) {

}

// Set a char relation (Ex: char '{' has a relation to char '}', '$' or ';')
func (m Mutation) SetCharRelation(char rune, charRelation ...rune) {
	if charRelation == nil {
		return
	}
	// Code...
}

func (m Mutation) SetCharsIgnore(chars map[rune]int) {

}

func (f Mutation) AppendFeedback() {

}

func (f Mutation) GetPayload() {

}
