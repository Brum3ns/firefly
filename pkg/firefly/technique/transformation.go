package technique

import (
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/Brum3ns/firefly/pkg/functions"
	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
	"github.com/Brum3ns/firefly/pkg/storage"
	"gopkg.in/yaml.v2"
)

type transformationData struct {
	Transformation []string `yaml:"transformations"`
	Payloads       []string `yaml:"payloads"`
}

type result struct {
	m_transformation map[string]string
	s_body           string
	s_regx           string
	s_payload        string
	s_payloadCompare string
}

/**Scan the fuzzed respones for any payload transformation and return the result with the original payload and the formated result (if any)*/
func TransformationScan(tfmtPayloads map[string]string, resp storage.Response) (map[string]string, bool) {
	payloadClear := functions.PayloadClearPattern(resp.Payload)
	data := &result{
		m_transformation: make(map[string]string),
		s_body:           string(resp.Body[:]),
		s_regx:           strings.Replace(regexp.QuoteMeta(functions.PayloadPattern("__REPALCE__")), "__REPALCE__", "*(.*?)*", 1),
		s_payload:        payloadClear,
		s_payloadCompare: tfmtPayloads[payloadClear],
	}

	//If the payload isen't included in the transformation payload list:
	if len(tfmtPayloads[data.s_payload]) <= 0 {
		return nil, false
	}

	//Setup regex to extract position of the payload reflectation:
	regx := regexp.MustCompile(data.s_regx)
	positions := regx.FindAllStringIndex(data.s_body, -1)

	for _, l_pos := range positions {
		if len(l_pos)%2 != 0 {
			//[TODO] log error to file
			return nil, false

		}

		for p1, p2 := 0, 1; p1 < len(l_pos); p1, p2 = p1+2, p2+2 {

			payloadExtract := functions.PayloadClearPattern(data.s_body[l_pos[p1]:l_pos[p2]])

			if tfmtPayloads[data.s_payload] == payloadExtract {
				data.m_transformation[data.s_payload] = payloadExtract
				return data.m_transformation, true
			}
		}
	}

	return nil, false
}

/**Setup the yml file that includes the payload and the transformations to be expected*/
func Transformation(ymlFile string) ([]string, map[string]string) {
	var (
		l_payloads []string
		m_result   = make(map[string]string)
		tfmt       transformationData
	)

	//Read the yaml file:
	data, err := ioutil.ReadFile(ymlFile)
	functions.IFError("f", err)

	//Setup yaml data to struct:
	err = yaml.Unmarshal(data, &tfmt)
	functions.IFError("f", err)

	if len(tfmt.Payloads) != len(tfmt.Transformation) {
		return nil, nil
	}

	for i := 0; i < len(tfmt.Payloads); i++ {
		l_payloads = append(l_payloads, tfmt.Payloads[i])
		m_result[tfmt.Payloads[i]] = tfmt.Transformation[i]
	}

	G.Amount_Item += len(l_payloads)

	return l_payloads, m_result
}
