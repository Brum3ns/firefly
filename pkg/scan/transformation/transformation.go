package transformation

import (
	"errors"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/Brum3ns/firefly/pkg/design"
	"gopkg.in/yaml.v2"
)

var (
	PREFIX = "1229"
	SUFFIX = "1345"
)

type Transformation struct {
	Storage map[string][2]string //(Expected payload[Transformed payload|Description])
	Regex   *regexp.Regexp
}

type Properties struct {
}

type Result struct {
	OK      bool
	Desc    string
	Payload string
	Format  string
}

// Create a new transformation
func NewTransformation(yamlFile string) Transformation {
	storage, err := getYamlMap(yamlFile)
	if err != nil {
		log.Println(design.STATUS.ERROR, err)
	}
	return Transformation{
		Storage: storage,
		Regex:   regexp.MustCompile(PREFIX + `(.*?)` + SUFFIX),
	}
}

func (t Transformation) Detect(body, payload string) Result {
	reflectedPayloads := t.Regex.FindAllString(body, -1)
	if reflectedPayloads == nil {
		return Result{OK: false}
	}

	//Search for a payload transformation:s
	payload = RmPrefixSuffix(payload)
	if arr, ok := t.Storage[payload]; ok {
		expectedPayload := arr[0]
		desc := arr[1]

		//Check if any valid transformation was discovered from all the reflected payload patterns:
		for _, i := range reflectedPayloads {
			transformationPayload := RmPrefixSuffix(i)
			if transformationPayload == expectedPayload {
				return Result{
					OK:      true,
					Desc:    desc,
					Payload: payload,
					Format:  transformationPayload,
				}
			}
		}
	}
	return Result{OK: false}
}

func RmPrefixSuffix(s string) string {
	return s[len(PREFIX):(len(s) - len(SUFFIX))]
}
func GetWordlist(yamlFile string) []string {
	var wordlist []string

	m, err := getYamlMap(yamlFile)
	if err != nil {
		log.Fatal(design.STATUS.ERROR, err)
	}
	for payload, _ := range m {
		wordlist = append(wordlist, (PREFIX + payload + SUFFIX))
	}
	return wordlist
}

// Read the transformation yaml file and return the full payloads, payload transformation map or error (if any)
func getYamlMap(yamlFile string) (map[string][2]string, error) {

	var storage = make(map[string][2]string)

	// Read the YAML file.
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return storage, errors.New("error reading transformation yaml file that was given")
	}

	// Unmarshal the YAML data into a map.
	tmpMap := make(map[string][][2]string)
	if err := yaml.Unmarshal(data, &tmpMap); err != nil {
		return storage, errors.New("error when unmarshaling the given yaml file, make sure that the yaml file do not contain any syntax errors")
	}

	for expectedPayload, lst := range tmpMap {
		for _, arr := range lst {

			if len(arr) != 2 {
				return storage, errors.New("invalid yaml syntax. The list must contain a paylaod and a desciption.")
			}
			payload := arr[0]
			desc := arr[1]

			storage[payload] = [2]string{expectedPayload, desc}
		}
	}

	return storage, nil
}
