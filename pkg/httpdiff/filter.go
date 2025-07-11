package httpdiff

import "github.com/Brum3ns/firefly/pkg/httpnode"

type HeaderFilter struct {
	Header httpnode.HeaderMergeNode
}

func (hf HeaderFilter) GetHeader(header string) httpnode.HeaderInfo {
	return hf.Header[header]
}

func (hf HeaderFilter) HasHeader(header string) bool {
	_, ok := hf.Header[header]
	return ok
}
