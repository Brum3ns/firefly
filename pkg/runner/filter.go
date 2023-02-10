package runner

import (
	"net/http"
	"regexp"
	"strconv"

	fc "github.com/Brum3ns/FireFly/pkg/functions"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

/*
type FilterData struct {
	m_filter map[string][]string
	m_match  map[string][]string
} */
type FResult struct {
	m map[string]bool
	t string
}

type FData struct {
	typ string
	lst []string
}

/**[Filter] - *Match* data within result given by user input*/
func FilterCheck(rt storage.Response, jobs <-chan FData, c chan<- FResult) {
	m := make(map[string]bool)
	for j := range jobs {
		switch j.typ {
		case "sc":
			m[j.typ] = FCode(j.lst, rt.Status)
		case "bs":
			m[j.typ] = FEqualSize(j.lst, rt.BodyByteSize)
		case "wc":
			m[j.typ] = FEqualSize(j.lst, fc.WordCount(string(rt.Body[:])))
		case "lc":
			m[j.typ] = FEqualSize(j.lst, fc.LineCount(string(rt.Body[:])))
		case "re":
			m[j.typ] = FRegex(j.lst[0], string(rt.Body[:]))
		}

		c <- FResult{
			m: m,
			t: j.typ,
		}
	}
}

/**[Filter] - *Filter* data within result given by user input*/

/**Filter - [M/F] response status code*/
func FCode(lsc []string, sc int) bool {
	return fc.InLst(lsc, strconv.Itoa(sc))
}

/** Filter check if given regex is inside body*/
func FRegex(re string, b string) bool {
	ok, err := regexp.MatchString(re, b)
	fc.IFError("p", err)
	return ok
}

/** Filter body size (bs), wordCount (wc), lineCount (lc) if equal to response*/
func FEqualSize(lsz []string, sz int) bool {
	for _, i := range lsz {
		bsize, _ := strconv.Atoi(i)
		if bsize == sz {
			return true
		}
	}
	return false
}

//[TODO]
func FHeader(header *http.Header, str string) bool {
	/** Filter response headers
	 */
	return false
}
