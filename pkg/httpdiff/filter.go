package httpdiff

type Filter struct {
	Headers []string
	//HTMLFilter
}

/*
type HeaderFilter struct {
	Header httpnode.HeaderMergeNode
}
*/
/* func (hf HeaderFilter) GetHeader(header string) httpnode.HeaderInfo {
	return hf.Header[header]
}

func (hf HeaderFilter) HasHeader(header string) bool {
	_, ok := hf.Header[header]
	return ok
}
*/
