package extract

import (
	"bufio"
	"log"
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

	//Known Result //<- Known Patterns/Regex that have been discovered
}

type Properties struct {
	Threads         int
	PrefixPatterns  []string
	WordlistPattern map[string][]string //Inc: key="A shared prefix that all words in the list have has", value="wordlist itself"
	WordlistRegex   map[string][]string // -/-
}

type Result struct {
	OK             bool
	TotalHits      int
	PatternBody    map[string]int
	PatternHeaders map[string]int
	RegexBody      map[string]int
	RegexHeaders   map[string]int
}

// !Note : (MUST be the same name as the "Result")
type ResultCombine struct {
	PatternBody    map[string][]int `json:"PatternBody"`
	PatternHeaders map[string][]int `json:"PatternHeaders"`
	RegexBody      map[string][]int `json:"RegexBody"`
	RegexHeaders   map[string][]int `json:"RegexHeaders"`
}

type process struct {
	done      bool
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

func NewCombine() ResultCombine {
	return ResultCombine{
		PatternBody:    make(map[string][]int),
		PatternHeaders: make(map[string][]int),
		RegexBody:      make(map[string][]int),
		RegexHeaders:   make(map[string][]int),
	}
}

func (e *Extract) AddJob(body, headers string) {
	//e.Known = known
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
		TotalHits:      0,
		PatternBody:    make(map[string]int),
		PatternHeaders: make(map[string]int),
		RegexBody:      make(map[string]int),
		RegexHeaders:   make(map[string]int),
	}

	//Check if the sources empty, if so, then end:
	if len(e.sources["body"]) == 0 || len(e.sources["headers"]) == 0 {
		return result
	}

	var (
		wg             sync.WaitGroup
		processChannel = make(chan process)
		JobQueue       = make(chan job)
	)

	for i := 0; i < e.Threads; i++ {
		go e.analyze(&wg, JobQueue, processChannel)
	}

	go func(wg *sync.WaitGroup) {
		var mutex sync.Mutex
		for {
			r := <-processChannel

			if r.ok {
				result.TotalHits += r.hits
				if r.typ == "body" {
					//Check method for body:
					if r.method == "pattern" {
						mutex.Lock()
						result.PatternBody[r.foundItem] = r.hits
						mutex.Unlock()

					} else if r.method == "regex" {
						mutex.Lock()
						result.RegexBody[r.foundItem] = r.hits
						mutex.Unlock()
					}
				} else if r.typ == "headers" {
					//Check method for headers:
					if r.method == "pattern" {
						mutex.Lock()
						result.PatternHeaders[r.foundItem] = r.hits
						mutex.Unlock()

					} else if r.method == "regex" {
						mutex.Lock()
						result.RegexHeaders[r.foundItem] = r.hits
						mutex.Unlock()
					}
				}

			} else if r.done {
				wg.Done()
			}
		}
	}(&wg)

	e.appendJobs(&wg, JobQueue)

	//Wait for all processes to end:
	wg.Wait()

	//Provide success status
	if result.TotalHits > 0 {
		result.OK = true
	}

	return result
}

// Give job to the extract process
func (e *Extract) appendJobs(wg *sync.WaitGroup, j chan<- job) {
	m := map[string]map[string][]string{
		"pattern": e.WordlistPattern,
		"regex":   e.WordlistRegex,
	}
	for _, mt := range []string{"pattern", "regex"} {
		for pfx, wl := range m[mt] { //Note : (Just extracting the map - key/value)
			wg.Add(1)
			j <- job{
				prefix:   pfx,
				method:   mt,
				wordlist: wl,
			}
		}
	}
}

// Take the current extracted map result and compare it with a known map result.
// Return the "current" map and all the unique items with their unique values
func GetUnique(current map[string]int, known map[string][]int, payload string) (map[string]int, int) {
	hit := 0
	for item, ValueCurrent := range current {
		//If the key exists and they share the same value, delete the key from the "current" map:
		uniuqe := true
		if lstValue, ok := known[item]; ok {
			if len(lstValue) == 1 && ValueCurrent == lstValue[0] {
				uniuqe = false

			} else {
				for _, value := range lstValue {
					if ValueCurrent == value {
						uniuqe = false
						break
					}
				}
			}
		}
		if uniuqe && !strings.Contains(payload, item) {
			hit += ValueCurrent
		}
		delete(current, item)
	}
	return current, hit
}

// Take a list that contains an array of two maps. The first map in the array is the *current* map and the secound map is a map that contains known items.
// Compare the two maps in the array for all arrays in the list and delete all known items that was detected inside the *current* map.
// Return a list of maps in the same order as given that only contains the unique items of the *current map*.
// !Note : (The order for input is important since it effects the return order of the final list of maps)
func GetMultiUnique(current []map[string]int, known []map[string][]int, payload string) ([]map[string]int, int) {
	if len(current) != len(known) {
		log.Fatal("length was different from \"current\" and \"known\" map list given in extract.")
	}

	//Extract the full list and take the two maps in each list to compare the differences within them:
	storageDiff := []map[string]int{} //<-List to store the differences (same order as it was given in)
	totalHits := 0

	for i := 0; i < len(current); i++ {
		uniqueItems, hit := GetUnique(current[i], known[i], payload)

		storageDiff = append(storageDiff, uniqueItems)
		totalHits += hit
	}

	return storageDiff, totalHits
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
		result <- process{done: true}
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
