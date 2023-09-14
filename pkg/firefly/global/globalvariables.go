// Global variables - (no dynamic variables)
package global

import (
	"os"
)

// File directory global variables (Never ever reassign any of these):
var (
	//Root directories:
	DIR_HOME, _ = os.UserHomeDir()
	DIR_CONFIG  = (DIR_HOME + "/.config/firefly/")
	DIR_DB      = (DIR_CONFIG + "/db/")

	//Resource directories:
	DIR_TAMPERS   = (DIR_DB + "/tampers/")
	DIR_RESOURCE  = (DIR_DB + "/resources/")
	DIR_DETECTION = (DIR_DB + "/resources/detection/")
	DIR_WORDLIST  = (DIR_DB + "/wordlists/")

	//Resource files:
	FILE_RANDOMAGENT    = (DIR_RESOURCE + "/randomUserAgent.txt")
	FILE_SKIP_HEADERS   = (DIR_RESOURCE + "/skipheaders.txt")
	FILE_TRANSFORMATION = (DIR_RESOURCE + "/transformation.yml")
)

// Major of the variables are declared in the validation process in 'options.go'.
var (

	//Request
	VERIFY          int
	INSERT          string
	PAYLOAD_PATTERN string
	PayloadChars    = make(map[rune]string)
	DefaultProto    = "http"

	RANDOMNESS_WHITELIST = lettersDigits()
	CONSONANTS_DIGITS    = constantsDigits()

	//Random
	RANDOM_INSERT = make(map[string]int)
	RANDOM_OPTION string

	//Input and output:
	Pipe bool

	//Payload:
	PayloadMark = "__FIREFLY_PAYLOAD__"

	//Verbose:
	DEBUG   bool
	VERBOSE bool

	//Input:
	Lst_rawData map[string][]string

	//Output:
	OutputType   string
	OutputFile   string
	OutputFileOS *os.File

	Total int //DELETE to save RAM memory

	CHECK_HEADERS []string
	RANDOM_AGENTS []string

	//Amount
	AMOUNT_ITEM int

	//Verification tags
	TAG_VERIFYPAYLOAD = "verifypayload"
	TAG_VERIFYCHAR    = "verifychar"
)

// Make and return a map containing rune and string of character [a-zA-Z0-9]
func lettersDigits() map[rune]string {
	var m = make(map[rune]string)
	for az, AZ, O9 := 'a', 'A', 48; az <= 'z' && AZ <= 'Z'; az, AZ, O9 = (az + 1), (AZ + 1), (O9 + 1) {
		if O9 <= 57 {
			v := rune(O9)
			m[v] = string(v)
		}
		m[az], m[AZ] = string(az), string(AZ)
	}
	//Add some special characters that is common in tokens and or hashes (URL encoded and or separeted for 'x' length etc...)
	//( _, -, :, ;, %, . =)
	for _, rn := range []rune{95, 45, 58, 59, 37, 46, 61} {
		m[rn] = string(rn)
	}
	return m
}

func constantsDigits() map[rune]string {
	var (
		m = make(map[rune]string) //Return
		l = []rune{
			'b', 'c', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'y', 'z',
			'B', 'C', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X', 'Y', 'Z',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		}
	)
	for _, rn := range l {
		m[rn] = string(rn)
	}
	return m
}
