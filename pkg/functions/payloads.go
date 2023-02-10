package functions

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"

	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
)

/**Remove known verified payload patterns from a string with*/
/* func PayloadClear(lst_payloadPatterns []string, str, replaceWith string) string {
	for _, payload := range lst_payloadPatterns {
		payload = PayloadPattern(payload)

		if strings.Contains(str, payload) {
			str = strings.ReplaceAll(str, payload, replaceWith)
		}
	}

	return str
} */

/**(QuoteMeta included) The amount of replace will be done "i" times (-1 resulting in all). Return the regex pattern used to extract reflected payloads*/
func PayloadRegexMark(i int) string {
	//QuoteMeta is used becuase the pattern can be added by user input:
	return strings.Replace(regexp.QuoteMeta(PayloadPattern("__REPALCE__")), "__REPALCE__", "*(.*?)*", i)
}

func PayloadPattern(p string) string {
	return (G.PayloadPattern + (p) + G.PayloadPattern)
}

/**Clear the payload pattern*/
func PayloadClearPattern(p string) string {
	ln := len(G.PayloadPattern)
	if (len(p) >= ln*2) && (strings.HasPrefix(p, G.PayloadPattern) && strings.HasSuffix(p, G.PayloadPattern)) {
		p = p[ln : len(p)-ln]
	}
	return p
}

func PayloadChar(r rune) string {
	return fmt.Sprintf("%d%d%dF%sF%d%d%d", r, r, r, string(r), r, r, r)
}

func PayloadURLNormalize(p string) string {
	var (
		l_find        = []string{" ", "\t", "\n", "#", "&", "?"}
		l_URLEncodeTo = []string{"%20", "%09", "%0a", "%23", "%26", "%3F"}
	)

	for i := 0; i < len(l_URLEncodeTo); i++ {
		if strings.Contains(p, l_find[i]) {
			p = strings.ReplaceAll(p, l_find[i], l_URLEncodeTo[i])
		}
	}
	return p
}

/**Insert the payload to a string*/
func PayloadInsert(s, p string) string {
	return strings.ReplaceAll(s, G.Insert, p)
}

//[TODO] <== Verify "PayloadEncode" function:
func PayloadEncode(p, encodeType string) string {

	/* [Payload encoder]
	* Encode payloads to adjust the fuzzing process.
	* The payloads can be encoded as many time as the user want and can be mixed
	* Smart encode only encode special chars that is likely to either make false
	* positives and will behave like a browser url normalization.
	 */

	//var lst_use []string

	//[TODO] Make user able to select chars ONLY be encoded:
	var lst_specialChars = []string{"!", "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", "\\", ":", ";", "<", "=", ">", "?", "@", "[", "\"", "]", "^", "_", "`", "{", "|", "}", "~"}
	var smart string

	var lst_encode = strings.Split(encodeType, ",")

	//Check if the encode was given "smart" option by the user input flag: //[TODO] - (Do not work)
	if strings.Contains(lst_encode[len(lst_encode)-1], ":") {
		lst_smart := strings.Split(lst_encode[len(lst_encode)-1], ":")
		smart = strings.ToLower(lst_smart[1])

		lst_encode[len(lst_encode)-1] = lst_smart[0]

		//Encode special chars by order:
		if smart == "smart" || smart == "s" {
			for _, char := range lst_specialChars {
				fmt.Print(char) // <===========================[START HERE]

				/* [TODO]
				- Smart function encode
				- Multi encode
				- Simple encode (single)
				*/
			}
		}
	}

	//Encode payload "x" amount of times depends on the user option:
	for _, encode := range lst_encode {

		if encode == "b64" {
			encode = "base64"
		}
		if encode == "b32" {
			encode = "base32"
		}

		switch strings.ToLower(encode) {
		case "url":
			p = url.QueryEscape(p)

		case "durl":
			p = strings.ReplaceAll(url.QueryEscape(p), "%", "%25")

		case "base64":
			p = base64.StdEncoding.EncodeToString([]byte(p))

		case "base32":
			p = base32.StdEncoding.EncodeToString([]byte(p))

		case "html":
			p = html.EscapeString(p)

		case "htmle": //html Equivalent decode (Replace " > &quot;)
			p = html.EscapeString(p)
			p = strings.ReplaceAll(p, "&#34;", "&quot;")

		case "hex":
			p = hex.EncodeToString([]byte(p))

		case "json":
			pJson, _ := json.Marshal(p)
			p = string(pJson)

		case "binary":
			var pBinary string

			for _, c := range p {
				pBinary = fmt.Sprintf("%s%.8b", pBinary, c)
			}
			p = pBinary
		}
	}
	return p
}

func PayloadTamper(p, tamper string) string {
	var (
		isLetter = regexp.MustCompile(`^[a-z]$`).MatchString
		and      = regexp.MustCompile(`(AND|and)`)
		or       = regexp.MustCompile(`(OR|or)`)
		s        = " "
		lo_and   = "&&"
		lo_or    = "||"
		lst      []string
	)

	if strings.Contains(tamper, ",") {
		lst = strings.Split(tamper, ",")
	} else {
		lst = append(lst, tamper)
	}

	for _, t := range lst {

		switch t {

		//Space:
		case "s2n":
			p = strings.ReplaceAll(p, s, "")
		case "s2t":
			p = strings.ReplaceAll(p, s, "	")
		case "s2p":
			p = strings.ReplaceAll(p, s, "+")
		case "s2u":
			p = strings.ReplaceAll(p, s, "%20")
		case "s2ut":
			p = strings.ReplaceAll(p, s, "%09")
		case "s2un":
			p = strings.ReplaceAll(p, s, "%00")
		case "s2ul":
			p = strings.ReplaceAll(p, s, "%0a")
		case "s2nl":
			p = strings.ReplaceAll(p, s, "\\n")
		case "s2nrl":
			p = strings.ReplaceAll(p, s, "\\n\\r")
		case "s2c":
			p = strings.ReplaceAll(p, s, "/**/")
		case "s2mc":
			p = strings.ReplaceAll(p, s, "/**0**/")

		//Logical operators:
		case "l2ao":
			p = and.ReplaceAllString(p, "&&")
			p = or.ReplaceAllString(p, "||")

		case "l2bao":
			p = strings.ReplaceAll(p, lo_and, "AND")
			p = strings.ReplaceAll(p, lo_or, "OR")

			//

		//Case modification:
		case "c2r":
			pRandCase := ""
			j := ""
			for _, i := range p {
				j = string(i)

				if isLetter(string(i)) {
					j = RandomCase(string(i))
				}
				pRandCase += j
			}
			p = pRandCase

		//Quote:
		case "q2d":
			p = strings.ReplaceAll(p, "'", "\"")
		case "q2s":
			p = strings.ReplaceAll(p, "\"", "'")

		//Backslash:
		case "b2q":
			p = strings.ReplaceAll(p, "'", "\\'")
			p = strings.ReplaceAll(p, "\"", "\\\"")
		case "b2qd":
			p = strings.ReplaceAll(p, "'", "\\\\'")
			p = strings.ReplaceAll(p, "\"", "\\\\\"")

		case "b2qt":
			p = strings.ReplaceAll(p, "'", "\\\\\\'")
			p = strings.ReplaceAll(p, "\"", "\\\\\\\"")

		case "b2qq":
			p = strings.ReplaceAll(p, "'", "\\\\\\\\'")
			p = strings.ReplaceAll(p, "\"", "\\\\\\\\\"")
		}
	}

	return p
}

/* [TAMPERS]
Space
=====================================
1.      s2n     - <space> : <null>
2.      s2t     - <space> : [tab]
3.      s2p     - <space> : +
4.      s2u     - <space> : %20
5.      s2ut    - <space> : %09
6.      s2un    - <space> : %00
7       s2ul    - <space> : %0a
8.      s2nl    - <space> : \n
9.      s2nrl   - <space> : \n\r
10.     s2c     - <space> : /··/
11.     s2mc    - <space> : /··0··/

Logical operators
=====================================
1.      l2ao    - AND, OR : &&, ||
2.      l2bao   - &&, || : AND, OR

Random
=====================================
1.      l2r     - payload : PaYlOaD

Quote
=====================================
1. q2d  -   ' : "
3. q2s  -   " : '

Backslash
=====================================
1. b2q	-	'," : \',\"
2. b2dq -	'," : \\',\\"
3. b2qq	-	'," : \\\\',\\\\"

*/
