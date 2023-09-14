package extract

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/Brum3ns/firefly/pkg/files"
)

var (
	WILDCARD = "WILDCARD"
)

type Extract struct {
	Properties
	jobAmount int
	sources   map[string]string                              //(Body|Headers)
	fn_check  map[string]func(item, stringToTest string) int //Note : (contains the pattern/regex "check" function)

	Known Result //<- Known Patterns/Regex that have been discovered
}

type Properties struct {
	Threads         int
	PrefixPatterns  []string
	WordlistPattern map[string][]string //Inc: key="A shared prefix that all words in the list have has", value="wordlist itself"
	WordlistRegex   map[string][]string // -/-
}

type Result struct {
	OK        bool
	TotalHits int
	Pattern   source
	Regex     source
}
type source struct {
	Body    map[string]int
	Headers map[string]int
}

type process struct {
	ok        bool
	hits      int
	typ       string
	method    string
	foundItem string
}

type job struct {
	prefix   string
	method   string
	wordlist []string
}

func NewExtract(p Properties) Extract {
	return Extract{
		fn_check: map[string]func(item, stringToTest string) int{
			"pattern": checkPattern,
			"regex":   checkRegex,
		},
		jobAmount:  (len(p.WordlistPattern) + len(p.WordlistRegex)) * 2, // Note : ("2" : Is used because we have two sources to check (body/headers) )
		Properties: p,
	}
}

func (e *Extract) AddJob(body, headers string, known Result) {
	e.Known = known
	e.sources = map[string]string{
		"body":    body,
		"headers": headers,
	}
}

// Extract all patterns that was found within the response body and/or the response headers from the current target response.
// Want: "Threads to use for pattern extracation", sources: "array of 2 {body, headers}", givenData: "Extract structure".
// Return: status, amountTotal, amountType, headersPatterns, bodyPatterns.
func (e Extract) Run() Result {
	result := Result{
		Pattern: source{
			Body:    make(map[string]int),
			Headers: make(map[string]int),
		},
		Regex: source{
			Body:    make(map[string]int),
			Headers: make(map[string]int),
		},
		TotalHits: 0,
	}

	//Check if the sources empty, if so, then end:
	if len(e.sources["body"]) == 0 || len(e.sources["headers"]) == 0 {
		return result
	}

	var (
		wg       sync.WaitGroup
		process  = make(chan process)
		JobQueue = make(chan job)
	)

	for i := 0; i < e.Threads; i++ {
		go e.analyze(&wg, JobQueue, process)
	}

	go func(wg *sync.WaitGroup) {
		var mutex sync.Mutex
		for {
			r := <-process
			if r.ok {
				result.TotalHits += r.hits
				if r.typ == "body" {
					//Check method for body:
					if r.method == "pattern" {
						mutex.Lock()
						result.Pattern.Body[r.foundItem]++
						mutex.Unlock()

					} else if r.method == "regex" {
						mutex.Lock()
						result.Regex.Body[r.foundItem]++
						mutex.Unlock()
					}
				} else if r.typ == "headers" {
					//Check method for headers:
					if r.method == "pattern" {
						mutex.Lock()
						result.Pattern.Headers[r.foundItem]++
						mutex.Unlock()

					} else if r.method == "regex" {
						mutex.Lock()
						result.Regex.Headers[r.foundItem]++
						mutex.Unlock()
					}
				}
				wg.Done()
			}
		}
	}(&wg)

	e.appendJobs(JobQueue)

	//Remove known patterns/Regex from the result:
	//TODO

	//Wait for all processes to end:
	wg.Wait()

	//Provide success status
	if result.TotalHits > 0 {
		result.OK = true
	}

	return result
}

// Give job to the extract process
func (e *Extract) appendJobs(j chan<- job) {
	m := map[string]map[string][]string{
		"pattern": e.WordlistPattern,
		"regex":   e.WordlistRegex,
	}
	for _, mt := range []string{"pattern", "regex"} {
		for pfx, wl := range m[mt] { //Note : (Just extracting the map - key/value)
			j <- job{
				prefix:   pfx,
				method:   mt,
				wordlist: wl,
			}
		}
	}
}

// TODO
// Compare verified known patterns with discovered to only extract the unique regarding to the fuzzed bahvior:
func (e *Extract) uniqueItem(p process) bool {
	return true
}

// Check for pattern inside the response body & headers:
// This function uses a prefix technique that takes advantage of patterns that share the same prefix to provide a faster and more lightweight analyze.
func (e *Extract) analyze(wg *sync.WaitGroup, jobs <-chan job, result chan<- process) {
	for j := range jobs {
		for _, t := range []string{"body", "headers"} {
			stringToTest := e.sources[t]
			if strings.Contains(stringToTest, j.prefix) {

				//A new wordlist is in need of being analyzed. Add the wordlist length to the listener:
				for _, item := range j.wordlist {

					//Calculate how many times the item was within the content source.
					//Then check if the item is a common (known behavior), if not, then send it to the listener:
					if hits := e.fn_check[j.method](item, stringToTest); hits > 0 {
						wg.Add(1)
						result <- process{
							ok:        true,
							typ:       t,
							hits:      hits,
							foundItem: item,
							method:    j.method,
						}
					}
				}
			}
		}
	}
}

// Regex, string - check regex:
func checkRegex(re, s string) int {
	if match, _ := regexp.MatchString(re, s); match {
		return 1
	}
	return 0
}

// Pattern, string - check amount of pattern in string:
func checkPattern(ptn, s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, ptn)
}

// Create a map from words within a list.
// Return a list of all the shared prefix and a map. The key of the map is a prefix of all the words associated within the words added to the map list.
// Note : (The list that only contain prefix is used for better preformance within loops.)
func CreatePrefixMap(lst []string) ([]string, map[string][]string) {
	var (
		m        = make(map[string][]string)
		lst_pfx  []string
		wildcard bool
	)
	for _, i := range lst {
		k := i[:3]
		m[k] = append(m[k], i)
	}
	for k, l := range m {
		if len(l) == 1 {
			wildcard = true
			//Add all alone items into a wildcard key ("WILDCARD"):
			//Note : This words only if the max prefix for the other items are set to 3 in length.
			m[WILDCARD] = append(m[WILDCARD], k)

			//Delete the old list with only one item:
			delete(m, k)
		} else {
			lst_pfx = append(lst_pfx, k)
		}
	}
	//If the map contained wildcard prefix, then add it at the end:
	if wildcard {
		lst_pfx = append(lst_pfx, WILDCARD)
	}

	return lst_pfx, m
}

// Take a folder that have files (wordlists) with a prefix of: "ptn_" (pattern) OR "_re" (regex).
// Return two wordlist : (Patterns|Regex)
func MakeWordlists(folder string) ([]string, []string) {
	lst_files, _ := files.InDir(folder)

	wordlists := make(map[string][]string)
	for _, f := range lst_files {
		fpath := (folder + f)
		typ := ""
		if strings.HasPrefix(f, "ptn_") {
			typ = "ptn"
		} else if strings.HasPrefix(f, "re_") {
			typ = "re"
		} else { //Simply ignore the other files that miss the prefix
			continue
		}
		//Read file and append each item to the map "wordlists":
		fcontent, _ := os.Open(fpath)
		scanner := bufio.NewScanner(fcontent)
		for scanner.Scan() {
			item := scanner.Text()
			if len(item) > 0 {
				wordlists[typ] = append(wordlists[typ], item)
			}
		}
		fcontent.Close()
	}

	return wordlists["ptn"], wordlists["re"]
}
