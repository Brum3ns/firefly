package config

import (
	"log"

	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/internal/option"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/functions"
	"github.com/Brum3ns/firefly/pkg/httpfilter"
	"github.com/Brum3ns/firefly/pkg/payloads"
	"github.com/Brum3ns/firefly/pkg/randomness"
	"github.com/Brum3ns/firefly/pkg/request"
	"github.com/Brum3ns/firefly/pkg/transformation"
)

type Configure struct {
	Httpfilter httpfilter.Filter
	HttpMatch  httpfilter.Filter
	Option     *option.Options
	Wordlist   *payloads.Wordlist
	Scanner    *Scanner
}

// Scanner properties (static storage)
// Note : (This structure should not be modified once it's defined).
type Scanner struct {
	OK_Extract         bool
	OK_Diff            bool
	OK_Transformation  bool
	DisablesTechniques bool
	Extract            extract.Extract
	Transformation     transformation.Transformation
	Randomness         randomness.Randomness
	HttpDiffFilter     httpfilter.Filter
}

func NewConfigure(opt *option.Options) (*Configure, error) {
	wl_transformation, err := transformation.GetWordlist(opt.TransformationYAMLFile)
	if err != nil {
		log.Fatal(err)
	}

	// Configure HTTP filter
	filter, err := httpfilter.NewFilter(httpfilter.Config{
		Mode:                  opt.FilterMode,
		HeaderRegex:           opt.FilterHeaderRegex,
		BodyRegex:             opt.FilterBodyRegex,
		StatusCodes:           functions.SplitEscape(opt.FilterCode, ','),
		WordCounts:            functions.SplitEscape(opt.FilterWord, ','),
		LineCounts:            functions.SplitEscape(opt.FilterLine, ','),
		ResponseSizes:         functions.SplitEscape(opt.FilterSize, ','),
		ResponseTimesMillisec: functions.SplitEscape(opt.FilterTime, ','),
		Header:                request.LstToHeaders(LstToKeyMap(functions.SplitEscape(opt.FilterHeader, ','))),
	})
	if err != nil {
		return &Configure{}, err
	}

	// Configure HTTP match filter
	match, err := httpfilter.NewFilter(httpfilter.Config{
		Mode:                  opt.MatchMode,
		HeaderRegex:           opt.MatchHeaderRegex,
		BodyRegex:             opt.MatchBodyRegex,
		StatusCodes:           functions.SplitEscape(opt.MatchCode, ','),
		WordCounts:            functions.SplitEscape(opt.MatchWord, ','),
		LineCounts:            functions.SplitEscape(opt.MatchLine, ','),
		ResponseSizes:         functions.SplitEscape(opt.MatchSize, ','),
		ResponseTimesMillisec: functions.SplitEscape(opt.MatchTime, ','),
		Header:                request.LstToHeaders(LstToKeyMap(functions.SplitEscape(opt.MatchHeader, ','))),
	})
	if err != nil {
		return &Configure{}, err
	}

	conf := &Configure{
		Option:     opt,
		Httpfilter: filter,
		HttpMatch:  match,
		Wordlist: payloads.NewWordlist(
			&payloads.Wordlist{
				Files:              opt.WordlistPaths,
				TransformationList: wl_transformation,
				Verify: payloads.Verify{
					Payload: opt.VerifyPayload,
					Amount:  opt.VerifyAmount,
				},
				PayloadProperties: payloads.PayloadProperties{
					Tamper:         opt.Tamper,
					Encode:         opt.Encode,
					PayloadPattern: opt.PayloadPattern,
					PayloadPrefix:  opt.PayloadPrefix,
					PayloadSuffix:  opt.PayloadSuffix,
					PayloadReplace: opt.PayloadReplace,
				},
			},
		),
	}

	//Return a *pointer* of the "Scanner" [struct]ure:
	conf.Scanner, err = conf.newScanner()
	if err != nil {
		return &Configure{}, err
	}

	return conf, nil

}

func (conf *Configure) newScanner() (*Scanner, error) {
	//Setup scanner technique resources:
	wlPtn, wlRegex := extract.MakeWordlists(global.DIR_DETECTION)
	wlPatternPrefix, wlPatterns := extract.CreatePrefixMap(wlPtn)

	rand, err := randomness.NewRandomness(randomness.DefaultConfig())
	if err != nil {
		return &Scanner{}, err
	}

	// Configure the transformation structure
	transform, err := transformation.NewTransformation(conf.Option.TransformationYAMLFile)
	if err != nil {
		return &Scanner{}, err
	}

	// Configure HTTP difference filter
	httpDiffFilter, err := httpfilter.NewFilter(httpfilter.Config{
		Header: request.LstToHeaders(LstToKeyMap(functions.SplitEscape(conf.Option.FilterDiffHeader, ','))),
	})
	if err != nil {
		return &Scanner{}, err
	}

	return &Scanner{
		OK_Extract:         conf.Option.Techniques["E"],
		OK_Diff:            conf.Option.Techniques["D"],
		OK_Transformation:  conf.Option.Techniques["T"],
		DisablesTechniques: conf.Option.Techniques["X"],

		Extract: extract.NewExtract(extract.Properties{
			Threads:         conf.Option.ThreadsExtract,
			PrefixPatterns:  wlPatternPrefix,
			WordlistPattern: wlPatterns,
			WordlistRegex:   map[string][]string{extract.WILDCARD: wlRegex},
		}),
		Transformation: transform,
		Randomness:     rand,
		HttpDiffFilter: httpDiffFilter,

		//Difference: *difference.NewDifference(difference.Properties{}),//<-Not needed ATM
	}, nil
}

func LstToKeyMap(lst []string) map[string]string {
	var m = make(map[string]string)
	for _, i := range lst {
		m[i] = ""
	}
	return m
}
