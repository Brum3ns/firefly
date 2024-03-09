package knowledge

import (
	"reflect"

	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/prepare"
)

type Knowledge struct {
	PayloadVerify string
	Responses     []output.Response
	Requests      []output.Request
	Combine       Combine
}

type Combine struct {
	Extract  extract.ResultCombine
	HTMLNode prepare.HTMLNodeCombine
}

type Learnt struct {
	Payload  string
	HTMLNode prepare.HTMLNode
	Extract  extract.Result

	Response output.Response
	Request  output.Request
}

func NewKnowledge() *Knowledge {
	return &Knowledge{}
}

func NewCombine() Combine {
	return Combine{
		Extract:  extract.NewCombine(),
		HTMLNode: prepare.NewCombineHTMLNode(),
	}
}

func GetKnowledge(learnt map[string][]Learnt) map[string]Knowledge {
	var storedKnowledge = make(map[string]Knowledge)
	c := NewCombine()

	for hashId, data := range learnt {
		k := Knowledge{}
		for _, d := range data {
			k.PayloadVerify = d.Payload
			k.Requests = append(k.Requests, d.Request)
			k.Responses = append(k.Responses, d.Response)

			k.Combine.Extract = combineAppendMaps(reflect.ValueOf(&c.Extract), d.Extract).(extract.ResultCombine)
			k.Combine.HTMLNode = combineAppendMaps(reflect.ValueOf(&c.HTMLNode), d.HTMLNode).(prepare.HTMLNodeCombine)
		}
		storedKnowledge[hashId] = k
	}
	return storedKnowledge
}

// Take a structure and combine all "map[string]int" into a map[string][]int and return the combined map:
func combineAppendMaps(combineData reflect.Value, data any) interface{} {
	combineData = combineData.Elem()
	dataValue := reflect.ValueOf(data)
	t := dataValue.Type()

	//Extract all field from the given "data":
	for i := 0; i < dataValue.NumField(); i++ {
		data_field := dataValue.Field(i)
		data_name := t.Field(i).Name

		//In case the field is a correct map that can be combined, then procceed:
		if data_map, ok := data_field.Interface().(map[string]int); ok {

			//Extract the same field (by name) from "combineData" that was recently extracted from "data":
			combineData_field := combineData.FieldByName(data_name)

			//Make sure the "cData" field is a correct map that can be used to compare the original map from "data":
			if combineData_map, ok := combineData_field.Interface().(map[string][]int); ok {

				//Extract the key value and the key's value. Then add only the unique items from "newData" to "combineData"
				for k, v := range data_map {
					combineData_map[k] = appendUniqueInt(combineData_map[k], v)
				}
			}
		}
	}
	return combineData.Interface()
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
