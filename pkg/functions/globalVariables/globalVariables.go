package globalvariables

import (
	"fmt"
	"os"
)

/*Global variables [static] - No dynamic variables!
Major of the variables are declared in the validation process in 'options.go'.*/
var (
	//Request
	Verify         int
	TotalReq       int
	TotalVReq      int
	Insert         string
	PayloadPattern string
	PayloadChars   = make(map[rune]string)
	DefaultProto   = "http"

	//Verify

	//Task in total Verify/fuzz process:
	Total int

	//Random
	M_Random      = DefaultRandomKeyword()
	RandomOptions string

	//Input and output:
	Pipe bool

	//Payload:
	EncodeChar  string
	PayloadMark = "__FIREFLY_PAYLOAD__"

	//Verbose:
	Verbose bool

	//Input:
	Lst_rawData map[string][]string

	//Output:
	Output       bool
	OutputType   string
	OutputFile   string
	OutputFileOS *os.File

	//Filter/Match:
	Lst_FMCheck = []string{"sc", "lc", "wc", "bs"}
	Lst_mFilter map[string][]string
	Lst_mMatch  map[string][]string
	MF_Mode     int

	//Filter/Match regex:
	/*Lst_mRegexFilter = make(map[string]string)//DELETE? - old filter maps
	Lst_mRegexMatch  = make(map[string]string)*/

	Lst_CheckHeaders []string
	Lst_RandomAgent  []string
	RegexHeaders     bool

	//Regex
	Regex string

	//Counters
	Amount_Lst  int
	Amount_LstG int
	Amount_Item int
)

func DefaultRandomKeyword() map[string]int {
	//Setup default first:
	m := make(map[string]int)
	m["s"] = 8
	m["n"] = 8

	return m
}

//Debug functions (I'm lazy)
func PL(l []string) {
	for n, i := range l {
		fmt.Println("DEBUG lst:", n, "|", i)
	}
}
