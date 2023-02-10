package technique

import (
	"strings"

	"github.com/Brum3ns/FireFly/pkg/design"
	"github.com/Brum3ns/FireFly/pkg/functions"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

//[TODO] Fix XSS detection better (print if an XSS as a high chance to work)

//Detect reflected values and check how high the possibility is for a reflected XSS
func Reflect(result storage.Response) (string, string) {
	var lst_xssChar = []string{"'", "\"", "<", ">", "\\", "</"}
	var color, IconXSS, status string

	if strings.Contains(string(result.Body), result.Payload) {

		//Check how likely it is to be a [Reflected XSS]:
		color = Detect_reflectedXSS(result.Payload, lst_xssChar)

		//Check how many "chars" of "lst_xssChar" was reflected by the target back:
		if color == "0" {
			IconXSS = "---"
		} else {
			IconXSS = strings.Replace(design.XSS, "_COLOR_", color, 1)

		}
		status = design.Reflect + "|" + IconXSS

		//Else give verbose back that it only reflect without any detection of possible XSS:
		return design.Success, status

	} else {
		return design.Info, design.Null
	}
}

//Move this and structure it better next: [TODO]
var Attack_techq []string

//Give back -> Vulnerable, vulnerable response data:

func Detect_reflectedXSS(str string, lst []string) string {
	var (
		lst_char = strings.Split(str, "")

		color = ""
		chars = ""
		hit   = 0
	)

	if functions.CheckSubstrings(str, "\"'<\\/") == 6 || (functions.CheckSubstrings(str, "\"'") == 2 && !strings.Contains(str, "\\")) {
		hit = 4

	} else if (functions.CheckSubstrings(str, "\"<") == 2 || functions.CheckSubstrings(str, "'<") == 2) && !strings.Contains(str, "\\") {
		hit = 4

	} else if (strings.Contains(str, "</") && !strings.Contains(str, "\\/")) || (functions.CheckSubstrings(str, "\"<") == 2 || functions.CheckSubstrings(str, "'<") == 2) {
		hit = 3

	} else if strings.Contains(str, "\\\\\"") || strings.Contains(str, "\\\\'") {
		hit = 2

	} else {
		for _, item := range lst {
			for _, char := range lst_char {

				if char == item {
					chars += char
					hit++
					break
				}
			}
		}
	}

	if hit > 1 {
		switch hit {
		case 2:
			color = "130;44"

		case 3:
			color = "130;43"

		default: // 4 >
			color = "130;41"
		}
	} else {
		color = "0"
	}

	return color
}
