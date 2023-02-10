package parse

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Brum3ns/FireFly/pkg/design"
	fc "github.com/Brum3ns/FireFly/pkg/functions"
	G "github.com/Brum3ns/FireFly/pkg/functions/globalVariables"
	"github.com/projectdiscovery/goflags"
)

//Options variable must be defined here before it's added as a argument inside the function "userArguments()": ("-h/--help")
type Options struct {

	//Standard input
	Url          string
	ReqRaw       string
	Method       string
	Protocol     string
	Headers      string
	PostData     string
	PostDataType string
	UserAgent    string

	//Debug
	Version bool

	//Display
	ShowDiff bool

	//Verify
	Verify         int
	VerifyHeader   bool
	VerifyMethod   string
	VerifyProtocol string
	VerifyPayload  string //<====[TODO] Make it possible to a list of "int", "bool", "string" to test better behaviors. Also fix the "replace function in the CompareDiff when it replace the "payload".
	VerifyChar     string

	//Wordlists
	Wordlist       string
	WordlistFolder string

	//Yaml file
	YmlTransformationFile string

	//Attack techniques
	Attack           string
	AutoDetectParams string
	SplitParam       string

	//Adjustments
	Insert string
	Suffix string
	Prefix string

	//Payload conf
	PayloadSuffix     string
	PayloadPrefix     string
	PayloadPattern    string
	PayloadReplace    string
	PayloadEncodeChar string
	Tamper            string
	Encode            string

	//Priotice
	Priotice string

	//Random
	Random string

	//Match list
	MatchCode   string
	MatchLine   string
	MatchWord   string
	MatchSize   string
	MatchRegex  string
	MatchHeader string

	//Filter list
	FilterCode   string
	FilterLine   string
	FilterWord   string
	FilterSize   string
	FilterRegex  string
	FilterHeader string

	//File input
	FilePayloadYML     string //[TODO]
	FilenameUrls       string
	FilenameRawUrls    string
	FileCompareHeaders string

	//Output data
	/*
		outStatusCode bool
		OutHeaders bool
		OutContentType bool
		OutContentLength bool
		OutRespTime bool
		OutBody bool
		OutBodyByteSize bool
		OutLineCount bool
		OutWordCount bool
		OutRespErr bool
		OutRespDiff bool
		OutFilter bool
		OutHeadersMatch bool
		OutWc   bool
		OutDiff bool */

	//Output file
	OutputJson string
	Output     string

	//Storage of items that came from file user input:
	Lst_urls []string

	Target map[string][]string
	/** Target map will include:
	* methods
	* protocols
	* headersRaw
	* headers
	* compareHeaders
	* postData
	 */

	//Preformance
	Threads int
	Timeout int
	Delay   int
	//Jitter  string

	RandomAgent     bool
	FollowRedirect  bool
	Listtamper      bool
	ShowPayload     bool
	SilenceDarkness bool
	Silence         bool
	Verbose         bool
	ShowConfig      bool
	RegexHeaders    bool
	ShowUrl         bool
	Force           bool
	Adapt           bool
	AdaptPayload    bool
	LiveAnalyze     bool
	Batch           bool
	Overwrite       bool
}

func UserArguments() *Options {
	opt := &Options{}

	design.Banner()

	//Flag menu description
	flagSet := goflags.NewFlagSet()

	desctxt := fmt.Sprintf(`Usage: ./Firefly -u "http://www.example.com/?query=FUZZ" ... [OPTIONS]%v
		%v%vFirely%v is an advance blackbox testing tool to detect respones behavior and adjust it's payloads.
		%vIt use unique wordlists that contains (Fuzz, reflected, Timebased ...) payloads to infect the target
		%vand to minimize false negatives.`, "\r", "\n", "\033[33m", "\033[0m", "\r", "\r")

	flagSet.SetDescription(desctxt)

	CreateGroup(flagSet, "info", "Info",
		flagSet.BoolVarP(&opt.Version, "version", "V", false, "Show current and latest version of FireFly"),
	)

	//User input declaration & description:
	CreateGroup(flagSet, "input", "Input",
		/*OK, pipline=OK*/ flagSet.StringVarP(&opt.Url, "url", "u", "", "URL to attack"),
		/*FIX, pipline=FIX*/ flagSet.StringVarP(&opt.FilenameUrls, "list", "l", "", "File containing URL's to test"),
		/*OK*/ flagSet.StringVarP(&opt.ReqRaw, "raw", "r", "", "HTTP Request raw data to be sent. In quotes separated by new lines. (Addicted of the \"protocol\" option)"),
		//TODO//flagSet.StringVarP(&opt.FilenameRawUrls, "fileraw", "rb", "", "Burp Suite file with HTTP raw requests - (saved item[s])"),
	)

	CreateGroup(flagSet, "request configs", "Request configurations",
		/*OK*/ flagSet.StringVarP(&opt.Protocol, "protocol", "p", "", "Protocol to use. "+design.Example+" [http,https,ftp,gopher ...]"),
		/*OK*/
		flagSet.StringVarP(&opt.Method, "method", "m", "GET", "HTTP method(s) to use separated by comma."),
		/*OK*/ flagSet.StringVarP(&opt.Headers, "headers", "H", "", "Header(s) to include in all requests. Inside quotes separated by new lines"),
		/*OK*/ flagSet.StringVarP(&opt.UserAgent, "user-agent", "ua", "FireFly", "User-Agent to be use"),
		/*OK*/ flagSet.BoolVarP(&opt.RandomAgent, "random-agent", "rua", false, "Random User-Agents"),
		/*OK*/ flagSet.StringVarP(&opt.PostData, "data", "d", "", "Post data to be used in the request. (Best to use with with the \"data-type\" (-dt) option. It will adjust the Content-Type header)"),
		/*FIX*/ //flagSet.StringVarP(&opt.PostDataType, "data-type", "dt", "", "Post data type [post,json,xml] or \"none\" for no Content-Type header to be added."),

		/*FIX, DELETE?*/
		flagSet.StringVarP(&opt.Prefix, "prefix", "px", "", "Add a string to the beginning of the URL(s). "+design.Example+" \"dev.www.target.com/\" (\"dev\")"),
		/*FIX, DELETE?*/ flagSet.StringVarP(&opt.Suffix, "suffix", "sx", "", "Add a string to the end of the URL(s). "+design.Example+" \"/search.html?query=FUZZ;--+-\" (\";--+-\")"),
		//[TODO]//flagSet.BoolVarP(&opt.FollowRedirect, "follow-redirect", "fr", true, "Follow redirects. This makes it possible to check response content in first response and at the redirected destination"),
	)
	CreateGroup(flagSet, "parameters", "Parameters",
		/*OK [TODO] Add cookie*/ flagSet.StringVarP(&opt.AutoDetectParams, "auto-detect", "au", "", "Auto detect GET/POST parameter(s), techniques: [a]ppend / [r]eplace"),
		/*OK (do not detect properly insert points)*/ flagSet.StringVarP(&opt.SplitParam, "chars", "cs", "?&", "Split GET/POST parameters by char"+fmt.Sprintln(`
	---
	1. Default (all)  → ?&
	4. GET params     → get[?&]
	3. POST params    → post[&]
	5. all params     → 'get[?&] post[?&] cookie[?&:]' (Inside quotes separated by [space])
	---`)),
	)

	CreateGroup(flagSet, "verify", "Verify",
		//[NOTE]: Invalid default values is reset in: 'configure.go' -> 'SetupVerify()'
		/*OK*/
		flagSet.IntVarP(&opt.Verify, "verify", "vf", 5, "Verify the original behavior. The amount of verification request to be sent (Recommended amount: 5-9)"),
		/*OK*/ flagSet.StringVarP(&opt.VerifyMethod, "verify-method", "vm", "First original method (-m)", "Verification request method(s) to be used separated by comma"),
		/*FIX - Do test all protocols and not only the given*/ flagSet.StringVarP(&opt.VerifyProtocol, "verify-protocol", "vp", "First original protocol (-p)", "Verification request protocol(s) to be used separated by comma"),
		/***[FIX]***/ //flagSet.BoolVarP(&opt.VerifyHeader, "verify-header", "vh", false, "Use same as the original headers added (-H) and inside raw (-r) if used. If 'false' don't use any user added header(s)"),
		/*OK*/
		flagSet.StringVarP(&opt.VerifyPayload, "verify-payload", "vfp", "13333337", "Verification payload to be used in the process (should be a simple payload of [a-zA-Z0-9])"),
		/*OK*/ flagSet.StringVarP(&opt.VerifyChar, "verify-chars", "vfc", "~!@#$%^&*() -_+={}][|,.\\/?;:`'\"<>", "Verify how special characters are encoded/filtered to detect differences in the reflection of the payload."),
		//flagSet.BoolVarP(&opt.Batch, "batch", "b", false, "Automatically answear questions that pops up during the process ([y/n] inputs)"),
	)

	CreateGroup(flagSet, "payload", "Payload",
		/*OK*/ flagSet.StringVarP(&opt.PayloadReplace, "payload-replace", "pr", "", "Use regex (RE2) to replace parts within the payloads. Use ( => ) as a \"replace to\" indicator. (Spaces are needed) "+design.Example+" \"'\\([0-9]+=[0-9]+\\) => (13=(37-24))'\". Will resul in: From=\"Z'or(1=1)--+-\" To=\"Z'or(13=(37-24))--+-\""),
		/*OK*/ flagSet.StringVarP(&opt.PayloadPattern, "payload-pattern", "pt", "9182", "Payload pattern to use. This makes it possible to spot reflectation/payload changes in the response(s). "+design.Example+" \"9182\" → 9182{PAYLOAD}9182"),
		/*OK*/ flagSet.StringVarP(&opt.PayloadSuffix, "payload-suffix", "ps", "", "Add string to the end of the payload"),
		/*OK*/ flagSet.StringVarP(&opt.PayloadPrefix, "payload-prefix", "pp", "", "Add string to the beginning of the payload"),
		/*OK*/ flagSet.StringVarP(&opt.Tamper, "tamper", "ta", "", "Tamper(s) to use within all the payloads. Multiple tampers can be used separated by a comma. "+design.Example+" \"s2c,q2u\""),
		/*OK*/ flagSet.BoolVarP(&opt.Listtamper, "list-tamper", "lt", false, "List all built in tampers"),
		/*[TODO] FireFly v1.1*/ //flagSet.StringVarP(&opt.FilePayloadYML, "payload-yml", "py", "", "Yaml (.yml) file containing your payload configs for advanced payload fuzzing"),
		//flagSet.StringVarP(&opt.PayloadEncodeChar, "encode-chars", "ec", "", "Encode the specific char[s] within the payloads separated by a comma. "+design.Example+" ',\",\\"),
		/*OK*/
		flagSet.StringVarP(&opt.Encode, "encode", "e", "", "Encode type to be used within the payload (order matter)"+fmt.Sprintln(`
	---
	1. url    [URL encode]
	2. durl   [DoubleURL encode]
	3. b64    [Base64 encode]
	4. b32    [Base64 encode]
	5. html   [HTML encode]
	6. htmle  [HTML Equivalent - HTML encode when " → &quot;]
	7. hex    [Hex encode]
	8. json   [Json encode]
	9. binary [Binary encode]
	---`)),
	)

	CreateGroup(flagSet, "transformation", "Transformation",
		flagSet.StringVar(&opt.YmlTransformationFile, "yml-transformation", "db/yml/transformation.yml", "Yaml file with payload transformation config"),
	)

	CreateGroup(flagSet, "random", "Random",
		flagSet.StringVarP(&opt.Random, "random", "ra", "", "Random [s]tring / [n]umber with length separated by colon (:)"+fmt.Sprint(`
		Adding keywords: "#RANDOM#" and/or "#RANDOMNUM#" in the request will be replaced by a random value using the rules specified.
		(Default: s8,n8)`)),
	)

	CreateGroup(flagSet, "match", "Match",
		flagSet.StringVarP(&opt.MatchCode, "match-code", "mc", "", "Match status code"),
		flagSet.StringVarP(&opt.MatchSize, "match-size", "ms", "", "Match response size"),
		flagSet.StringVarP(&opt.MatchRegex, "match-regex", "mr", "", "Match regex (RE2)"),
		flagSet.StringVarP(&opt.MatchHeader, "match-header", "mh", "", "[In development...] Match header regex and/or string"),
		flagSet.StringVarP(&opt.MatchLine, "match-line", "ml", "", "Match line count"),
		flagSet.StringVarP(&opt.MatchWord, "match-word", "mw", "", "Match word count"),
	)

	CreateGroup(flagSet, "filter", "Filter",
		flagSet.StringVarP(&opt.FilterCode, "filter-code", "fc", "", "Filter status code"),
		flagSet.StringVarP(&opt.FilterSize, "filter-size", "fs", "", "Filter response size"),
		flagSet.StringVarP(&opt.FilterRegex, "filter-regex", "fr", "", "Filter regex (RE2)"),
		flagSet.StringVarP(&opt.FilterHeader, "filter-header", "fh", "", "[In development...] Filter header regex"),
		flagSet.StringVarP(&opt.FilterLine, "filter-line", "fl", "", "Filter line count"),
		flagSet.StringVarP(&opt.FilterWord, "filter-word", "fw", "", "Filter word count"),
	)

	CreateGroup(flagSet, "preformance", "Preformance",
		flagSet.IntVarP(&opt.Timeout, "timeout", "T", 2000, "Timeout in milliseconds (ms) before giving up on the response"),
		flagSet.IntVarP(&opt.Threads, "threads", "t", 15, "Threads/Concurrency to use"),
		flagSet.IntVarP(&opt.Delay, "delay", "dl", 0, "Delay in milliseconds (ms) between each request each thread"),
		//flagSet.IntVarP(&opt.Rate, "rate-limit", "rl", 0, "Delay in milliseconds (ms) between each request each thread"),
		/*[TODO]*/ //flagSet.StringVarP(&opt.Jitter, "d", "", "Random delay range [ms] between requests inside each threads"),
		//flagSet.BoolVarP(&opt.Force, "force", "f", false, "Force the process to continue (This can cause instable results)"),
	)

	/*CreateGroup(flagSet, "verbose", "Verbose") //flagSet.BoolVarP(&opt.ShowUrl, "show-url", "su", false, "Show full URL"),
	//flagSet.BoolVar(&opt.Silence, "s", false, "Do not print output result on screen")
	//flagSet.BoolVar(&opt.SilenceDarkness, "sS", false, "Do not print anything, pure darkness"),


		CreateGroup(flagSet, "Adapt", "Adapt",
			flagSet.BoolVarP(&opt.Adapt, "adapt", "ad", true, "Adapt FireFly technqiues from detected behaviors and or patterns"),
			flagSet.BoolVarP(&opt.AdaptPayload, "adapt-payload", "ap", true, "Adapt payloads from detected behaviors and or patterns"),
		)
	*/

	CreateGroup(flagSet, "attacksTechniques", "Attacks Techniques",
		/*OK*/ flagSet.StringVarP(&opt.Insert, "insert", "i", "FUZZ", "Keyword to be replaced by the payload(s). "+design.Example+" \"http://www.example.com/index.php?query=FUZZ\""),
		/*FIX*/ flagSet.StringVarP(&opt.Attack, "attack", "a", "fuzz", "Only fuzz work for now, [In development...]"), /*  fmt.Sprintln(`Attack mode to use:
		---
		1. fuzz    → Special chars to detect behavior (Errors, Pattern, Crashes, Vuln snippers, Leak etc...)
		2. reflect → Reflected user input (XSS, HTMLI, TI, Entity hijack etc...)
		3. header  → Check/Guess for header that effect the behavior
		4. time    → Time based injection (Code injection, SQLi)
		---`)), */
	)

	CreateGroup(flagSet, "display", "Display",
		/*OK*/ flagSet.BoolVar(&opt.ShowDiff, "show-diff", false, "Display response diff in live view"),
		/*OK*/ flagSet.BoolVar(&opt.ShowConfig, "show-config", false, "Display all configured parses and their values before the process starts"),
		/*OK*/ flagSet.BoolVar(&opt.ShowPayload, "show-payload", false, "Show all payloads and their modifications"),
	)

	CreateGroup(flagSet, "debug", "Debug",
		/*OK*/ flagSet.BoolVarP(&opt.Verbose, "verbose", "v", false, "Verbose output"),
	)

	CreateGroup(flagSet, "wordlist_detection", "Wordlist & Detection",
		/*FIX*/ //flagSet.StringVarP(&opt.FileCompareHeaders, "check-headers", "cH", "db/resources/checkHeaders.txt", "File with headers to check for in the responses"),
		/*OK*/
		flagSet.StringVarP(&opt.WordlistFolder, "wordlistf", "wf", "db/wordlists/", "Folder containing the wordlist(s) to be used. All the wordlist filenames prefix has to contain the attack type (<TYPE>_mywordlist.txt)"),
		/*OK*/ flagSet.StringVarP(&opt.Wordlist, "wordlist", "w", "", "Custom wordlist to be use [fuzz,reflect,time,header] "+design.Example+" \"/path/to/wordlist.txt:fuzz\""),
	)

	CreateGroup(flagSet, "output", "Output",
		flagSet.StringVarP(&opt.OutputJson, "output-json", "oJ", "", "output in JSON format (Recommended, you can use the tool \"jq\" to grep data easier)"),
		flagSet.StringVarP(&opt.Output, "output", "o", "", "Output in plaintext format"),
		/*OK*/ flagSet.BoolVar(&opt.Overwrite, "overwrite", false, "Overwrite existing filename (use carefully)"),
	)

	_ = flagSet.Parse()

	if opt.Version {
		VersionFireFly()
	}

	//Validate user input:
	proceed, msg := ValidateInput(opt)
	if msg != "" {
		fmt.Println(fc.IFFail(msg))
		os.Exit(1)

	} else if !proceed {
		os.Exit(0)
	}

	return opt
}

func ValidateInput(opt *Options) (bool, string) {
	/** Validate user input options
	* Control that the user input is correct set
	* so they can be used in future processes.
	 */
	var err error

	//Declare static values to be global:
	G.RandomOptions = opt.Random
	G.Verify = opt.Verify
	G.Verbose = opt.Verbose
	G.EncodeChar = opt.PayloadEncodeChar
	G.Insert = opt.Insert
	G.PayloadPattern = opt.PayloadPattern

	if len(G.Insert) <= 0 {
		return false, "-i"
	}

	//Create core map to store all request configs:
	opt.Target = Target(opt.Target)
	/*opt.Target["urls"] = []string{}
	opt.Target["protocols"] = []string{}
	opt.Target["methods"] = []string{}
	opt.Target["headers"] = []string{}
	opt.Target["headersRaw"] = []string{}
	opt.Target["compareHeaders"] = []string{}
	opt.Target["postData"] = []string{}*/

	//Declare static header list to a global list:
	if len(opt.FileCompareHeaders) > 0 {
		opt.Target["compareHeaders"] = fc.FileToArray(opt.FileCompareHeaders, opt.Target["compareHeaders"])

		if len(opt.Target["compareHeaders"]) > 0 {
			G.Lst_CheckHeaders = opt.Target["compareHeaders"]
		}
	}

	//Filter & Match regex setup:
	/* G.Lst_mRegexFilter["re"] = opt.FilterRegex //DELETE? - old filter
	G.Lst_mRegexFilter["rh"] = opt.FilterHeader

	G.Lst_mRegexMatch["re"] = opt.MatchRegex
	G.Lst_mRegexMatch["rh"] = opt.MatchHeader */

	//Check *Filter* & *Match* and make a *map[string][]string* of the once being used:
	lst_Fr := map[string]string{
		"sc": opt.FilterCode, "bs": opt.FilterSize,
		"lc": opt.FilterLine, "wc": opt.FilterWord,
		"re": opt.FilterRegex,
	}

	lst_Mh := map[string]string{
		"sc": opt.MatchCode, "bs": opt.MatchSize,
		"lc": opt.MatchLine, "wc": opt.MatchWord,
		"re": opt.MatchRegex,
	}

	//Filter & Match assigned to global values:
	G.Lst_mFilter = MapToMapLst(lst_Fr)
	G.Lst_mMatch = MapToMapLst(lst_Mh)

	//fmt.Println("F:", len(G.Lst_mFilter), "|", "M:", len(G.Lst_mRegexFilter), "| regex: F", len(G.Lst_mMatch), "M:", len(G.Lst_mRegexMatch)) //DEBUG

	//Check Filter/Match [0/1] from response data:
	/* if len(G.Lst_mFilter)+len(G.Lst_mRegexFilter) > 2 { //DELETE? - old filter
		G.MF_Mode = 1
	} else if len(G.Lst_mMatch)+len(G.Lst_mRegexMatch) > 2 {
		G.MF_Mode = 2
	} */

	//Check output format:
	if len(opt.OutputJson) > 0 || len(opt.Output) > 0 {
		G.Output = true
		G.OutputFile, G.OutputType = OutputSetup(opt.OutputJson, opt.Output)

		//Check dosen't exist and if it does ask if the user want to overwrite it.
		if _, err := os.Stat(G.OutputFile); !errors.Is(err, os.ErrNotExist) && !opt.Overwrite {
			return false, "-ov"
		}

		G.OutputFileOS, err = os.Create(G.OutputFile)
		fc.IFError("Could not output data to file.", err)
	}

	//Check Threads:
	if opt.Threads <= 0 {
		return false, "-t"
	}

	//List all tampers:
	if opt.Listtamper {
		fmt.Println(ListTamper())
		return false, ""
	}

	if opt.PayloadReplace != "" && !strings.Contains(opt.PayloadReplace, " => ") {
		return false, "-pr"
	}

	//Read required inputs (Stdin):
	reader, _ := os.Stdin.Stat()
	if (reader.Mode() & os.ModeNamedPipe) > 0 {
		G.Pipe = true
	}

	//Check required inputs:
	if opt.Url == "" && opt.ReqRaw == "" && opt.FilenameRawUrls == "" && opt.FilenameUrls == "" && !G.Pipe {
		return false, "-u"

	}

	//Check verify amount:
	/*if G.Verify <= 1 {
		fmt.Println(design.Fail, "Let verify request atleast make three [3] times to make sure the original response behavior is not dynamic.",
			"\n    Recomended to use atleast five [5+] times to make a stable ground of the target default response[s]")
		return false
	}*/

	return true, ""
}

func CreateGroup(flagSet *goflags.FlagSet, groupName, description string, flags ...*goflags.FlagData) {
	flagSet.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}

func Target(t map[string][]string) map[string][]string {
	/** Target setup
	*  Data that will be used to request target(s)
	 */

	t = make(map[string][]string)

	/*===[Fuzz process]===*/
	t["urls"] = []string{}
	t["protocols"] = []string{}
	t["methods"] = []string{}
	t["headers"] = []string{}
	t["headersRaw"] = []string{}
	t["compareHeaders"] = []string{}
	t["postData"] = []string{}

	/*===[Verify process]===*/
	t["vurls"] = []string{} //Protocol depended
	t["vmethods"] = []string{}
	t["vprotocols"] = []string{}
	t["vheaders"] = []string{}
	//t["headersRaw"] => same
	//t["compareHeaders"] => (not included)
	//t["postData"] => same

	return t
}

func ListTamper() string {

	tamper := fmt.Sprintln(`
Tampers are used to bypass filter and or web application firewalls (WAF's).
This makes it possible to keep a stable process with undetectable payloads.

The following commands change the payloads in different ways to bypass 
filter checks. Before using any tamper[s] a manual test should be performed to
verify the target response and it's filter checks.

- More than one tamper can be used and are separated with "," without space.

Space
=====================================
1.	s2n` + "\t" + `- <space> : <null>
2.	s2t` + "\t" + `- <space> : [tab]
3.	s2p` + "\t" + `- <space> : +
4.	s2u` + "\t" + `- <space> : %20
5.	s2ut` + "\t" + `- <space> : %09
6.	s2un` + "\t" + `- <space> : %00
7	s2ul` + "\t" + `- <space> : %0a
8.	s2nl` + "\t" + `- <space> : \n
9.	s2nrl` + "\t" + `- <space> : \n
10.	s2c` + "\t" + `- <space> : /**/ )
11.	s2mc` + "\t" + `- <space> : /**0**/

Logical operators
=====================================
1.	l2ao` + "\t" + `- AND, OR : &&, ||
2.	l2bao` + "\t" + `- &&, || : AND, OR

Case modification
=====================================
1.	c2r` + "\t" + `- payload : PaYlOaD

Quote
=====================================
1. q2d` + "\t" + `-   ' : "
3. q2s` + "\t" + `-   " : '

Backslash
=====================================
1. b2q	-	'," : \',\"
2. b2qd -	'," : \\',\\"
3. b2qt	-	'," : \\\',\\\"
4. b2qq	-	'," : \\\\',\\\\"`)

	return tamper
}
