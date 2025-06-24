package knowledge

import (
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/httpreflect"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

type Knowledge struct {
	Storage map[string][]Storage
	Merged  map[string]MergedKnowledge
}

type MergedKnowledge struct {
	Canary    string
	Responses []rhttp.Response
	//Requests                []output.Request
	HttpReflectSurroundings []httpreflect.Surrounding
	Merge                   Merge
}

func NewKnowledge() Knowledge {
	return Knowledge{
		Storage: make(map[string][]Storage),
		Merged:  make(map[string]MergedKnowledge),
	}
}

func (k *Knowledge) SetMergeKnowledge() {
	k.Merged = GetKnowledge(k.Storage)
}

func (k *Knowledge) AppendKnowledge(targetHash, payload string, resp rhttp.Response) {
	k.Storage[targetHash] = append(
		k.Storage[targetHash],
		Storage{
			Payload:    payload,
			HTMLNode:   httpnode.GetHTMLNode(resp.Body),
			HeaderNode: httpnode.GetHeaderNode(resp.Headers),
			//Extract: ,
			Response: resp,
			//Request: result.job.,
			//HttpReflectSurroundings: ,
		},
	)
}
