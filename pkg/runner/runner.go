package runner

import (
	"net/http"
	"sync"

	"github.com/Brum3ns/firefly/pkg/design"
	fc "github.com/Brum3ns/firefly/pkg/functions"
	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
	"github.com/Brum3ns/firefly/pkg/output"
	"github.com/Brum3ns/firefly/pkg/parse"
	"github.com/Brum3ns/firefly/pkg/storage"
)

//Combine `struct(s)` from other pkg sources into struct `Runner`
type Runner struct {
	mutex   sync.Mutex
	options *parse.Options
	client  *http.Client
	//tmp_verifyData *storage.Temp_VerifyData
	verifyData *storage.VerifyData
	verifier   *Verifier
	resp       *storage.Response
	wordlist   *storage.Wordlists
	collection *storage.Collection //[TODO] delete?
	payloads   *storage.Payloads   //[TODO] delete?
	//data     *storage.Analyze
	/* analyze  *storage.Analyze
	errors   *storage.Errors
	Patterns *storage.Patterns */
	verify bool
	count  int
}

type Verifier struct {
	Data map[int]struct {
		Tag           string
		URL           string
		Payload       string
		Body          string
		HeadersString string
		Headers       http.Header
	}
}

var VResp = &storage.VerifyData{}

/**FireFly verify/fuzz runner
* The runner spin up the core engine and execute all tasks including their own engine */
func New(opt *parse.Options, wl *storage.Wordlists, V bool) (*Runner, error) {
	runner := &Runner{
		options:  opt,
		wordlist: wl,
		client:   Client(opt),
		verifier: SetupVerifier(),
		resp:     storage.ConfResp(),
		//data:     &storage.Analyze{},
		verify: V,
		count:  0,
	}

	if runner.verify {
		runner.verifyData = storage.ConfVerifyData()
	} else {
		runner.verifyData = VResp
	}

	var (
		err     error
		reqFail = 0
	)

	if G.Total, err = runner.SelectTotal(runner.verify); err != nil {
		fc.IFError("p", err)
	}

	//Setup channels & Close all [chan]nels when everything is done (defer):
	c_jobs := make(chan storage.Target, len(opt.Target["urls"]))
	c_results := make(chan storage.Response, G.Total)
	c_collection := make(chan storage.Collection, G.Total)
	defer close(c_jobs)
	defer close(c_results)
	defer close(c_collection)

	//[Threads|Requests] - Spinning up threads & send the jobs to the spawned threads:
	for t := 1; t <= opt.Threads; t++ {
		go Request(t, runner.client, opt, c_jobs, c_results)
	}

	//Give jobs to the request task (verify/fuzz):
	if V {
		JobsVerification(runner, c_jobs)
	} else {
		Jobs(runner, c_jobs)
	}

	//[Engine] start up the engine to run all tasks in the background:
	go func() {
		for ; runner.count < (G.Total); runner.count++ {
			result := <-c_results
			//If the response 'status code' was "OK" without any error(s), then procce to engine:
			if result.Status >= 1 {

				//[Engine]a Start the engine for all internal tasks:
				runner.Engine(result.Id, result, c_collection) //[TODO] <= ('result\.R\.Id' is not nessacary u got 'result')

				//If the response failed/timed-out return the status and reason of failure + count request fails:
			} else {
				reqFail++
				c_collection <- storage.Collection{
					Status: false,
					ErrMsg: result.ErrMsg,
				}
			}
		}
	}()

	//[Listener] Intercept tasks that are done and prepare verbose & output:
	for id := 0; id <= (G.Total - 1); id++ {
		count := (id + 1)
		result := <-c_collection

		//If 50% verification requests failed, then exit: (not able to fuzz without knowing the target default behavior)
		if runner.verify {
			if reqFail >= (G.Total / 2) {
				fc.IFFail("vreq")
			}
			design.Loadingbar_Verify(id, G.Total)

			runner.mutex.Lock()
			runner.verifier.Data[result.ID] = struct {
				Tag           string
				URL           string
				Payload       string
				Body          string
				HeadersString string
				Headers       http.Header
			}{
				Tag:           result.Tag,
				URL:           result.UrlNoPayload,
				Payload:       result.Payload,
				Body:          string(result.Body[:]),
				HeadersString: fc.HeadersToStr(result.Headers),
				Headers:       result.Headers,
			}
			runner.mutex.Unlock()

		} else {
			design.DisplayInfo(count, result)
			/* if runner.options.ShowDiff {
				//displayDiff := fmt.Sprintln(design.BlueLight, strings.Join(result.RespDiff[result.UrlNoPayload][result.Payload], "\n"), design.White)
				fmt.Println(displayDiff)
			} */

			//Output result for each result:
			if len(G.OutputFile) > 0 {
				err := output.Output(result)
				fc.IFError("p", err)
			}
		}
	}

	//Dynamic content extraction: (Don't run in a [go]rutine process)
	if runner.verify {
		VResp = runner.Verify()
	}

	return runner, nil
}

/**Setup the Verifier struct for the verification process*/
func SetupVerifier() *Verifier {
	verifier := &Verifier{}
	verifier.Data = make(map[int]struct {
		Tag           string
		URL           string
		Payload       string
		Body          string
		HeadersString string
		Headers       http.Header
	})
	return verifier
}
