package parse

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Brum3ns/FireFly/pkg/design"
	"github.com/Brum3ns/FireFly/pkg/firefly/technique"
	fc "github.com/Brum3ns/FireFly/pkg/functions"
	G "github.com/Brum3ns/FireFly/pkg/functions/globalVariables"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

//[TODO]
/*type Config struct {
 opt  *Options
	wl   *storage.Wordlists
	conf storage.Configure
}*/

//Configure user options before taking action in the process:
func Configure(opt *Options, wl *storage.Wordlists) (bool, string) {
	/** Configure all user parse
	+-------------------------------+
	| - Valid/Invalid input conf	|
	| - Wordlist conf				|
	| - Payload conf				|
	| - Attack conf					|
	| - Detection conf				|
	| - Filter conf					|
	| - Request conf				|
	+-------------------------------+
	*/

	//[Verbose] Setup all detection wordlists & User-Agents:
	fmt.Println(design.Info, "Grep Wordlist[s] setup in process")
	fc.LstGSetup(wl)

	G.Lst_RandomAgent = fc.FileToLst("db/resources/randomUserAgent.txt")
	print("\r") //<= Delete the loading junk left over.

	//Setup 'Chars' (Option: -verify-char) to be tested:
	if msg, err := Conf_CharPayload(opt.VerifyChar, wl); err != nil {
		return false, msg
	}

	//Verify '#RANDOM#' insert keyword(s):
	if msg := Conf_RandomKeyword(opt.Random); msg != "" {
		return false, msg
	}

	// If custom attack techniques was given add them: [TODO] - (Fix more attack methods)
	if msg := Conf_AttackToUse(opt.Attack, wl); msg != "" {
		return false, msg
	}

	//Setup up request(s):
	if msg := Conf_RequestSetup(opt); msg != "" {
		return false, msg
	}

	//If a post body is added with (-d) then use priotice it:
	if len(opt.PostData) > 0 {
		opt.Target["postData"] = []string{opt.PostData}
	}

	//Auto detect parameters and add insert (FUZZ):
	if au := strings.ToLower(opt.AutoDetectParams); len(au) > 0 {
		if (au == "replace" || au == "r") || (au == "append" || au == "a") {
			//Detect params to fuzz method given:
			var err error
			if len(opt.Target["urls"]) > 0 {
				for id, url := range opt.Target["urls"] {
					url, err = fc.FuzzParam(url, "get", opt.SplitParam, G.Insert, au)
					fc.IFError("p", err)

					if len(url) > 0 {
						opt.Target["urls"][id] = url
					}
				}
			}
			if len(opt.PostData) > 0 {
				l_AutoParam, err := fc.FuzzParam(opt.PostData, "post", opt.SplitParam, G.Insert, au)
				fc.IFError("p", err)
				opt.Target["postData"] = []string{l_AutoParam}
			}
		} else {
			return false, "-au"
		}
	}

	//Read & Setup yml file containing paylods and transformations:
	if msg := Conf_Wordlist(opt.YmlTransformationFile, wl); msg != "" {
		return false, msg
	}

	//Add custom wordlist[s] to use: //[TODO] - Clear and code better
	if opt.Wordlist != "" {
		if !strings.Contains(opt.Wordlist, ":") {
			return false, "-w"
		}
		fc.StringToArray_Attack(opt.PayloadReplace, opt.Encode, opt.Tamper, opt.Wordlist, wl, false)

	} else {
		fc.WordlistFolder_Declare(opt.PayloadReplace, opt.Encode, opt.Tamper, opt.WordlistFolder, wl)
	}

	//[TODO] Check file input and add it to the memory storage:
	/*if opt.FilenameUrls != "" {
		opt.Lst_urls = ToArray(opt, opt.FilenameUrls, opt.Lst_urls, "url")
	}

	//Burp HTTP raw file extract url & raw data and store it into lists:
	if opt.FilenameRawUrls != "" {
		opt.Lst_urls = ToArray(opt, opt.FilenameRawUrls, opt.Lst_urls, "raw")

	}*/

	if len(G.PayloadPattern) < 4 {
		return false, "Payload pattern (-ptn) must atleast be four chars long. This is needed to detect reflectations, payload changes and differences in the responses"
	}

	return true, "Configuration completed"
}

func Conf_Wordlist(ymlFile string, wl *storage.Wordlists) string {
	//Read yml file containing paylods and transformations:
	if ymlFile != "" {
		wl.TransformationPayloads, wl.TransformationCompare = technique.Transformation(ymlFile)
		return ""
	}
	return "tfmt"
}

func Conf_CharPayload(s string, wl *storage.Wordlists) (string, error) {
	m := make(map[rune]string)
	for _, r := range s {
		m[r] = fc.PayloadChar(r)
	}

	wl.VerifyPayloadChars = m
	return "-vfc", nil
}

func Conf_RandomKeyword(r string) string {
	/**Validate and setup the random keyword(s) - (%RANDOM%, %RANDOMSTR%, %RANDOMINT%)
	Return if the process failed/invalid input
	*/

	//If no user preference are set:
	if r == "" {
		return ""
	}

	//Adjust to user preference (Only 3 loops max so 'regexp.MatchString()' is fine to use here.):
	for _, i := range strings.Split(r, ",") {
		if ok, _ := regexp.MatchString(`^[sn]:[0-9]+$`, strings.ToLower(i)); !ok {
			return "-ra"
		}
		v := strings.Split(i, ":")
		ln, _ := strconv.Atoi(string(v[1]))
		G.M_Random[string(v[0])] = ln
	}

	return ""
}

func Conf_AttackToUse(a string, wl *storage.Wordlists) string {
	wl.UseTechniques = []string{"fuzz", "transformation"}

	//DELETE - Junk code...
	/* if a != "all" {
		attack := strings.ToLower(a)

		technique.Attack_techq = strings.Split(attack, ",")

		//Check if it's a valid attack technique. If true. Add it:
		for _, techq := range technique.Attack_techq {
			if fc.InLst(wl.Lst_Valid_Techniques, techq) {
				wl.Lst_UseTechniques = append(wl.Lst_UseTechniques, techq)

			} else {
				return "-a"
			}
		}
		//Use *all* -> default techniques:
	} else {
		wl.Lst_UseTechniques = wl.Lst_Valid_Techniques
	}*/
	return ""
}

func Conf_RequestSetup(opt *Options) string {

	var (
		lst_url []string
		msg     string
	)

	//[RAW] : Single URL/HTTPRawData given by "opt.url" or "opt.RawData" and add it to the list:
	if opt.ReqRaw != "" {
		//Only for user awareness
		if opt.Url != "" || G.Pipe {
			return "stdin"
		}

		//Extract the URL from the input HTTP Raw data: (Unless an other URL was speficied. Ex: for host header manupulation)
		G.Lst_rawData = make(map[string][]string)
		G.Lst_rawData, msg = fc.RawData(opt.ReqRaw)

		if len(msg) > 0 {
			return msg
		}

		//Define request HTTP raw data:
		opt.Target["headersRaw"] = G.Lst_rawData["headers"]
		opt.Target["postData"] = G.Lst_rawData["postData"]

		lst_url = append(lst_url, G.Lst_rawData["url"][0])

		fmt.Println("\n------\n\r" + strings.Join(G.Lst_rawData["headers"], "\n") + "\n------\n") //DEBUG

		//[URL] : Set single URL(s)
	} else if opt.Url != "" {
		if G.Pipe {
			return "stdin"
		}
		lst_url = append(lst_url, opt.Url)

		//[PIPELINE] : Read stdin data - URL(s):
	} else if G.Pipe {
		lst_url, msg = Pipeline()
		if len(msg) > 0 {
			return msg
		}

		//[FILE] : Take URL(s) from a file:
	} else if opt.FilenameUrls != "" {
		lst_url = fc.FileToLst(opt.FilenameUrls)

		//Invalid input or none:
	} else {
		return "-u"
	}

	//Check HTTP method to use:
	if opt.Method != "" {
		opt.Target["methods"] = append(opt.Target["methods"], fc.ToLstSplit(opt.Method, ",")...)
	}

	//Check if a custom protocol has been added (-p). If so, add it/them:
	if opt.Protocol != "" {
		opt.Target["protocols"] = append(opt.Target["protocols"], fc.ToLstSplit(opt.Protocol, ",")...)

	} else if lst_protocols := EAProto(lst_url); len(lst_protocols) > 0 {
		opt.Target["protocols"] = append(opt.Target["protocols"], lst_protocols...)
	} else {
		opt.Target["protocols"] = append(opt.Target["protocols"], G.DefaultProto)
	}

	//Add static header provided by the user option:
	if opt.Headers != "" {
		opt.Target["headers"] = append(opt.Target["headers"], fc.ToLstSplit(opt.Headers, "\n")...)
	}

	//Add postdata priotice the user input (If raw postData is in place)
	if opt.PostData != "" {
		opt.Target["postData"] = []string{opt.PostData} //Do not use append.
	}

	//Setup & configure the verify requests that is used to check default response(s) from target(s):
	if ok, msg := SetupVerify(opt); !ok {
		fc.IFFail(msg)
	}

	// =====[Setup Verify & Fuzz URL(s)]=====
	opt.Target["vurls"] = AddUrls(opt.Target["vprotocols"], lst_url)
	opt.Target["urls"] = AddUrls(opt.Target["protocols"], lst_url)

	return ""
}

func CheckProto(u string) (string, bool) {
	/**Check for valid protocols within a url
	* return the protocol if any
	 */

	if re, _ := regexp.MatchString("[a-zA-Z]+://*(.*?)*($|/)", u); re {
		p := fc.RegexBetween(u, "^*(.*?)*://")
		if p != "" {
			return p, true
		}
	}
	return "", false
}

func EAProto(lstu []string) []string {
	/** Extract & Add (EA) protocols from url(s) if none was added.
	* Return list of protocols gathered from the url(s)
	 */
	var lp []string
	for _, u := range lstu {
		if p, ok := CheckProto(u); ok {
			lp = append(lp, p)
		}
	}
	return lp
}

func AddUrls(lstProto, lstUrl []string) []string {
	/** Collect all protocol(s) and URL(s)
	* Return the new list with all unique URL(s)
	 */

	var l, lp, lu []string

	//Clear url & protocol and add to the seprated lists:
	lp = append(lp, lstProto...)

	for _, u := range lstUrl {
		if p, ok := CheckProto(u); ok {
			lp = append(lp, p)
			u = strings.Replace(u, (p + "://"), "", 1)
		}
		lu = append(lu, u)
	}

	//Remove duplicates
	lp = fc.LstRmDups(lp)
	lu = fc.LstRmDups(lu)

	for _, p := range lp {
		for _, u := range lu {
			l = fc.ToLst_s(l, ((p + "://") + u))
		}
	}

	return l
}

func Pipeline() ([]string, string) {
	/** Append pipeline input add it to a lst and return it
	* The input reads stdin add it to a list that will alter be used as the target
	 */

	var (
		lst []string
		scn *bufio.Scanner
	)

	//Read stdin data:
	reader, _ := os.Stdin.Stat()
	if reader.Mode()&os.ModeNamedPipe > 0 {
		scn = bufio.NewScanner(os.Stdin)

	} else {
		return nil, "stdin"
	}

	//Extract and verify data
	for scn.Scan() {

		i := scn.Text()
		if !fc.InLst(lst, i) || !strings.ContainsAny(i, "|| |\t|\n") {
			lst = append(lst, i)
		}
	}
	return lst, ""
}

//[TODO] I don't know what the FUCK I wrote here but change it...
func ToArray(opt *Options, file string, lst []string, typ string) []string {

	var (
		item    string
		scanner *bufio.Scanner
	)

	//Read URL[s] from stdin pipeline input:
	if typ == "stdin" {
		reader, _ := os.Stdin.Stat()

		if reader.Mode()&os.ModeNamedPipe > 0 {
			scanner = bufio.NewScanner(os.Stdin)

		} else {
			fmt.Println("[\033[31mx\033[0m] Could not read stdin (pipeline) input.")
			os.Exit(0)
		}

		//The stdin is used to only import URL[s] set "typ" to url for the next process:
		typ = "url"

		//Standard URL list
	} else if len(file) > 0 {
		reader, _ := os.Open(file)
		scanner = bufio.NewScanner(reader)

	}

	for scanner.Scan() {
		if scanner.Text() == "" || scanner.Text() == " " || scanner.Text() == "\t" || scanner.Text() == "\n" {
			continue
		}

		if typ == "url" {
			item = opt.Prefix + scanner.Text() + opt.Suffix
		} else {
			item = scanner.Text()
		}

		if !fc.InLst(lst, item) {

			//Extract the original URL protocol and add it to the "Target["protocols"]" unless "-fp":
			protocol := fc.RegexBetween(item, `^(.*)://`)
			if len(protocol) > 0 && !fc.InLst(opt.Target["protocols"], protocol) {
				opt.Target["protocols"] = append(opt.Target["protocols"], []string{protocol}...)
				//opt.Target["protocols"] = append(opt.Target["protocols"], protocol)
			}

			//Clear the protocol from the URL
			item = strings.Replace(item, (protocol + "://"), "", 1)

			//Add all protocol from the "Target["protocols"]":
			for _, proto := range opt.Target["protocols"] {
				lst = append(lst, (proto + "://" + item))
			}
		}
	}
	return lst
}

func OutputSetup(s1, s2 string) (string, string) {
	/** Output options
	*	s1	:	opt.OutputJson			:	json
	*	s2	: 	opt.Output				:	output
	 */

	var (
		lst = []string{s1, s2}
		f   string
	)

	for i := range lst {
		if len(lst[i]) > 0 {
			switch i {
			case 0:
				f = s1
				return f, "json"

			case 1:
				f = s2
				return f, "output"

			default:
				break
			}
		}
	}
	return f, ""
}

func MapToMapLst(lst map[string]string) map[string][]string {
	var l = map[string][]string{}

	for x, y := range lst {
		if len(y) > 0 {
			l[x] = strings.Split(y, ",")
		}
	}
	return l
}

func SetupVerify(opt *Options) (bool, string) {
	/**Setup the verify request(s) variables. The once that aren't included use same variables as the original request variables.
	* Update the options[.]Target map.
	 */
	var (
		vm []string
		vp []string
		vh []string
	)

	//[vMethod] - (Remove default junk text)
	if opt.VerifyMethod == "First original method (-m)" {
		vm = append(vm, opt.Target["methods"][0])
	} else {
		vm = fc.ToLstSplit(opt.VerifyMethod, ",")
	}
	//[vProtocol] - (Remove default junk text)
	if opt.VerifyProtocol == "First original protocol (-p)" {
		vp = append(vp, opt.Target["protocols"]...)
	} else {
		vp = fc.ToLstSplit(opt.VerifyProtocol, ",")
	}

	//[vHeaders]
	if opt.VerifyHeader {
		vh = opt.Target["headers"]
	}

	//Check so atleast one method & protocol is within the original request variables and add the matched once:
	vp = fc.LstAppendSame(opt.Target["protocols"], vp)
	vm = fc.LstAppendSame(opt.Target["methods"], vm)

	if len(vp) <= 0 || len(vm) <= 0 {
		return false, "v:mp"
	}

	//Set all configured value(s) to the opt[.]Target map:
	opt.Target["vmethods"] = vm
	opt.Target["vprotocols"] = vp
	opt.Target["vheaders"] = vh

	return true, ""
}
