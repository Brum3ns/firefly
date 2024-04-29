package option

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Brum3ns/firefly/internal/fail"
	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/internal/info"
	"github.com/Brum3ns/firefly/internal/update"
	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/Brum3ns/firefly/pkg/files"
	"github.com/Brum3ns/firefly/pkg/parameter"
	"github.com/Brum3ns/firefly/pkg/request"
	"golang.org/x/exp/slices"
)

// [README]
// The Option struct holds all options (flags) provided by the user.
// - The order of this struct will effect the helpmenu displayed to the CLI.
// - Local variables are *only* set to be included in the helpmenu.
type Options struct {
	Input
	Request
	Verify
	Wordlist
	Payload
	Filter
	Display
	Preformance
	General
	File
}

// ////////////// Input //////////////// //
type Input struct {
	url        string          `flag:"u" errorcode:"0"`  //<-local
	technique  string          `flag:"tq" errorcode:"0"` //<-local
	Techniques map[string]bool `flag:"" errorcode:"1003"`
	ReqRaw     string          `flag:"r" errorcode:"1001"`
}

// ////////////// Diff //////////////// //
/* type diff struct {
	LevelDiff       int  `flag:"lvl" errorcode:"2001"`
	DiffBody_skip   bool `flag:"skip-body" errorcode:"2002"`
	DiffHeader_skip bool `flag:"skip-header" errorcode:"2003"`
} */

// ////////////// Filter //////////////// //
type Filter struct {
	MatchMode/*(OR|AND)*/ string         `flag:"mmode" errorcode:"3001"`
	MatchCode                     string `flag:"mc" errorcode:"3002"`
	MatchLine                     string `flag:"ms" errorcode:"3003"`
	MatchWord                     string `flag:"ml" errorcode:"3004"`
	MatchSize                     string `flag:"mw" errorcode:"3005"`
	MatchTime                     string `flag:"mt" errorcode:"3006"`
	MatchBodyRegex                string `flag:"mr" errorcode:"3007"`
	MatchHeaderRegex              string `flag:"mh" errorcode:"3008"`
	MatchHeaderAmount             string `flag:"mH" errorcode:"3009"`
	FilterMode/*(OR|AND)*/ string        `flag:"fmode" errorcode:"3010"`
	FilterCode                    string `flag:"fc" errorcode:"3011"`
	FilterLine                    string `flag:"fs" errorcode:"3012"`
	FilterWord                    string `flag:"fl" errorcode:"3013"`
	FilterSize                    string `flag:"fw" errorcode:"3014"`
	FilterTime                    string `flag:"ft" errorcode:"3015"`
	FilterBodyRegex               string `flag:"fr" errorcode:"3016"`
	FilterHeaderRegex             string `flag:"fh" errorcode:"3017"`
	FilterHeaderAmount            string `flag:"fH" errorcode:"3018"`
}

// ////////////// Output //////////////// //
type File struct {
	Output string `flag:"o" errorcode:"4001"`
}

// ////////////// Display //////////////// //
type Display struct {
	Version     bool `flag:"version" errorcode:"5001"`
	NoDisplay   bool `flag:"no-display" errorcode:"5002"`
	Listtampers bool `flag:"list-tampers" errorcode:"5003"`
	ShowConfig  bool `flag:"show-config" errorcode:"5004"`
	TerminalUI  bool `flag:"tui" errorcode:"5005"`
	Verbose     bool `flag:"v" errorcode:"5006"`
	Detail      bool `flag:"detail" errorcode:"5007"`
}

// ////////////// Payload //////////////// //
type Payload struct {
	PayloadReplace string   `flag:"pr" errorcode:"8001"`
	PayloadPattern string   `flag:"pt" errorcode:"8002"`
	PayloadSuffix  string   `flag:"ps" errorcode:"8003"`
	PayloadPrefix  string   `flag:"pp" errorcode:"8004"`
	Tamper         string   `flag:"tamper" errorcode:"8005"`
	Encode         []string `flag:"e" errorcode:"8006"`
	encode         string   `flag:"e" errorcode:"0"` //<-local
	InsertKeyword  string   `flag:"insert" errorcode:"8007"`
}

// ////////////// Wordlist //////////////// //
type Wordlist struct {
	wordlistPath           string   `flag:"w" errorcode:"9001"`
	WordlistPaths          []string `flag:"w" errorcode:"9001"`
	TransformationYAMLFile string   `flag:"yml-tfmt" errorcode:"9003"`
	wordlistValid          bool
}

// ////////////// Request //////////////// //
type Request struct {
	header       string                          `flag:"H" errorcode:"0"`      // <-Local
	method       string                          `flag:"X" errorcode:"0"`      // <-Local
	scheme       string                          `flag:"scheme" errorcode:"0"` // <-Local
	randomInsert string                          `flag:"random" errorcode:"0"` // <-Local
	paramRules   map[string]parameter.QueryRules `flag:"au" errorcode:"10009"` // <-Local
	URLs         []string                        `flag:"" errorcode:"10001"`
	Methods      []string                        `flag:"" errorcode:"10002"`
	Scheme       []string                        `flag:"" errorcode:"10003"`
	Proxy        string                          `flag:"proxy" errorcode:"10004"`
	PostData     string                          `flag:"d" errorcode:"10006"`
	UserAgent    string                          `flag:"ua" errorcode:"10007"`
	SkipHeaders  string                          `flag:"sH" errorcode:"10008"`
	Delay        int                             `flag:"delay" errorcode:"10011"`
	Timeout      int                             `flag:"timeout" errorcode:"10012"`
	HTTP2        bool                            `flag:"timeout" errorcode:"10010"`
	RandomAgent  bool                            `flag:"rua" errorcode:"10013"`
	Headers      [][2]string                     `flag:"" errorcode:"100015"`
	Random       map[string]int                  `flag:"" errorcode:"100014"`
	Hosts        map[string]request.Host         `flag:"" errorcode:"10000"`
	Params       map[string]parameter.Parameter  `flag:"" errorcode:"100016"`
}

// ////////////// Preformance //////////////// //
type Preformance struct {
	Threads             int `flag:"t" errorcode:"11001"`
	ThreadsScanner      int `flag:"tS" errorcode:"11003"`
	ThreadsExtract      int `flag:"tE" errorcode:"11002"`
	MaxIdleConns        int `flag:"idle" errorcode:"11004"`
	MaxIdleConnsPerHost int `flag:"idle-host" errorcode:"11005"`
	MaxConnsPerHost     int `flag:"conn-host" errorcode:"11006"`
	//RateLimit           int `flag:"rate" errorcode:"11007"`
}

// ////////////// General //////////////// //
type General struct {
	//Color          bool `flag:"c" errorcode:"12001"`
	Overwrite      bool `flag:"overwrite" errorcode:"12002"`
	UpdateResource bool `flag:"uR" errorcode:"12003"`
}

// ////////////// Verify //////////////// //
type Verify struct {
	VerifyAmount  int    `flag:"vf" errorcode:"13001"`
	VerifyPayload string `flag:"vP" errorcode:"13002"`
	//VerifyChar    string `flag:"vC" errorcode:"13003"`
}

func NewOptions() *Options {
	opt := &Options{}
	opt.paramRules = make(map[string]parameter.QueryRules)
	opt.Params = make(map[string]parameter.Parameter)
	opt.Hosts = make(map[string]request.Host)
	opt.Methods = []string{"GET"}
	opt.Scheme = []string{"http"}
	opt.URLs = []string{}
	opt.Techniques = map[string]bool{
		"D": true,
		"E": true,
		"T": true,
		"X": false,
	}
	opt.Random = map[string]int{
		"s": 8,
		"n": 8,
	}

	//TODO
	//flag.BoolVar(&opt.Color, "c", false, "Add colors to the screen output")
	//flag.StringVar(&opt.SkipHeaders, "sH", global.FILE_SKIP_HEADERS, "Header(s) to threat as uninteresting in the response when doing difference checks")
	//flag.StringVar(&opt.Tamper, "tamper", "", "Tamper(s) to use within all the payloads. Multiple tampers can be used *separated by a comma*. "+exampleValues(" \"s2c,q2u\""))
	//flag.IntVar(&opt.RateLimit, "rate", 0, "Request rate limit (0 = unlimited)")
	//flag.StringVar(&opt.VerifyChar, "vC", "~!@#$%^&*() -_+={}][|,.\\/?;:`'\"<>", "Verify how special characters are encoded/filtered to detect differences in the reflection of the payload.")
	/* flag.StringVar(&opt.SplitParam, "pS", "?&", "Split GET/POST parameters by char"+fmt.Sprintln(`
	---
	1. Default (all)  → ?&
	4. GET params     → get[?&]
	3. POST params    → post[&]
	5. all params     → 'get[?&] post[?&] cookie[?&:]' (Inside quotes (') *separated by [SPACE]*)
	---`)) */

	flag.BoolVar(&opt.Version, "version", false, "Show current and latest version of FireFly, then exit.")
	flag.Func("u", "The URL(s) to preform black-test on *separated by comma*, if a comma is used wihtin the header value simply escape it with a backslash (\\,)", opt.setURLs_NoScheme)
	flag.Func("scheme", "HTTP scheme to use http[s]", opt.setScheme)
	flag.Func("H", "Header(s) to include in all requests *separated by comma*, if a comma is used wihtin the header value simply escape it with a backslash (\\,)", opt.setHeaders)
	flag.Func("X", "HTTP method(s) to use *separated by comma* (all = all methods except \"DELETE\". To add method \"DELETE\", do \"all,delete\")", opt.setMethods)
	flag.Func("r", "HTTP Request raw data to be sent. In quotes *separated by new lines*. (Addicted of the \"scheme\" option)", opt.setRaw)
	flag.Func("random", `Random [s]tring / [n]umber with a digit at the end to set the length. Both can be set *separeted by a comma*. The keyword(s): "#RANDOM#" / "#RANDOMNUM#" will be replaced with a random value`, opt.setRandomInsert)
	flag.Func("e", "Encode type to be used within the payload (order matter) *separated by a comma*. "+support_encodes(), opt.setEncode)
	flag.Func("au", "Auto detect parameters. More than one can be added *separated by comma*. "+support_autoParameters()+". "+support_format("{param_postion}:{[r]eplace|[a]ppend}:{separators}")+"\n\t\tThe last option (separators) is optional. Note that in \"url\" the \"?\" is added by default. In case you must use \":\" as a separator escape it as \"\\:\".\n\t\t"+exampleValues("url:replace:& | body:a | body:append,url:replace:&;,cookie:replace")+"\n", opt.setAutoParamRules)

	flag.StringVar(&opt.technique, "tq", "ETD", "Technique(s) to be used within the process ([D]iff, [E]xtract, [T]ransformation or [X] to disable all techniques) by letter")

	//- [ Request ] -
	flag.BoolVar(&opt.HTTP2, "http2", false, "Use HTTP/2 otherwise use HTTP/1.1")
	flag.StringVar(&opt.UserAgent, "ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36 Firefly", "User-Agent to be use")
	flag.BoolVar(&opt.RandomAgent, "rua", false, "Random User-Agents")
	flag.StringVar(&opt.PostData, "d", "", "Post data to be used in the request")
	flag.StringVar(&opt.Proxy, "proxy", "", "Proxy to use, "+exampleValues("http://127.0.0.1:8080"))

	//- [ Verify ] -
	flag.IntVar(&opt.VerifyAmount, "vf", 10, "Verify the original behavior. The amount of verification request to be sent (Recommended amount: 5-9)")
	flag.StringVar(&opt.VerifyPayload, "vP", "13333337", "Verification payload to be used in the process (should be a simple payload of [a-zA-Z0-9])")

	flag.StringVar(&opt.InsertKeyword, "insert", "FUZZ", "Payload insert point to be replaced with the payload")
	flag.StringVar(&opt.PayloadReplace, "pr", "", "Use regex (RE2) to replace parts within the payloads. Use ( => ) as a \"replace to\" indicator. (Spaces are needed) "+exampleValues(" \"'\\([0-9]+=[0-9]+\\) => (13=(37-24))'\". Will resul in: From=\"Z'or(1=1)--+-\" To=\"Z'or(13=(37-24))--+-\""))
	flag.StringVar(&opt.PayloadPattern, "pt", "9182", `Pattern of payload to be used. If this is set to none, it will be harder to detect payload reflected payload changes in the response(s). `+exampleValues("\"9182\" → 9182{PAYLOAD}9182"))
	flag.StringVar(&opt.PayloadSuffix, "ps", "", "Add string to the end of the payload")
	flag.StringVar(&opt.PayloadPrefix, "pp", "", "Add string to the beginning of the payload")
	flag.BoolVar(&opt.Listtampers, "list-tampers", false, "List all built in tampers")

	//- [ Transformation ] -
	flag.StringVar(&opt.TransformationYAMLFile, "yml-tfmt", global.FILE_TRANSFORMATION, "Yaml file with payload transformation config")

	//- [ Match ] -
	flag.StringVar(&opt.MatchMode, "mmode", "or", "Match mode (AND|OR)")
	flag.StringVar(&opt.MatchCode, "mc", "", "Match status code")
	flag.StringVar(&opt.MatchSize, "ms", "", "Match response size")
	flag.StringVar(&opt.MatchLine, "ml", "", "Match line count")
	flag.StringVar(&opt.MatchWord, "mw", "", "Match word count")
	flag.StringVar(&opt.MatchTime, "mt", "", "Match response time")
	flag.StringVar(&opt.MatchBodyRegex, "mr", "", "Match body regex (RE2)")
	flag.StringVar(&opt.MatchHeaderRegex, "mh", "", "Match header regex (RE2)")
	flag.StringVar(&opt.MatchHeaderAmount, "mH", "", "Match header amount")

	//- [ Filter ] -
	flag.StringVar(&opt.FilterMode, "fmode", "or", "Filter mode (AND|OR)")
	flag.StringVar(&opt.FilterCode, "fc", "", "Filter status code")
	flag.StringVar(&opt.FilterSize, "fs", "", "Filter response size")
	flag.StringVar(&opt.FilterLine, "fl", "", "Filter line count")
	flag.StringVar(&opt.FilterWord, "fw", "", "Filter word count")
	flag.StringVar(&opt.FilterTime, "ft", "", "Filter response time")
	flag.StringVar(&opt.FilterBodyRegex, "fr", "", "Filter body regex (RE2)")
	flag.StringVar(&opt.FilterHeaderRegex, "fh", "", "Filter header regex (RE2)")
	flag.StringVar(&opt.FilterHeaderAmount, "fH", "", "Filter header amount")

	//- [ Preformance ] -
	flag.IntVar(&opt.Timeout, "timeout", 11, "Timeout in secounds before giving up on the response")
	flag.IntVar(&opt.Threads, "t", 50, "Threads (requests)")
	flag.IntVar(&opt.ThreadsScanner, "tS", 3, "Number of processes to be run in the scanner (this can take up a lot of CPU usage if the value is too high)")
	flag.IntVar(&opt.ThreadsExtract, "tE", 2, "Threads to be used to extract patterns from target response data (hardware)")
	flag.IntVar(&opt.Delay, "delay", 0, "Delay in milliseconds (ms) between each request each thread")
	flag.IntVar(&opt.MaxIdleConns, "idle", 1000, "Controls the maximum number of idle (keep-alive) connections across all hosts")
	flag.IntVar(&opt.MaxIdleConnsPerHost, "idle-host", 500, "Controls the maximum idle (keep-alive) connections to keep per-host")
	flag.IntVar(&opt.MaxConnsPerHost, "conn-host", 500, "Limits the total number of connections per host")

	//- [ Display ] -
	/*In development*/ //flag.BoolVar(&opt.TerminalUI, "tui", false, "Use advanced terminal user interface (UI)")
	flag.BoolVar(&opt.Detail, "detail", false, "Show the difference discovered in an unexpected behavior")
	flag.BoolVar(&opt.NoDisplay, "no-display", false, "Do not display result to screen")
	flag.BoolVar(&opt.ShowConfig, "show-config", false, "Display all configured parses and their values before the process starts")

	//- [ Debug ] -
	flag.BoolVar(&opt.Verbose, "v", false, "Display Verbose")

	//- [ Wordlist ] -
	flag.StringVar(&opt.wordlistPath, "w", global.DIR_WORDLIST, "Wordlist to be used. A single wordlist can be selected or a folder containing wordlists (files must have an \"txt\" extension if a folder is used) "+exampleValues("\"/path/to/wordlist.txt\""))

	//- [ Output ] -
	flag.StringVar(&opt.Output, "o", "", "Output result to given file (JSON format)")
	flag.BoolVar(&opt.Overwrite, "overwrite", false, "Overwrite the existing file name to be used as the output file (use carefully)")

	//- [ Update ] -
	flag.BoolVar(&opt.UpdateResource, "uR", false, "Update the resources that Firefly uses for discovery")

	//Set custom usage menu output:
	flag.Usage = opt.customUsage
	flag.Parse()

	//Show Firefly version OR update the resources:
	switch {
	case opt.Version:
		fmt.Println(info.VERSION)
		os.Exit(0)
		//VersionFirefly()
	case opt.UpdateResource:
		update.Resources()
	}

	//Set Option values that wasen't possible in the flag process direcly:
	if err := opt.setTechniques(); err != nil {
		log.Fatal(err)
	}

	//Read input from pipeline STDIN:
	if err := opt.readStdin(); err != nil {
		log.Fatalf("Invalid stdin was given: %s", err)
	}

	//Setup the parameter object for each supported position within the HTTP request
	if err := opt.makeURLs(); err != nil {
		log.Fatal(err)
	}
	//Setup all params in related to their host and supported position within the HTTP request
	if err := opt.makeParams(); err != nil {
		log.Fatal(err)
	}

	opt.Headers = append(opt.Headers, [2]string{"user-agent", opt.UserAgent})

	if err := opt.makeWordlist(); err != nil {
		log.Fatal(design.STATUS.ERROR, err)
	}

	//Configure all the options (user input):
	configuredOptions, errcode := Configure(opt)
	if errcode > 0 {
		fail.IFFail(errcode)
	}

	//Preview options to be shown on the screen:
	if opt.ShowConfig {
		opt.showConfigOnScreen()
	}

	return configuredOptions
}

// Read stdin and add the given input to the 'Options struct' (if any)
func (opt *Options) readStdin() error {
	if data, err := os.Stdin.Stat(); err != nil {
		return err

	} else if data.Mode()&os.ModeNamedPipe > 0 {
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			//Add collected URL to the 'Options struct' variable:
			if u, ok := opt.validateURL(scanner.Text()); ok && !slices.Contains(opt.URLs, u) {
				opt.URLs = append(opt.URLs, u)
			}
		}
	}
	return nil
}

// Validate the URL
// Note : (no HTTP scheme validation)
func (opt *Options) validateURL(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" || s == " " || s == "\t" || s == "\n" {
		return s, false
	}
	return s, true
}

// Add wordlist(s) given by a path to a file or a folder that contain wordlists files.
// !Note : (The function do not validate the wordlist(s), if the wordlists is empty it wont return an error)
func (opt *Options) makeWordlist() error {
	//Check that the given wordlist was set:
	if len(opt.wordlistPath) == 0 {
		return nil
	}

	//Check if the wordlist given is a file or folder:
	typ, err := files.FileOrFolder(opt.wordlistPath)
	if err != nil {
		return err
	}

	switch typ {
	case "file":
		fSize, _ := files.FileSize(opt.wordlistPath)
		if fSize > 0 {
			opt.WordlistPaths = append(opt.WordlistPaths, opt.wordlistPath)
			return nil
		}
		err = errors.New("the wordlist file given is empty")

	case "folder":
		folder := opt.wordlistPath
		flst, _ := files.InDir(folder)
		for _, f := range flst {
			//Atleast one file has to be valid and inside the folder:
			pathToFile := (folder + "/" + f)
			if fSize, _ := files.FileSize(pathToFile); fSize > 0 && filepath.Ext(f) == ".txt" {
				opt.WordlistPaths = append(opt.WordlistPaths, pathToFile)
			}
		}
		if len(opt.WordlistPaths) > 0 {
			if opt.Verbose {
				log.Printf("Wordlists that will used: %v\n", opt.WordlistPaths)
			}
			return nil
		}
		err = errors.New("wordlist files inside the folder must have the extension \"txt\", otherwise ignored")
	}
	return err
}

// Update all provided URLs by appending all the schemes given (including the original scheme provided from the URLs)
func (opt *Options) makeURLs() error {
	//var l_urls []string
	for _, scheme := range opt.Scheme {
		for _, u := range opt.URLs {
			for _, method := range opt.Methods {
				Url := (scheme + "://" + u)
				hash := request.MakeHash(u, method)

				opt.Hosts[hash] = request.Host{
					URL:    Url,
					Method: method,
					Scheme: scheme,
				}
			}
		}
	}
	return nil
}

// Set HTTP Method(s) to be used within all future requests.
// all : preform all supported (*golang http package*) HTTP methods except DELETE
// all,delete : same as 'all' but include the DELETE method
func (opt *Options) setMethods(s string) error {
	all := "GET,HEAD,POST,PUT,TRACE,CONNECT"
	switch {
	case len(s) == 0:
		s = "GET"
	case strings.ToLower(s) == "all":
		s = all
	case strings.ToLower(s) == "all,delete":
		s = (all + ",DELETE")
	}
	opt.Methods = strings.Split(s, ",")
	return nil
}

func (opt *Options) setURLs_NoScheme(s string) error {
	for _, u := range splitEscape(s, ',') {
		//Delete the HTTP scheme in the URL (if any)
		if scheme := request.ContainScheme(u); scheme != "" {
			u = strings.Replace(u, (scheme + "://"), "", 1)

			//Add scheme to 'Option.Scheme' variable (if not already presented)
			if request.ValidHTTPScheme(scheme) && !slices.Contains(opt.Scheme, scheme) {
				opt.Scheme = append(opt.Scheme, scheme)
			}
		}
		if validURL, ok := opt.validateURL(u); ok && !slices.Contains(opt.URLs, u) {
			opt.URLs = append(opt.URLs, validURL)
		}

	}
	return nil
}

// Set HTTP scheme(s) that will be used within all future request
func (opt *Options) setScheme(s string) error {
	opt.Scheme = []string{}
	for _, scheme := range strings.Split(strings.ToLower(strings.TrimSpace(s)), ",") {
		if request.ValidHTTPScheme(scheme) {
			if !slices.Contains(opt.Scheme, scheme) {
				opt.Scheme = append(opt.Scheme, scheme)
			}
		} else {
			log.Fatalf("Unkown scheme: %s", scheme)
		}
	}
	return nil
}

// Set headers that will be used within all future request
// !Note : (all header names will be in transformed to lowercase)
func (opt *Options) setHeaders(s string) error {
	for _, headers := range splitEscape(s, ',') {
		if h := strings.SplitN(strings.TrimSpace(headers), ":", 2); len(h) == 2 {
			headerName := strings.ToLower(h[0])
			headerValue := h[1]
			opt.Headers = append(opt.Headers, [2]string{headerName, headerValue})
		} else {
			log.Fatalf("Invalid header can't make a header of value: %s", headers)
		}
	}
	return nil
}

// Extract the given HTTP request data string and set the data to the needed struct variables
func (opt *Options) setRaw(s string) error {
	var (
		Url       string
		Host      string
		HttpRaw   = strings.Split(s, "\n")
		Firstline = strings.SplitN(HttpRaw[0], " ", 3)
	)

	//Check if a url is set from the options
	if len(opt.URLs) > 1 {
		log.Fatalf("HTTP Raw request can only handle a single URL set: %v", opt.URLs)

	} else if len(opt.URLs) > 0 {
		Url = opt.URLs[0]
	}

	if len(Firstline) == 3 {
		opt.setMethods(strings.TrimSpace(Firstline[0]))                                     //Set Method
		opt.HTTP2, _ = regexp.MatchString("^HTTP/2(|.0)$", strings.TrimSpace(Firstline[2])) //Set Protocol

	} else {
		log.Fatal("Invalid HTTP Raw request given. Can't set Method, URL and/or Protocol from it.")
	}

	//Set Headers and Postdata:
	// Note : (Go reverse order to detect the postdata at the bottom of the HTTP Raw request)
	for i := 1; i <= len(HttpRaw)-1; i++ {
		v := HttpRaw[len(HttpRaw)-i] //<- (Header or Postdata)

		if len(opt.PostData) == 0 && (len(v) > 0 && (len(HttpRaw) > i && len(HttpRaw[len(HttpRaw)-(i+1)]) == 0)) {
			opt.PostData = v

		} else if len(v) > 0 {
			if h := strings.SplitN(v, ":", 2); len(h) == 2 {
				headerName := strings.ToLower(h[0])
				headerValue := h[1]
				opt.Headers = append(opt.Headers, [2]string{headerName, headerValue})

				//Check if it's the host header:
				if strings.ToLower(h[0]) == "host" {
					Host = h[1]
				}

			} else {
				log.Fatalf("Invalid header in HTTP Raw data. can't make a header of value: %s", v)
			}
		}
	}

	// Set URL:
	endpoint := strings.TrimSpace(Firstline[1])
	if len(Host) > 0 && len(opt.URLs) == 0 {
		//Check forward slashes to ignore double frontslashes (*unless it's by purpose*)
		if endpoint[0] == '/' && (len(Url) > 0 && Url[len(Url)-1] == '/') {
			endpoint = endpoint[1:]
		}
		opt.setURLs_NoScheme(Host + endpoint)
	}

	return nil
}

// Set the techniques to be used in the scanner process
func (opt *Options) setTechniques() error {
	//Reset the technique map:
	opt.Techniques = map[string]bool{
		"D": false,
		"E": false,
		"T": false,
		"X": false,
	}
	//Add technique related to user preference:
	for _, r := range strings.ToUpper(opt.technique) {
		techq := string(r)
		if _, ok := opt.Techniques[techq]; ok {
			opt.Techniques[techq] = true
		} else {
			return errors.New("Invalid technique was used: " + techq)
		}
	}
	return nil
}

// Detect and add insert points within the URLs GET parameters.
func (opt *Options) makeParams() error {
	_, cookieQuery := request.GetHeader(opt.Headers, "cookie")

	for hash, host := range opt.Hosts {
		//Add parameters to the related URL hash:
		parameter, err := parameter.NewParameter(opt.paramRules, opt.InsertKeyword)
		if err != nil {
			return err
		}
		for position, _ := range opt.paramRules {

			//Set the params in relation to their position for each host:
			switch position {
			case "url":
				if URLQuery, err := request.GetRawQuery(host.URL); len(URLQuery) > 0 && err == nil {
					parameter.SetURLparams(URLQuery)
				} else {
					return err
				}
			case "body":
				if len(opt.PostData) > 0 {
					parameter.SetBodyparams(opt.PostData)
				}
			case "cookie":
				if len(cookieQuery) > 0 {
					parameter.SetCookieparams(cookieQuery)
				}
			default:
				return errors.New("no valid position find for auto params")
			}

			//Set the parameters associated with the host (host hash)
			opt.Params[hash] = parameter
		}
	}
	return nil
}

func (opt *Options) setAutoParamRules(s string) error {
	var err error
	for _, paramRule := range splitEscape(s, ',') {
		paramRule = strings.ToLower(paramRule)

		if lst := splitEscape(paramRule, ':'); len(lst) == 2 || len(lst) == 3 {
			var ( //Extracted part of definitions from the core input value:
				position   = lst[0]
				method     = lst[1]
				separators = []rune{} //<- Optional
			)
			if len(lst) == 3 {
				//Append all the separators (chars/runes):
				for _, r := range lst[2] {
					separators = append(separators, r)
				}
			}
			//Validate and set rules in relation to the param position specified:
			//Note : Validate each defined value (separators are optional and can be anything):
			if ok, _ := regexp.MatchString(`^(replace|append|r|a)$`, method); !ok { //Method validation
				log.Fatal(design.STATUS.FAIL, " Invalid method used in auto detection of parameters, supported \"replace\" or \"append\" (r/a), Option: (\033[33m-au\033[0m)")
			}

			if ok, _ := regexp.MatchString(`^(url|body|cookie)$`, position); !ok { //Postion validation
				log.Fatal(design.STATUS.FAIL, " Invalid postion used in auto detection of parameters, supported \"url\", \"body\" or \"cookie\" (r/a), Option: (\033[33m-au\033[0m)")
			}
			//Set the auto parameter value for the option [struct]ure:
			opt.paramRules[position], err = parameter.NewRules(opt.InsertKeyword, method, separators)
			if err != nil {
				return err
			}

		} else {
			log.Fatal(design.STATUS.FAIL, " Syntax error on auto detection of parameter(s), Option: (\033[33m-au\033[0m)")
		}
	}
	return nil
}

func (opt *Options) setEncode(s string) error {
	opt.Encode = strings.Split(s, ",")
	return nil
}

// validate and setup the random keyword(s) - (%RANDOM%, %RANDOMSTR%, %RANDOMINT%) Return if the process failed/invalid input.
// Return the map as "nil" if an error was triggered in the process.
func (opt *Options) setRandomInsert(s string) error {
	if len(s) > 0 {
		for _, i := range strings.Split(s, ",") {
			//[Note] :  Only <= 3 loops so 'regexp.MatchString()' is fine to use here
			if ok, _ := regexp.MatchString(`^[sn]:[0-9]+$`, strings.ToLower(i)); ok {
				rules := strings.Split(i, ":")
				length, _ := strconv.Atoi(rules[1])
				opt.Random[rules[0]] = length

			} else {

				log.Fatalf("Invalid value for randominsert: %s. %v : Makes a random string with the length of 10.", s, exampleValues("s:10"))
			}
		}
	}
	return nil
}

func (opt *Options) showConfigOnScreen() {
	fmt.Println(strings.Repeat("_", 64), "\n\r")
	flag.VisitAll(func(f *flag.Flag) {
		var (
			lenName = len(fmt.Sprintf("%v", f.Value))
			Name    = fmt.Sprintf("%v", f.Name)
			value   string
		)
		if lenName > 0 && len(Name) > 2 {
			if Name == "raw" {
				value = "HTTP Raw"
			} else {
				value = fmt.Sprintf("%v", f.Value)
			}
			fmt.Printf(" - %s :: %s\n", strings.Title(f.Name), value)
		}
	})
	fmt.Println(strings.Repeat("_", 64), "\n\r")
}

// Display the supported encoders that can be used to payloads
func support_encodes() string {
	s := `1. url    [URL encode]
	2. durl   [DoubleURL encode]
	3. base64 [Base64 encode]
	4. base32 [Base32 encode]
	5. html   [HTML encode]
	6. htmle  [HTML Equivalent - HTML encode when " → &quot;]
	7. hex    [Hex encode]
	8. json   [Json encode]
	9. binary [Binary encode]`
	return s
}

func support_autoParameters() string {
	lst := []string{}
	for _, i := range parameter.SUPPORTED_PARAM_POSITIONS {
		lst = append(lst, "\033[1;33m"+i+"\033[0m")
	}
	return ("Support: (" + strings.Join(lst, ",") + ")")
}

func support_format(s string) string {
	return fmt.Sprintf("Format: (\033[1;33m%s\033[0m)", s)
}

// Split a string by comma but ignore escaped comma characters (\,) to be splitted.
// Return a string based list of all the items.
func splitEscape(s string, sep rune) []string {
	var (
		l   []string
		str string
	)
	for idx, r := range s {
		if len(s) >= 2 && r == sep {
			if s[idx-1] != '\\' {
				l = append(l, str)
				str = ""
				continue

			} else if s[idx-1] == '\\' {
				str = str[:len(str)-1]
			}
		}
		str += string(r)

		if idx == len(s)-1 {
			l = append(l, str)
		}
	}
	return l
}

func exampleValues(s string) string {
	return fmt.Sprintf("\033[1;33mEx\033[0m: ( %s )", s)
}
