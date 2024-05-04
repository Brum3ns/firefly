package payloads

type Relation struct {
	Payload   map[string][]string
	Character map[rune][]string
}

func NewRelation() Relation {
	return Relation{}
}

// Find a related chars within the given string when compared to chars from other strings in the memory
func (r *Relation) findChar(s string) {

}
