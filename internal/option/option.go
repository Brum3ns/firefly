package option

import (
	"fmt"
	"time"

	"github.com/projectdiscovery/goflags"
)

type Option struct {
	input
	http
	filter
	knowledge
	payload
	placeholder
	performance
	output
}

type input struct {
	RawHTTPRequest string `errcode:"IN-REQ-001"`
	Url            string `errcode:"IN-URL-001"`
}

type output struct {
	Output  string `errcode:"OUT-FILE-001"`
	LogFile string `errcode:"OUT-LOG-001"`
}

type http struct {
	Methods                goflags.StringSlice `errcode:"HTTP-METHOD-001"`
	Headers                goflags.StringSlice `errcode:"HTTP-HEADER-001"`
	URIPath                string              `errcode:"HTTP-PATH-001"`
	Body                   string              `errcode:"HTTP-BODY-001"`
	Proxy                  string              `errcode:"HTTP-PROXY-001"`
	Timeout                int                 `errcode:"HTTP-TIMEOUT-001"`
	HttpDelay              int
	FollowRedirect         bool `errcode:"HTTP-REDIRECT-001"`
	MaxRedirects           int
	FollowHostRedirects    bool
	RespectHSTS            bool
	AutomaticHostHeader    bool
	AutomaticContentLength bool
	ProxyDialTimeout       int
}

type filter struct {
	MatchMode/*(OR|AND)*/ string                      `errcode:"FILTER-MODE-001"`
	MatchStatusCode               goflags.StringSlice `errcode:"FILTER-CODE-001"`
	MatchLineCount                goflags.StringSlice `errcode:"FILTER-LINE-001"`
	MatchWordCount                goflags.StringSlice `errcode:"FILTER-WORD-001"`
	MatchSize                     goflags.StringSlice `errcode:"FILTER-SIZE-001"`
	MatchTime                     goflags.StringSlice `errcode:"FILTER-TIME-001"`
	MatchHeader                   goflags.StringSlice `errcode:"FILTER-HEADER-001"`
	MatchBodyRegex                string              `errcode:"FILTER-REBODY-001"`
	MatchHeaderRegex              string              `errcode:"FILTER-REHEADER-001"`
	FilterMode/*(OR|AND)*/ string                     `errcode:"FILTER-MODE-002"`
	FilterStatusCode              goflags.StringSlice `errcode:"FILTER-CODE-002"`
	FilterLineCount               goflags.StringSlice `errcode:"FILTER-LINE-002"`
	FilterWordCount               goflags.StringSlice `errcode:"FILTER-WORD-002"`
	FilterSize                    goflags.StringSlice `errcode:"FILTER-SIZE-002"`
	FilterTime                    goflags.StringSlice `errcode:"FILTER-TIME-002"`
	FilterHeader                  goflags.StringSlice `errcode:"FILTER-HEADER-002"`
	FilterBodyRegex               string              `errcode:"FILTER-REBODY-002"`
	FilterHeaderRegex             string              `errcode:"FILTER-REHEADER-002"`
	//filterDiffHeader string   `flag:"fdH" errorcode:"3019"`
	//FilterDiffHeader []string `flag:"fdH" errorcode:"3019"`
}

type payload struct {
	FuzzPlaceholder string `flag:"w" errcode:"PAY-FUZZ-001"`
	PayloadWordlist string `flag:"w" errcode:"PAY-WORDLIST-001"`
	PayloadVerify   string `flag:"w" errcode:"PAY-VERIFY-001"`
}

type knowledge struct {
	VerifyKnowledge int `flag:"w" errcode:"KNOWLEDGE-VERIFY-001"`
}

type placeholder struct {
	ListPlaceholders    bool   `flag:"w" errcode:"GEN-PLACEHOLDER-001"`
	DisablePlaceholders bool   `flag:"w" errcode:"GEN-DISABLE-001"`
	TestPlaceholders    string `flag:"w" errcode:"GEN-TEST-001"`
}

type performance struct {
	ThreadsRequest      int
	ThreadsScanner      int
	BufferPoolRequest   int
	BufferPoolScanner   int
	BufferPoolResponse  int
	BufferPoolKnowledge int
	BufferPoolResult    int
}

func NewOption() (Option, error) {
	opt := Option{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription("Firefly is an advanced black-box fuzzer and not just a standard asset discovery tool. Firefly provides the advantage of testing a target with a large number of built-in checks to detect behaviors in the target.")

	// Group example
	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&opt.RawHTTPRequest, "raw", "r", "", "raw HTTP request"),
		flagSet.StringVarP(&opt.Url, "url", "u", "", "target url"),
	)

	flagSet.CreateGroup("http request", "HTTP Request",
		flagSet.VarP(&opt.Headers, "header", "H", "Output file in JSON format"),
		flagSet.StringSliceVarP(&opt.Methods, "method", "X", []string{"GET"}, "HTTP Method(s) (file,comma-separated)", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringVar(&opt.URIPath, "uripath", "/", "HTTP Request URI path"),
		flagSet.StringVarP(&opt.Body, "data", "d", "", "HTTP Request body"),
		flagSet.IntVarP(&opt.Timeout, "timeout", "T", 5000, "HTTP Timeout in milliseconds (ms)"),
		flagSet.IntVarP(&opt.HttpDelay, "delay", "D", 0, "Delay between HTTP requests in milliseconds (ms)"),
		flagSet.StringVar(&opt.Proxy, "proxy", "", "HTTP Proxy"),

		flagSet.IntVarP(&opt.ProxyDialTimeout, "proxy-timeout", "pT", 4000, "Proxy timeout in milliseconds (ms)"),

		flagSet.BoolVar(&opt.AutomaticContentLength, "aCL", true, "Add automatic content-length in HTTP request"),
		flagSet.BoolVar(&opt.AutomaticHostHeader, "aH", true, "Add automatic host header in HTTP request"),

		flagSet.BoolVarP(&opt.FollowHostRedirects, "redirect-host", "rFH", false, "Follow host redirect"),
		flagSet.BoolVarP(&opt.FollowRedirect, "redirect", "rF", false, "Follow redirect"),
		flagSet.IntVarP(&opt.MaxRedirects, "redirect-max", "rM", 3, "Follow redirect"),

		flagSet.BoolVar(&opt.RespectHSTS, "hsts", false, "Respect HTTP Strict Transport Security (HSTS)"),
	)

	flagSet.CreateGroup("filter", "Filter",
		flagSet.StringVar(&opt.FilterMode, "fMode", "or", "Filter mode"),
		flagSet.StringSliceVar(&opt.FilterStatusCode, "fc", []string{}, "Filter HTTP status code", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.FilterLineCount, "fl", []string{}, "Filter lines in HTTP response", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.FilterWordCount, "fe", []string{}, "Filter word", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.FilterSize, "fd", []string{}, "Filter HTTP body size", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.FilterTime, "ft", []string{}, "Filter HTTP response time", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.FilterHeader, "fh", []string{}, "Filter HTTP header name", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringVar(&opt.FilterBodyRegex, "fbr", "", "Filter response body by regex"),
		flagSet.StringVar(&opt.FilterHeaderRegex, "fhr", "", "Filter HTTP header by regex"),
	)

	flagSet.CreateGroup("match", "Match",
		flagSet.StringVar(&opt.MatchMode, "mMode", "or", "Match mode"),
		flagSet.StringSliceVar(&opt.MatchStatusCode, "mc", []string{}, "Match HTTP status code", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.MatchLineCount, "ml", []string{}, "Match lines in HTTP response", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.MatchWordCount, "mw", []string{}, "Match word", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.MatchSize, "ms", []string{}, "Match HTTP body size", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.MatchTime, "mt", []string{}, "Match HTTP response time", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.MatchHeader, "mh", []string{}, "Match HTTP header name", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringVar(&opt.MatchBodyRegex, "mrb", "", "Match response body by regex"),
		flagSet.StringVar(&opt.MatchHeaderRegex, "mhr", "", "Match HTTP header by regex"),
	)

	flagSet.CreateGroup("knowledge", "Knowledge",
		flagSet.IntVarP(&opt.VerifyKnowledge, "vK", "verify-knowledge", 3, "Amount of request to perform with the verify payload to fetch target knowledge"),
	)

	flagSet.CreateGroup("payload", "Payload",
		flagSet.StringVar(&opt.FuzzPlaceholder, "i", "FUZZ", "Placeholder to insert payload"),
		flagSet.StringVarP(&opt.PayloadWordlist, "wl-payload'", "w", "", "Payload wordlist"),
		flagSet.StringVar(&opt.PayloadVerify, "pV", "8Ajg1al2", "Payload verify canary"),
	)

	flagSet.CreateGroup("placeholder", "Placeholder",
		flagSet.BoolVarP(&opt.ListPlaceholders, "lP", "list-placeholders", false, "List all supported placeholder and exit"),
		flagSet.BoolVarP(&opt.DisablePlaceholders, "dP", "disable-placeholder", false, "Disable placeholder insert"),
		flagSet.StringVarP(&opt.TestPlaceholders, "tP", "test-placeholder", "", "Test a placeholder"),
	)

	flagSet.CreateGroup("performance", "Performance",
		flagSet.IntVarP(&opt.ThreadsRequest, "tR", "threads-req", 50, "Threads used within the request's worker pool"),
		flagSet.IntVarP(&opt.ThreadsScanner, "tS", "threads-scan", 20, "Threads used within the scanner's worker pool"),
		flagSet.IntVarP(&opt.BufferPoolRequest, "bR", "buffer-req", 3000, "Buffer pool for request handler"),
		flagSet.IntVarP(&opt.BufferPoolScanner, "bS", "buffer-scan", 2000, "Buffer pool for scanner handler"),
		flagSet.IntVarP(&opt.BufferPoolResponse, "bRs", "buffer-resp", 3000, "Buffer pool for response handler"),
		flagSet.IntVarP(&opt.BufferPoolKnowledge, "bk", "buffer-knowledge", 2000, "Buffer pool for knowledge handler"),
		flagSet.IntVarP(&opt.BufferPoolResult, "bRt", "buffer-result", 2000, "Buffer pool for result handler"),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&opt.Output, "output", "o", "", "Output file in JSON format"),
		flagSet.StringVar(&opt.LogFile, "log", getLogFile(), "File to store all logs"),
	)

	if err := flagSet.Parse(); err != nil {
		return opt, fmt.Errorf("could not parse flags: %v", err)
	}

	return Validate(opt)
}

func getLogFile() string {
	return fmt.Sprintf("/tmp/firefly/firefly-%s.log", time.Now().Format("2006-01-02-15-04-05"))
}
