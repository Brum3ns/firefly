package knowledge

import (
	"reflect"

	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/httpreflect"
)

func GetKnowledge(storage map[string][]Storage) map[string]MergedKnowledge {
	var storedKnowledge = make(map[string]MergedKnowledge)
	c := NewMerge()

	for hashId, data := range storage {
		k := MergedKnowledge{}
		for _, d := range data {
			k.Canary = d.Payload
			//k.Requests = append(k.Requests, d.Request)
			k.Responses = append(k.Responses, d.Response)
			k.HttpReflectSurroundings = httpreflect.MergeUniqueSurroundings(k.HttpReflectSurroundings, d.HttpReflectSurroundings)

			k.Merge.HeaderNode = c.HeaderNode.Merge(d.Response.Headers)
			k.Merge.Extract = mergeAppendMaps(reflect.ValueOf(&c.Extract), d.Extract).(extract.ResultCombine)
			k.Merge.HTMLNode = mergeAppendMaps(reflect.ValueOf(&c.HTMLNode), d.HTMLNode).(httpnode.HTMLNodeCombine)
		}
		storedKnowledge[hashId] = k
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
