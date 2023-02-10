package output

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	fc "github.com/Brum3ns/FireFly/pkg/functions"
	G "github.com/Brum3ns/FireFly/pkg/functions/globalVariables"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

type OutputData struct {
	data storage.Collection `json:"ID"`
}

func Output(clt storage.Collection) error {
	out := &OutputData{
		data: clt,
	}

	m_data := OutputTemplate(out)

	m_out := make(map[int]map[string]interface{})
	m_out[clt.ID] = make(map[string]interface{})
	m_out[clt.ID] = m_data

	var output string
	switch strings.ToLower(G.OutputType) {
	case "json":

		// outJson, err := json.MarshalIndent(out.data, "", "  ")
		outJson, _ := json.Marshal(m_out)
		output = string(outJson)

	case "output":
		l := []string{}
		for k, i := range m_data {
			l = append(l, fmt.Sprintf("%s: %v", k, i)) //DEBUG
		}
		output = strings.Join(l, "\n")
	}
	fmt.Fprintf(G.OutputFileOS, "%s\n", output)

	return nil
}

/**Setup what data that should be used within the output file*/
func OutputTemplate(o *OutputData) map[string]interface{} {
	var (
		l_ignore = []string{
			"body",
			"icon",
			"errors",
			"haserr",
			"threadid",
			"valid",
			"status",
			"errmsg",
			"tfmt_display",
			"taskstatus",
			"tfmt_ok",
			"filter",

			//[TODO] - Below will be added in later version(s)
			"resperr",
			"headersmatch",
			"respdiff",
			"regexmatch",
		}

		v    = reflect.ValueOf(o.data)
		data = make(map[string]interface{}, v.NumField())
		typ  = v.Type()
	)

	for i := 0; i < v.NumField(); i++ {
		N := strings.ToLower(typ.Field(i).Name)

		//If it's not in the 'ignore list' then add it to the output template:
		if !fc.InLst(l_ignore, N) {
			var V interface{}
			V = v.Field(i).Interface()

			//[TODO] - Add options to be able to save full bodies + headers
			/* if N == "body" {
				V = string(v.Field(i).Interface().([]byte))
			} else {
				V = v.Field(i).Interface()
			}
			*/
			data[typ.Field(i).Name] = V
		}
	}

	return data
}
