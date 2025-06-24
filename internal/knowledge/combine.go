package knowledge

import (
	"reflect"

	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/httpnode"
)

type Merge struct {
	Extract    extract.ResultCombine
	HTMLNode   httpnode.HTMLNodeCombine
	HeaderNode httpnode.HeaderNode
}

func NewMerge() Merge {
	return Merge{
		Extract:    extract.NewCombine(),
		HTMLNode:   httpnode.NewCombineHTMLNode(),
		HeaderNode: httpnode.NewHeader(),
	}
}

// Take a structure and combine all "map[string]int" into a map[string][]int and return the combined map:
func mergeAppendMaps(mergeData reflect.Value, data any) interface{} {
	mergeData = mergeData.Elem()
	dataValue := reflect.ValueOf(data)
	t := dataValue.Type()

	//Extract all field from the given "data":
	for i := 0; i < dataValue.NumField(); i++ {
		data_field := dataValue.Field(i)
		data_name := t.Field(i).Name

		//In case the field is a correct map that can be combined, then procceed:
		if data_map, ok := data_field.Interface().(map[string]int); ok {

			//Extract the same field (by name) from "combineData" that was recently extracted from "data":
			mergeData_field := mergeData.FieldByName(data_name)

			//Make sure the "cData" field is a correct map that can be used to compare the original map from "data":
			if mergeData_map, ok := mergeData_field.Interface().(map[string][]int); ok {

				//Extract the key value and the key's value. Then add only the unique items from "newData" to "combineData"
				for k, v := range data_map {
					mergeData_map[k] = appendUniqueInt(mergeData_map[k], v)
				}
			}
		}
	}
	return mergeData.Interface()
}
