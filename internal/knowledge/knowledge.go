package knowledge

import (
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

type Knowledge struct {
	Storage map[string][]Storage
	Merged  map[string]MergedKnowledge
}

type KnowledgeMeta struct {
	TargetHash          string
	Payload             string
	HTTPResponse        rhttp.Response
	ExtractResultHeader extract.Result
	ExtractResultBody   extract.Result
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

func (k *Knowledge) AppendKnowledge(meta KnowledgeMeta) {
	// Get the JSON node and ignore errors
	jsonNode, _ := httpnode.GetJSONNode(meta.HTTPResponse.Body)

	k.Storage[meta.TargetHash] = append(
		k.Storage[meta.TargetHash],
		Storage{
			Payload:        meta.Payload,
			HTMLNode:       httpnode.GetHTMLNode(meta.HTTPResponse.Body),
			HeaderNode:     httpnode.GetHeaderNode(meta.HTTPResponse.Headers),
			JSONNode:       jsonNode,
			ExtractBody:    meta.ExtractResultBody,
			ExtractHeaders: meta.ExtractResultHeader,
			Response:       meta.HTTPResponse,
			//Request: result.job.,
			//HttpReflectSurroundings: ,
		},
	)
}
