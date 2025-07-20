package option

import (
	"fmt"
	"strings"
	"time"

	"github.com/Brum3ns/firefly/pkg/encode"
	"github.com/projectdiscovery/goflags"
)

type Option struct {
	Input       input
	Scan        Scan
	Process     process
	Http        http
	Match       match
	Filter      filter
	Skip        skip
	Knowledge   knowledge
	Randomness  randomness
	Payload     payload
	Extract     extract
	Placeholder placeholder
	Performance performance
	Output      output
}

var (
	DEFAULT_FILTER_DIFF_HEADER = []string{
		"Age",
		"Alt-Svc",
		"Cache-Control",
		"CF-Cache-Status",
		"CF-Ray",
		"Content-Length",
		"Content-Encoding",
		"Date",
		"ETag",
		"Expires",
		"Last-Modified",
		"Location",
		"Pragma",
		"Retry-After",
		"Server-Timing",
		"Set-Cookie",
		"Strict-Transport-Security",
		"Transfer-Encoding",
		"Vary",
		"Via",
		"X-Amz-Cf-Id",
		"X-Amz-Request-Id",
		"X-AspNet-Version",
		"X-Cache",
		"X-Cache-Status",
		"X-Cloud-Trace-Context",
		"X-Content-Type-Options",
		"X-Correlation-ID",
		"X-Edge-Location",
		"X-Fastly-Request-Id",
		"X-Frame-Options",
		"X-Powered-By",
		"X-Request-ID",
		"X-Response-Time",
		"X-Runtime",
		"X-Serve-Time",
		"X-Server",
		"X-Served-By",
		"X-Timer",
		"X-Trans-ID",
		"X-UA-Compatible",
		"X-Varnish",
		"X-Vercel-Cache",
		"X-Vercel-ID",
		"X-XSS-Protection",
	}
)

type input struct {
	RawHTTPRequest string `errcode:"IN-REQ-001"`
	Url            string `errcode:"IN-URL-001"`
}

type Scan struct {
	DisableExtract   bool
	DisableDiff      bool
	ResponseTimeDiff int // TODO
}
type process struct {
	Silence bool `errcode:"PRO-SILENCE-001"`
}

type output struct {
	OutputFile             string `errcode:"OUT-FILE-001"`
	OutputFileAnalyze      string `errcode:"OUT-FILE-004"`
	Overwrite              bool   `errcode:"OUT-FILE-002"`
	OutputHTTPResponseBody bool   `errcode:"OUT-FILE-003"`
	LogFile                string `errcode:"OUT-LOG-001"`
}

type http struct {
	Version                       string              `errcode:"HTTP-PROTO-001"`
	Methods                       goflags.StringSlice `errcode:"HTTP-METHOD-001"`
	Headers                       goflags.StringSlice `errcode:"HTTP-HEADER-001"`
	HeadersBrowserPreset          string              `errcode:"HTTP-HEADER-002"`
	URIPath                       string              `errcode:"HTTP-PATH-001"`
	Body                          string              `errcode:"HTTP-BODY-001"`
	Proxy                         string              `errcode:"HTTP-PROXY-001"`
	Timeout                       int                 `errcode:"HTTP-TIMEOUT-001"`
	Delay                         int
	FollowRedirect                bool `errcode:"HTTP-REDIRECT-001"`
	MaxRedirects                  int
	FollowHostRedirects           bool
	RespectHSTS                   bool
	DisableAutomaticHostHeader    bool
	DisableAutomaticContentLength bool
	ProxyDialTimeout              int
}

type match struct {
	Mode/*(OR|AND)*/ string                     `errcode:"Match-MODE-001"`
	StatusCode              goflags.StringSlice `errcode:"Match-CODE-001"`
	LineCount               goflags.StringSlice `errcode:"Match-LINE-001"`
	WordCount               goflags.StringSlice `errcode:"Match-WORD-001"`
	BodySize                goflags.StringSlice `errcode:"Match-SIZE-001"`
	ResponeTime             goflags.StringSlice `errcode:"Match-TIME-001"`
	MatchHeader             goflags.StringSlice `errcode:"Match-HEADER-001"`
	BodyRegex               string              `errcode:"Match-REBODY-001"`
	HeaderRegex             string              `errcode:"Match-REHEADER-001"`
}

type filter struct {
	Mode/*(OR|AND)*/ string                     `errcode:"FILTER-MODE-001"`
	StatusCode              goflags.StringSlice `errcode:"FILTER-CODE-001"`
	LineCount               goflags.StringSlice `errcode:"FILTER-LINE-001"`
	WordCount               goflags.StringSlice `errcode:"FILTER-WORD-001"`
	BodySize                goflags.StringSlice `errcode:"FILTER-SIZE-001"`
	ResponeTime             goflags.StringSlice `errcode:"FILTER-TIME-001"`
	FilterHeader            goflags.StringSlice `errcode:"FILTER-HEADER-001"`
	BodyRegex               string              `errcode:"FILTER-REBODY-001"`
	HeaderRegex             string              `errcode:"FILTER-REHEADER-001"`
	FilterDiffHeaders       goflags.StringSlice `errcode:"FILTER-DIFFHEADER-001"`
}

type skip struct {
	DiffBody      bool `errcode:"SKIP-DIFF-001"`
	DiffHeader    bool `errcode:"SKIP-DIFF-002"`
	ExtractHeader bool `errcode:"SKIP-EXTRACT-001"`
	ExtractBody   bool `errcode:"SKIP-EXTRACT-002"`
}

type payload struct {
	Placeholder  string              `errcode:"PAY-FUZZ-001"`
	Wordlist     string              `errcode:"PAY-WORDLIST-001"`
	VerifyCanary string              `errcode:"PAY-VERIFY-001"`
	Prefix       string              `errcode:"PAY-PREFIX-001"`
	ReflectStart int                 `errcode:"PAY-REFLECT-001"`
	ReflectEnd   int                 `errcode:"PAY-REFLECT-002"`
	Suffix       string              `errcode:"PAY-SUFFIX-001"`
	Replace      goflags.StringSlice `errcode:"PAY-REPLACE-001"`
	Encoders     goflags.StringSlice `errcode:"PAY-ENCODE-001"`
	EncoderParts goflags.StringSlice `errcode:"PAY-ENCODE-002"`
}

type extract struct {
	Wordlist string `errcode:"EXT-WORDLIST-001"`
	Regex    string `errcode:"EXT-REGEX-001"`
}

type knowledge struct {
	VerifyAmount int `errcode:"KNOWLEDGE-VERIFY-001"`
	Jitter       int `errcode:"KNOWLEDGE-JITTER-001"`
}

type randomness struct {
	EntropyValue int `errcode:"RANDOM-ENTROPY-001"`
}

type placeholder struct {
	List    bool   `errcode:"GEN-PLACEHOLDER-001"`
	Disable bool   `errcode:"GEN-DISABLE-001"`
	Test    string `errcode:"GEN-TEST-001"`
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
	flagSet.SetDescription("Firefly is an advanced black-box fuzzer and not just a standard asset discovery tool. Firefly provides the advantage of testing a target with a large number of built-in checks to detect behaviors in a web application.")

	flagSet.CreateGroup("input", "Input",
		flagSet.StringVarP(&opt.Input.RawHTTPRequest, "raw", "r", "", "raw HTTP request"),
		flagSet.StringVarP(&opt.Input.Url, "url", "u", "", "target url"),
	)

	flagSet.CreateGroup("silence", "Silence",
		flagSet.BoolVarP(&opt.Process.Silence, "silence", "s", false, "Do not display info during the fuzzing process"),
	)

	flagSet.CreateGroup("http request", "HTTP Request",
		flagSet.VarP(&opt.Http.Headers, "header", "H", "Output file in JSON format"),
		flagSet.StringVarP(&opt.Http.Version, "http-version", "hv", "1.1", "HTTP proto version (Supported: 0.9, 1.0, 1.1)"),
		flagSet.StringVarP(&opt.Http.HeadersBrowserPreset, "headers-preset", "Hp", "", "Preset headers that refer to a browser header order behavior. Available:[Chrome, Firefox or Safari]"),
		flagSet.StringSliceVarP(&opt.Http.Methods, "method", "X", []string{"GET"}, "HTTP Method(s) (file,comma-separated)", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringVar(&opt.Http.URIPath, "uripath", "/", "HTTP Request URI path"),
		flagSet.StringVarP(&opt.Http.Body, "data", "d", "", "HTTP Request body"),
		flagSet.IntVarP(&opt.Http.Timeout, "timeout", "T", 5000, "HTTP Timeout in milliseconds (ms)"),
		flagSet.IntVarP(&opt.Http.Delay, "delay", "D", 0, "Delay between HTTP requests in milliseconds (ms)"),
		flagSet.StringVar(&opt.Http.Proxy, "proxy", "", "HTTP Proxy"),
		flagSet.IntVarP(&opt.Http.ProxyDialTimeout, "proxy-timeout", "pT", 4000, "Proxy timeout in milliseconds (ms)"),
		flagSet.BoolVarP(&opt.Http.DisableAutomaticContentLength, "disable-auto-contentlength", "dacl", false, "Add automatic content-length in HTTP request"),
		flagSet.BoolVarP(&opt.Http.DisableAutomaticHostHeader, "disable-auto-host", "dah", false, "Add automatic host header in HTTP request"),
		//flagSet.BoolVarP(&opt.Http.FollowHostRedirects, "redirect-host", "rFH", false, "Follow host redirect"),
		flagSet.BoolVarP(&opt.Http.FollowRedirect, "redirect", "rF", false, "Follow redirect"),
		flagSet.IntVarP(&opt.Http.MaxRedirects, "redirect-max", "rM", 3, "Follow redirect"),
		flagSet.BoolVar(&opt.Http.RespectHSTS, "hsts", false, "Respect HTTP Strict Transport Security (HSTS)"),
	)

	flagSet.CreateGroup("diff filter", "Diff Filter",
		flagSet.StringSliceVar(&opt.Filter.FilterDiffHeaders, "fdH", DEFAULT_FILTER_DIFF_HEADER, "Filter hevy dynamic HTTP headers to avoid false positives when performing difference scans", goflags.CommaSeparatedStringSliceOptions),
	)

	// TODO : Make them work
	flagSet.CreateGroup("scanner", "Scanner",
		flagSet.BoolVarP(&opt.Scan.DisableExtract, "disable-extract", "dE", false, "Disable extract using regex/wordlist from the HTTP response"),
		flagSet.IntVarP(&opt.Scan.ResponseTimeDiff, "time-diff", "dT", 0, "HTTP Response Time difference to be counted as a suspisious behavior (used in timing attacks)"),
		flagSet.BoolVarP(&opt.Scan.DisableDiff, "disable-diff", "dD", false, "Disable diff checks between knowledge and fuzzed HTTP response results"),
	)

	flagSet.CreateGroup("http filter", "HTTP Filter",
		flagSet.StringVar(&opt.Filter.Mode, "fMode", "or", "Filter mode"),
		flagSet.StringSliceVar(&opt.Filter.StatusCode, "fc", []string{}, "Filter HTTP status code", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Filter.LineCount, "fl", []string{}, "Filter lines in HTTP response", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Filter.WordCount, "fe", []string{}, "Filter word", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Filter.BodySize, "fs", []string{}, "Filter HTTP body size", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Filter.ResponeTime, "ft", []string{}, "Filter HTTP response time", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Filter.FilterHeader, "fh", []string{}, "Filter HTTP header name", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringVar(&opt.Filter.BodyRegex, "fbr", "", "Filter response body by regex"),
		flagSet.StringVar(&opt.Filter.HeaderRegex, "fhr", "", "Filter HTTP header by regex"),
	)

	flagSet.CreateGroup("http match", "HTTP Match",
		flagSet.StringVar(&opt.Match.Mode, "mMode", "or", "Match mode"),
		flagSet.StringSliceVar(&opt.Match.StatusCode, "mc", []string{}, "Match HTTP status code", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Match.LineCount, "ml", []string{}, "Match lines in HTTP response", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Match.WordCount, "mw", []string{}, "Match word", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Match.BodySize, "ms", []string{}, "Match HTTP body size", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Match.ResponeTime, "mt", []string{}, "Match HTTP response time", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&opt.Match.MatchHeader, "mh", []string{}, "Match HTTP header name", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringVar(&opt.Match.BodyRegex, "mrb", "", "Match response body by regex"),
		flagSet.StringVar(&opt.Match.HeaderRegex, "mhr", "", "Match HTTP header by regex"),
	)

	flagSet.CreateGroup("skip", "Skip",
		flagSet.BoolVarP(&opt.Skip.DiffHeader, "sdh", "skip-diff-header", false, "Skip diff scan in HTTP response headers"),
		flagSet.BoolVarP(&opt.Skip.DiffBody, "sdb", "skip-diff-body", false, "Skip diff scan in HTTP response body"),
		flagSet.BoolVarP(&opt.Skip.ExtractHeader, "seh", "skip-extract-header", false, "Skip to extract keywords in the HTTP response headers"),
		flagSet.BoolVarP(&opt.Skip.ExtractBody, "seb", "skip-extract-body", false, "Skip to extract keywords in the HTTP response body"),
	)

	flagSet.CreateGroup("knowledge", "Knowledge",
		flagSet.IntVarP(&opt.Knowledge.VerifyAmount, "k", "verify-knowledge", 11, "Amount of request to perform with the verify payload to fetch target knowledge"),
		flagSet.IntVarP(&opt.Knowledge.Jitter, "j", "jitter", 0, "Jitter delay in millisecond (ms) between verify requests to get more accurate HTTP responses stored in the knowledge"),
	)

	flagSet.CreateGroup("randomness", "Randomness",
		flagSet.IntVarP(&opt.Randomness.EntropyValue, "eV", "entropy", 3500, "Minimum entropy value to counts as a random value (noise) in the HTTP respons. (Ex: 3500 => 3.5 entropy)"),
	)

	flagSet.CreateGroup("payload", "Payload",
		flagSet.StringVar(&opt.Payload.Placeholder, "i", "FUZZ", "Placeholder to insert payload"),
		flagSet.StringVarP(&opt.Payload.Wordlist, "w", "wl-payload", "tests/wordlists/wordlist.txt", "Payload wordlist"),
		flagSet.StringVar(&opt.Payload.VerifyCanary, "pv", "8Ajg1al2", "Payload canary to verify target knowledge"),
		flagSet.StringVar(&opt.Payload.Suffix, "ps", "", "Suffix that will be added to all payloads in the given wordlist"),
		flagSet.StringVar(&opt.Payload.Prefix, "px", "", "Prefix that will be added to all payloads in the given wordlist"),
		flagSet.IntVar(&opt.Payload.ReflectStart, "rs", 6, "Payload surrounding length start"),
		flagSet.IntVar(&opt.Payload.ReflectEnd, "re", 6, "Payload surrounding length end"),
		flagSet.StringSliceVar(&opt.Payload.Replace, "pr", []string{}, "Replace values inside payloads in the given wordlist, separeted by comma", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVarP(&opt.Payload.EncoderParts, "ep", "encode-parts", []string{}, "Parts to be encoded within the payloads", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVarP(&opt.Payload.Encoders, "e", "encode", []string{}, fmt.Sprintf("Encoders to encode each payload with in order, separated with comma.Supported encoders:\n%v\n", strings.Join(encode.GetEncoders(), "\n")), goflags.CommaSeparatedStringSliceOptions),
	)

	flagSet.CreateGroup("extract", "Extract",
		flagSet.StringVarP(&opt.Extract.Regex, "eR", "extract-regex", "tests/wordlists/wordlist-regex.txt", "Wordlist of keywords that contains patters to extract in the HTTP response"),
		flagSet.StringVarP(&opt.Extract.Wordlist, "eW", "extract-wordlist", "tests/wordlists/wordlist-pattern.txt", "Wordlist of regexes that contains patters to extract in the HTTP response"),
	)

	flagSet.CreateGroup("placeholder", "Placeholder",
		flagSet.BoolVarP(&opt.Placeholder.List, "lP", "list-placeholders", false, "List all supported placeholder and exit"),
		flagSet.BoolVarP(&opt.Placeholder.Disable, "dP", "disable-placeholder", false, "Disable placeholder insert"),
		flagSet.StringVarP(&opt.Placeholder.Test, "tP", "test-placeholder", "", "Test a placeholder"),
	)

	flagSet.CreateGroup("performance", "Performance",
		flagSet.IntVarP(&opt.Performance.ThreadsRequest, "tR", "threads-req", 50, "Threads used within the request's worker pool"),
		flagSet.IntVarP(&opt.Performance.ThreadsScanner, "tS", "threads-scan", 20, "Threads used within the scanner's worker pool"),
		flagSet.IntVarP(&opt.Performance.BufferPoolRequest, "bR", "buffer-req", 3000, "Buffer pool for request handler"),
		flagSet.IntVarP(&opt.Performance.BufferPoolScanner, "bS", "buffer-scan", 2000, "Buffer pool for scanner handler"),
		flagSet.IntVarP(&opt.Performance.BufferPoolResponse, "bRs", "buffer-resp", 3000, "Buffer pool for response handler"),
		flagSet.IntVarP(&opt.Performance.BufferPoolKnowledge, "bk", "buffer-knowledge", 2000, "Buffer pool for knowledge handler"),
		flagSet.IntVarP(&opt.Performance.BufferPoolResult, "bRt", "buffer-result", 2000, "Buffer pool for result handler"),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.BoolVarP(&opt.Output.Overwrite, "output-overwrite", "oW", false, "Overwrite output file"),
		flagSet.StringVarP(&opt.Output.OutputFile, "output", "o", "", "Output file in JSON format"),
		flagSet.StringVarP(&opt.Output.OutputFileAnalyze, "output-analyze", "oA", "", "Output file for payload analyze in JSON format"),
		flagSet.StringVar(&opt.Output.LogFile, "log", getLogFile(), "File to store all logs"),
		flagSet.BoolVarP(&opt.Output.OutputHTTPResponseBody, "output-responsebody", "orb", false, "Save full HTTP response in output"),
	)

	if err := flagSet.Parse(); err != nil {
		return opt, fmt.Errorf("could not parse flags: %v", err)
	}
	return Validate(opt)
}

func getLogFile() string {
	return fmt.Sprintf("/tmp/firefly/firefly-%s.log", time.Now().Format("2006-01-02-15-04-05"))
}
