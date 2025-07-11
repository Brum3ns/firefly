package extract

type MergeExtractResult struct {
	Regexes  map[string][]int `json:"regexes"`
	Keywords map[string][]int `json:"keywords"`
}

func NewMerge() MergeExtractResult {
	return MergeExtractResult{
		Regexes:  make(map[string][]int),
		Keywords: make(map[string][]int),
	}
}
