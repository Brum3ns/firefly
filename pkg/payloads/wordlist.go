package payloads

import (
	"bufio"
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/Brum3ns/firefly/pkg/encode"
)

// Wordlist global tag names:
var (
	TAG_VERIFY         = "Verify"
	TAG_FUZZ           = "Fuzz"
	TAG_TRANSFORMATION = "Transformation"
	TAGS               = []string{TAG_VERIFY, TAG_FUZZ, TAG_TRANSFORMATION}
)

// Wordlist structure stores the wordlist and tags
type Wordlist struct {
	Wordlist           map[string][]string //(tag|wordlist)
	Files              []string
	TransformationList []string
	Verify             Verify
	PayloadProperties
}

type PayloadProperties struct {
	Tamper         string
	Encode         []string
	PayloadReplace string
	PayloadPattern string
	PayloadSuffix  string
	PayloadPrefix  string
}

type Verify struct {
	Payload string
	Amount  int
}

// Create a new wordlist object
func NewWordlist(wl *Wordlist) *Wordlist {
	wl.Wordlist = make(map[string][]string)

	//Create verify wordlist
	wl.Wordlist[TAG_VERIFY] = verifyWordlist(wl.Verify.Payload, wl.Verify.Amount)

	//Create fuzz wordlist by combining all wordlist files given (if multiple)
	for _, filename := range wl.Files {
		wl.Wordlist[TAG_FUZZ] = append(wl.Wordlist[TAG_FUZZ], wl.createPayloadWordlist(filename)...)
	}

	//Transformation wordlist:
	wl.Wordlist[TAG_TRANSFORMATION] = wl.TransformationList

	return wl
}

// Get a wordlist by tag
func (wl *Wordlist) Get(tag string) ([]string, error) {
	if wl, exist := wl.Wordlist[tag]; !exist {
		return wl, errors.New("can't find the tag")
	} else {
		return wl, nil
	}
}

// Return a map containing all the wordlists and tags (tag as the key)
func (wl *Wordlist) GetAll() map[string][]string {
	return wl.Wordlist
}

func verifyWordlist(verifyPayload string, amount int) []string {
	var lst []string
	for i := 0; i < amount; i++ {
		lst = append(lst, verifyPayload)
	}
	return lst
}

// Create a wordlist by a given file path
func (wl Wordlist) createWordlist(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	var (
		lst     []string
		scanner = bufio.NewScanner(file)
	)
	for scanner.Scan() {
		item := scanner.Text()
		if len(item) > 0 {
			lst = append(lst, item)
		}
	}
	file.Close()
	return lst
}

// Take a filename and return a payload adapted wordlist by using given rules in the payloadProperties [struct]ure:
func (wl Wordlist) createPayloadWordlist(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
	}
	var (
		lst     []string
		scanner = bufio.NewScanner(file)
	)
	for scanner.Scan() {
		payload := scanner.Text()
		if len(payload) > 0 {

			if len(wl.PayloadReplace) > 0 {
				payload = replaceRegex(payload, wl.PayloadReplace)
			}

			//Check if payload should be encoded:
			if len(wl.Encode) > 0 {
				payload = encode.Encode(payload, wl.Encode)
			}
			payload = (wl.PayloadPrefix + payload + wl.PayloadSuffix)
			lst = append(lst, payload)
		}
	}
	file.Close()
	return lst
}

func replaceRegex(p, regexReplace string) string {
	i := strings.Split(regexReplace, " => ")
	re := regexp.MustCompile(i[0])
	return re.ReplaceAllString(p, i[1])
}
