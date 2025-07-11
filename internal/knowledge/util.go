package knowledge

import (
	"reflect"

	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/httpreflect"
)

func GetKnowledge(storageMap map[string][]Storage) map[string]MergedKnowledge {
	var storedKnowledge = make(map[string]MergedKnowledge)
	merge := NewMerge()

	for hashId, storage := range storageMap {
		mk := MergedKnowledge{}
		// Json merge
		mk.Merge.JSONNode = httpnode.NewMergeJSONNode()
		for _, s := range storage {
			mk.Canary = s.Payload
			//k.Requests = append(k.Requests, d.Request)
			mk.Responses = append(mk.Responses, s.Response)
			mk.HttpReflectSurroundings = httpreflect.MergeUniqueSurroundings(mk.HttpReflectSurroundings, s.HttpReflectSurroundings)

			// Merge the extract result
			mk.Merge.ExtractBody = mergeAppendMaps(reflect.ValueOf(&merge.ExtractBody), s.ExtractBody).(extract.MergeExtractResult)
			mk.Merge.ExtractHeader = mergeAppendMaps(reflect.ValueOf(&merge.ExtractHeader), s.ExtractHeaders).(extract.MergeExtractResult)

			// Merge the HTTP response headers
			mk.Merge.HeaderNode = merge.HeaderNode.Merge(s.Response.Headers)

			// Detect response body format and merge HTTP response body
			switch httpnode.DetectResponseType(s.Response.Body) {
			case httpnode.ResponseTypeJSON:
				mk.Merge.JSONNode.Merge(s.JSONNode)
			// HTML, XML, Plain-text, Binary
			default:
				mk.Merge.HTMLNode = mergeAppendMaps(reflect.ValueOf(&merge.HTMLNode), s.HTMLNode).(httpnode.HTMLNodeMerge)
			}

		}
		storedKnowledge[hashId] = mk
	}

	return storedKnowledge
}

// Append a string to a list Works similar as append but do not append duplicates or empty strings
func appendUniqueInt(l []int, i int) []int {
	if len(l) == 0 {
		return append(l, i)
	}
	for _, item := range l {
		if item == i {
			return l
		}
	}
	return append(l, i)
}
