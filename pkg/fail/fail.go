package fail

import (
	"fmt"
	"log"

	"github.com/Brum3ns/firefly/pkg/design"
	globalvariables "github.com/Brum3ns/firefly/pkg/firefly/global"
)

// Failed messages:
var ERRORCODE_MESSAGES = map[int]string{
	13001:  design.STATUS.FAIL + " The verify responses are less than 50% in success (" + design.COLOR.RED + "There was to many request/response errors." + design.COLOR.WHITE + ")",
	1011:   design.STATUS.FAIL + " The filter syntax is invalid. The valid for each filter separeted by a comma (if any) are as following (not combined): \".\",\"-\",\"--\",\"++\"",
	1008:   design.STATUS.FAIL + " No input was detected (" + design.COLOR.ORANGE + "-u" + design.COLOR.WHITE + "," + design.COLOR.ORANGE + "-f" + design.COLOR.WHITE + ") or STDIN pipeline",
	1001:   design.STATUS.FAIL + " Invalid HTTP Raw data" + design.COLOR.ORANGE + "-r" + design.COLOR.WHITE + ")",
	10005:  design.STATUS.FAIL + " No insert points detected (" + design.COLOR.ORANGE + "-i" + design.COLOR.WHITE + ")",
	8001:   design.STATUS.FAIL + " The argument \"payload-replace\" (" + design.COLOR.ORANGE + "-pr" + design.COLOR.WHITE + ") do not contain the \" => \" (spaces included). Firefly dosen't know what to replace the regex/string with.",
	1006:   design.STATUS.FAIL + " Can't use a threads lower or equal to zero (" + design.COLOR.ORANGE + "-t" + design.COLOR.WHITE + ")",
	1005:   design.STATUS.WARNING + " This file already exist. If you want to overwrite it. Use option (" + design.COLOR.ORANGE + "-ov" + design.COLOR.WHITE + ")",
	10009:  design.STATUS.FAIL + " Invalid input for \"auto-detect\" (" + design.COLOR.ORANGE + "-au" + design.COLOR.WHITE + ")",
	100014: design.STATUS.FAIL + " Invalid random value Example usage: s8 (string with length as 8) or 's4,n8' to use both string and number(" + design.COLOR.RED + "Random: Invalid usage" + design.COLOR.WHITE + ")",
	13003:  design.STATUS.FAIL + " Can't setup the payload given (" + design.COLOR.ORANGE + "-verify-char" + design.COLOR.WHITE + ")",
	2001:   design.STATUS.FAIL + " The level has to be between 1-3 (" + design.COLOR.ORANGE + "-lv" + design.COLOR.WHITE + ")",
	3001:   design.STATUS.FAIL + " The match mode is invalid (" + design.COLOR.ORANGE + "-mmode" + design.COLOR.WHITE + "). Valid input: and, or",
	3010:   design.STATUS.FAIL + " The filter mode is invalid (" + design.COLOR.ORANGE + "-fmode" + design.COLOR.WHITE + "). Valid input: and, or",
	9003:   design.STATUS.FAIL + " Invalid transformation input",
	10001:  design.STATUS.FAIL + " Invalid URL(s) given (" + design.COLOR.ORANGE + "-u" + design.COLOR.WHITE + ")",
	10002:  design.STATUS.FAIL + " Invalid method(s) given (" + design.COLOR.ORANGE + "-X" + design.COLOR.WHITE + ")",
	10003:  design.STATUS.FAIL + " Invalid scheme(s) input (" + design.COLOR.ORANGE + "-scheme" + design.COLOR.WHITE + ")",
	9001:   design.STATUS.FAIL + " Invalid wordlist given. Make sure that the wordlist is not empty (" + design.COLOR.ORANGE + "-w" + design.COLOR.WHITE + ").",
	9002:   design.STATUS.FAIL + " Invalid wordlist folder given. Firefly coulen't find atleast one valid file (wordlist to use) in the folder. Make sure that the files in the folder are correct set and not empty (" + design.COLOR.ORANGE + "-wf" + design.COLOR.WHITE + ").",
	10012:  design.STATUS.FAIL + " Cannot use a timeout lower than zero",
	4001:   design.STATUS.FAIL + " The specified output file already exists. Use the overwrite option to overwrite it (be careful).",
}

// Check the failed type
// Return a desciption message on why it failed
func IFFail(code int) {
	if code <= 0 {
		return
	} else if msg, ok := ERRORCODE_MESSAGES[code]; ok {
		log.Fatal(msg)
	} else {
		log.Fatal(design.STATUS.FAIL + " Unkown failure!")
	}
}

// If there is an error then check what type to provide
func IFError(typ string, err error) (bool, string) {
	if err != nil {
		//Verbose error message
		var ErrMsg = fmt.Sprintf("%s %s", design.STATUS.FAIL, err)

		switch typ {
		case "f":
			log.Fatal(ErrMsg)

		case "p":
			log.Panic(ErrMsg)
		}

		if globalvariables.VERBOSE {
			fmt.Println(design.STATUS.FAIL, typ, err)
		}
		return true, typ
	}

	//No error was detected
	return false, typ
}
