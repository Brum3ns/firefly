package runner

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Brum3ns/firefly/internal/knowledge"
	"github.com/Brum3ns/firefly/internal/option"
	"github.com/Brum3ns/firefly/internal/scan"
	"github.com/panjf2000/ants/v2"

	"github.com/Brum3ns/firefly/pkg/httpfilter"
	"github.com/Brum3ns/firefly/pkg/payload"
	"github.com/Brum3ns/firefly/pkg/randomness"
	"github.com/Brum3ns/firefly/pkg/rhttp"
	"github.com/Brum3ns/firefly/pkg/waitgroup"
	"github.com/projectdiscovery/rawhttp"
)

const (
	mode_knowledge = "knowledge"
	mode_fuzz      = "fuzz"
)

type Runner struct {
	mode       string
	wg         waitGroup
	channel    channel
	payload    payload.Payload
	scan       scan.Scan
	request    rhttp.Http
	option     option.Option
	httpfilter httpfilter.Filter
	httpmatch  httpfilter.Filter
	randomness randomness.Randomness
	statistic  Statistic
	workerpool workerpool
	knowledge  knowledge.Knowledge
}

type workerpool struct {
	request *ants.Pool
	scanner *ants.Pool
}

type waitGroup struct {
	process      waitgroup.WaitGroup
	knowledge    waitgroup.WaitGroup
	result       waitgroup.WaitGroup
	scanner      waitgroup.WaitGroup
	httpRequest  waitgroup.WaitGroup
	httpResponse waitgroup.WaitGroup
}

type channel struct {
	httpRequest  chan jobHTTP
	scanner      chan jobScanner
	httpResponse chan jobHTTP
	knowledge    chan jobScanner
	result       chan jobScanner
}

func NewRunner(opt option.Option) (Runner, error) {
	var err error
	var r = Runner{
		option:    opt,
		knowledge: knowledge.NewKnowledge(),
		channel: channel{
			scanner:      make(chan jobScanner, opt.BufferPoolScanner),
			httpRequest:  make(chan jobHTTP, opt.BufferPoolRequest),
			httpResponse: make(chan jobHTTP, opt.BufferPoolRequest),
			knowledge:    make(chan jobScanner, opt.BufferPoolKnowledge),
			result:       make(chan jobScanner, opt.BufferPoolResult),
		},
	}

	// Set the payload instance
	if r.payload, err = payload.NewPayload(payload.Config{
		WordlistFile:  opt.PayloadWordlist,
		PayloadVerify: opt.PayloadVerify,
	}); err != nil {
		return r, err
	}

	// Configure HTTP request
	if r.request, err = rhttp.NewRequest(rhttp.Config{
		Methods:     opt.Methods,
		Headers:     rhttp.MakeHeaders(opt.Headers),
		Url:         opt.Url,
		URIPath:     opt.URIPath,
		Body:        opt.Body,
		RespectHSTS: opt.RespectHSTS,
	},
		&rawhttp.Options{
			Timeout:                time.Duration(opt.Timeout) * time.Millisecond,
			FollowRedirects:        opt.FollowRedirect,
			MaxRedirects:           3,
			AutomaticHostHeader:    opt.AutomaticHostHeader,
			AutomaticContentLength: opt.AutomaticContentLength,
			Proxy:                  opt.Proxy,
			ProxyDialTimeout:       time.Duration(opt.ProxyDialTimeout) * time.Millisecond,
		},
		/* &rawhttp.Options{
			Timeout:                7000 * time.Millisecond,
			FollowRedirects:        opt.FollowHostRedirects,
			MaxRedirects:           opt.MaxRedirects,
			AutomaticHostHeader:    opt.AutomaticHostHeader,
			AutomaticContentLength: opt.AutomaticContentLength,
			//CustomHeaders:          opt.Headers,
			ForceReadAllBody: true,
			//CustomRawBytes:         "",
			Proxy:            opt.Proxy,
			ProxyDialTimeout: 7000 * time.Millisecond,
			//SNI:                    "",
			//FastDialer:             "",
		}, */
	); err != nil {
		return r, err
	}

	// Configure HTTP filter
	r.httpfilter, err = httpfilter.NewFilter(httpfilter.Config{
		Mode:                  opt.FilterMode,
		HeaderRegex:           opt.FilterHeaderRegex,
		BodyRegex:             opt.FilterBodyRegex,
		StatusCodes:           opt.FilterStatusCode,
		WordCounts:            opt.FilterWordCount,
		LineCounts:            opt.FilterLineCount,
		ResponseSizes:         opt.FilterSize,
		ResponseTimesMillisec: opt.FilterTime,
		//TODO : //Header:                request.LstToHeaders(LstToKeyMap(opt.FilterHeader)),
	})
	if err != nil {
		return r, err
	}

	// Configure HTTP match filter
	r.httpmatch, err = httpfilter.NewFilter(httpfilter.Config{
		Mode:                  opt.MatchMode,
		HeaderRegex:           opt.MatchHeaderRegex,
		BodyRegex:             opt.MatchBodyRegex,
		StatusCodes:           opt.MatchStatusCode,
		WordCounts:            opt.MatchWordCount,
		LineCounts:            opt.MatchLineCount,
		ResponseSizes:         opt.MatchSize,
		ResponseTimesMillisec: opt.MatchTime,
		//TODO : //Header:                request.LstToHeaders(LstToKeyMap(opt.MatchHeader)),
	})
	if err != nil {
		return r, err
	}

	// Setup randomness config
	r.randomness, err = randomness.NewRandomness(randomness.Config{
		InRow:      randomness.DEFAULT_INROW,
		Vocal:      randomness.DEFAULT_VOCAL,
		Digit:      randomness.DEFAULT_DIGIT,
		Consonant:  randomness.DEFAULT_CONSONANT,
		Blacklist:  randomness.DEFAULT_BLACKLISTS,
		BlackRegex: randomness.DEFAULT_BLACKREGEX,
		Spaces:     []rune{' ', '_', '-', '.'},
	})
	if err != nil {
		log.Println(err)
	}

	// Define worker pools
	r.workerpool.request, err = ants.NewPool(opt.ThreadsRequest)
	r.workerpool.scanner, err = ants.NewPool(opt.ThreadsScanner)

	// Make the runner
	if err := r.make(); err != nil {
		return r, err
	}

	return r, nil
}

func (r *Runner) make() error {
	if err := r.request.MakeCoreRequests(); err != nil {
		return err
	}
	if err := r.payload.MakeWordlist(); err != nil {
		return err
	}
	return nil
}

func (r *Runner) FetchKnowledge() error {
	r.setMode(mode_knowledge)
	err := r.run()
	// Merge the knowledge
	if r.mode == mode_knowledge {
		r.knowledge.SetMergeKnowledge()
	}
	return err
}

func (r *Runner) RunFuzz() (Statistic, error) {
	if !r.hasKnowledge() {
		return r.statistic, errors.New("can not run fuzz, no knowledge of target")
	}

	r.setMode(mode_fuzz)
	if err := r.run(); err != nil {
		return r.statistic, err
	}
	return r.statistic, nil
}

func (r *Runner) run() error {
	// Create the log file to store all logs in
	/* fileLog, err := makeLog(r.Option.LogFile)
	if err != nil {
		return err
	}
	defer fileLog.Close() */

	// TODO : Move to each global run func and mode knowledge handler to its run
	ctx, cancel := context.WithCancel(context.Background())

	// Start handlers
	if r.mode == mode_knowledge {
		r.wg.process.Add(1)
		go r.handleKnowledge(ctx)
	}
	r.wg.process.Add(4)
	go r.handlerJobRequest(ctx)
	go r.handleJobScanner(ctx)
	go r.handlerResponse(ctx)
	go r.handlerResult(ctx)

	r.sendCoreJobs()
	r.waitJobs()

	// safely kill all running processes
	cancel()
	r.wg.process.Wait()

	return nil
}

func (r *Runner) waitJobs() {
	for {
		if !r.wg.httpRequest.HasJob() && len(r.channel.httpRequest) == 0 &&
			!r.wg.httpResponse.HasJob() && len(r.channel.httpResponse) == 0 &&
			!r.wg.scanner.HasJob() && len(r.channel.scanner) == 0 &&
			!r.wg.knowledge.HasJob() && len(r.channel.knowledge) == 0 &&
			!r.wg.result.HasJob() && len(r.channel.result) == 0 {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

// Set mode knowledge / fuzz
func (r *Runner) setMode(mode string) {
	r.mode = strings.ToLower(mode)
}

// Get mode knowledge / fuzz
func (r *Runner) getMode() string {
	return r.mode
}
