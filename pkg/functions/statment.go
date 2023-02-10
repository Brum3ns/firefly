package functions

import (
	"fmt"
	"log"
	"os"
	"strings"

	d "github.com/Brum3ns/firefly/pkg/design"
	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
)

/** If there is an error then check what type to provide*/
func IFError(typ string, err error, msg ...string) (bool, string) {

	if err != nil {
		//Verbose error message
		var ErrMsg = fmt.Sprintf("%s %s", d.Fail, err)

		switch typ {
		case "f":
			log.Fatal(ErrMsg)

		case "p":
			log.Panic(ErrMsg)
		}

		if G.Verbose {
			fmt.Println(d.Fail, typ, err)
		}
		return true, typ
	}

	//No error was detected
	return false, typ
}

/** Check option error type Invalid user input detected and return msg*/
func IFFail(s string) string {

	var (
		msg  string
		exit = false
	)
	t := strings.Split(s, ":")

	switch t[0] {
	case "stdin":
		msg = fmt.Sprintf("%s Can't use %s with %s option\n", d.Fail, d.Colortxt("stdin", "o", false), d.Colortxt("-f", "o", false))

	case "v":
		switch t[1] {
		case "mp":
			msg = fmt.Sprintf("%s The verified method and protocol must both contain at least one value that is the same as the verified method and protocol specified. %s, %s\n", d.Fail, d.Colortxt("-vm", "o", false), d.Colortxt("-vp", "o", false))
			exit = true
		}
	case "-u":
		msg = fmt.Sprintf("%s No input was detected %s, %s or stdin\n", d.Fail, d.Colortxt("-u", "o", false), d.Colortxt("-f", "o", false))
	case "-r":
		msg = d.Fail + " No host header in " + d.Colortxt("-a", "o", false)
	case "-i":
		msg = d.Fail + " No insert points detected " + d.Colortxt("-i", "o", false)
	case "-w":
		msg = d.Fail + " Wordlist do not contain \":\". " + d.Colortxt("-w", "o", false) + " Ex: (-w wordlist.txt:fuzz)"
		exit = true
	case "-pr":
		msg = d.Fail + " The argument \"payload-replace\" " + d.Colortxt("-pr", "o", false) + "do not contain the \" => \" (spaces included). FireFly dosen't know what to replace the regex/string with."
	case "-t":
		msg = d.Fail + " Can't use a threads lower or equal to \"[0]\" " + d.Colortxt("-t", "o", false)
	case "-ov":
		msg = (d.Warning + " This file already exist. If you want to overwrite it. Use option " + d.Colortxt("-ov", "o", false))
	case "-a":
		msg = d.Fail + " Invalid attack technique detected " + d.Colortxt("-a", "o", false)
	case "-au":
		msg = d.Fail + " Invalid input for \"auto-detect\" " + d.Colortxt("-au", "o", false)
		exit = true
	case "-ra":
		msg = d.Fail + " Invalid random value Ex usage: s:8 (string with length as 8) " + d.Colortxt("Random: Invalid usage", "r", false)
		exit = true

	case "-vfc":
		msg = d.Fail + " Can't setup the payload given " + d.Colortxt("Invalid -verify-char", "r", false)
		exit = true
	case "vresp":
		msg = d.Fail + " Could not " + d.Colortxt("verify default response behaviors", "o", false)
	case "vreq":
		msg = d.Fail + " The verify responses are less than 50% in success." + d.Colortxt("There was to many request/response errors.", "r", false)
		exit = true

	case "tfmt":
		msg = d.Fail + " The amount of payloads do not match the amount of transformations " + d.Colortxt("The amount must be equal to be comparable", "r", false)
		exit = true

	case "eiID": //critical
		msg = d.Critical + " The engine tasks mixed up processes & RespID's " + d.Colortxt("Engine: response ID verification", "r", false)
	case "eiLen": //critical
		msg = d.Critical + " The engine wasen't able to proceed all tasks " + d.Colortxt("Engine: Not all tasks was done related to a response (id)", "r", false)

	case "taskcheck": //critical
		msg = d.Critical + "The task process inside the engine failed. " + d.Colortxt("Task result collided with the struct memory", "r", false)
		exit = true

	case "eiStatus": //critical
		msg = d.Critical + " The engine proceed all tasks but the result for one or more failed to return " + d.Colortxt("Engine: All task completed but not all where had a valid status", "r", false)
		exit = true

	case "runner": //critical
		msg = d.Critical + d.Colortxt(" The runner core process had an error in the process. Make sure that no configure is invalid or contains a syntax error", "r", false)
		exit = true

	case "vrunner": //critical
		msg = d.Critical + d.Colortxt(" The runner verify process had an error in the process. Make sure that no configure is invalid or contains a syntax error", "r", false)
		exit = true

	case "filter": //critical
		msg = d.Critical + d.Colortxt(" The runner filter/match failed to detect channel type.", "r", false)
		exit = true

	default:
		msg = ""
	}

	//[TODO] If the future process is in need of the failed input, exit.
	if exit {
		fmt.Println(msg)
		os.Exit(1)
	}

	return msg
}
