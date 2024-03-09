package config

import (
	"log"

	"github.com/Brum3ns/firefly/internal/global"
	"github.com/Brum3ns/firefly/internal/option"
	"github.com/Brum3ns/firefly/pkg/extract"
	"github.com/Brum3ns/firefly/pkg/filter"
	"github.com/Brum3ns/firefly/pkg/payloads"
	"github.com/Brum3ns/firefly/pkg/transformation"
)

type Configure struct {
	*option.Options
	*payloads.Wordlist
	*filter.Filter
	*Scanner
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
	//Difference         difference.Difference//<-Not needed ATM
}

func NewConfigure(opt *option.Options) *Configure {
	wl_transformation, err := transformation.GetWordlist(opt.TransformationYAMLFile)
	if err != nil {
		log.Fatal(err)
	}

	conf := &Configure{
		Options: opt,
		Filter:  filter.NewFilter(opt),
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
	conf.Scanner = conf.newScanner()

	return conf

}

func (conf *Configure) newScanner() *Scanner {
	//Setup scanner technique resources:
	wlPtn, wlRegex := extract.MakeWordlists(global.DIR_DETECTION)
	wlPatternPrefix, wlPatterns := extract.CreatePrefixMap(wlPtn)

	transform, err := transformation.NewTransformation(conf.TransformationYAMLFile)
	if err != nil {
		log.Fatal(err)
	}

	return &Scanner{
		OK_Extract:         conf.Options.Techniques["E"],
		OK_Diff:            conf.Options.Techniques["D"],
		OK_Transformation:  conf.Options.Techniques["T"],
		DisablesTechniques: conf.Options.Techniques["X"],

		Extract: extract.NewExtract(extract.Properties{
			Threads:         conf.Options.ThreadsExtract,
			PrefixPatterns:  wlPatternPrefix,
			WordlistPattern: wlPatterns,
			WordlistRegex:   map[string][]string{extract.WILDCARD: wlRegex},
		}),
		Transformation: transform,

		//Difference: *difference.NewDifference(difference.Properties{}),//<-Not needed ATM
	}
}
