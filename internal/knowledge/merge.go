package knowledge

import (
	"reflect"

	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpnode"
	"github.com/Brum3ns/firefly/pkg/httpreflect"
	"github.com/Brum3ns/firefly/pkg/rhttp"
)

type MergedKnowledge struct {
	Canary    string
	Responses []rhttp.Response
	//Requests                []output.Request
	HttpReflectSurroundings []httpreflect.Surrounding
	Merge                   Merge
}

type Merge struct {
	//Extract    extract.ResultMerge
	HTMLNode      httpnode.HTMLNodeMerge
	JSONNode      httpnode.JSONNodeMerge
	HeaderNode    httpnode.HeaderMergeNode
	ExtractBody   extract.MergeExtractResult
	ExtractHeader extract.MergeExtractResult
}

func NewMerge() Merge {
	return Merge{
		ExtractHeader: extract.NewMerge(),
		ExtractBody:   extract.NewMerge(),
		HeaderNode:    httpnode.NewHeader(),
		HTMLNode:      httpnode.NewMergeHTMLNode(),
		JSONNode:      httpnode.NewMergeJSONNode(),
	}
}

// Take a structure and combine all "map[string]int" into a map[string][]int and return the combined map:
func mergeAppendMaps(mergeData reflect.Value, data any) interface{} {
	mergeData = mergeData.Elem()
	dataValue := reflect.ValueOf(data)
	t := dataValue.Type()

	//Extract all field from the given "data":
	for i := 0; i < dataValue.NumField(); i++ {
		dataField := dataValue.Field(i)
		dataName := t.Field(i).Name

		//In case the field is a correct map that can be combined, then procceed:
		if data, ok := dataField.Interface().(map[string]int); ok {

			//Extract the same field (by name) from "combineData" that was recently extracted from "data":
			mergeDataField := mergeData.FieldByName(dataName)

			//Make sure the "cData" field is a correct map that can be used to compare the original map from "data":
			if mergeData, ok := mergeDataField.Interface().(map[string][]int); ok {

				//Extract the key value and the key's value. Then add only the unique items from "newData" to "combineData"
				for k, v := range data {
					mergeData[k] = appendUniqueInt(mergeData[k], v)
				}
			}
		}
	}
	return mergeData.Interface()
}
