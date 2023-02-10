package functions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

/**Take any string and encode to valid utf-8/ASCII. Return new string*/
func RuneToASCII(str string) string {
	var (
		s string
		c string
	)
	for _, r := range str {
		if r != utf8.RuneError {
			c = string(r)
		} else {
			c = strconv.QuoteRuneToASCII(r)
			c = c[1 : len(c)-1] //Remove quotes
		}
		s += c
	}
	return s
}

func MapToString(m map[string]string, sep string) string {
	var str string

	for key, value := range m {
		str += fmt.Sprintf("%s%s%s", key, sep, value)
	}

	return str
}

func RandomCase(str string) string {
	var (
		randInt   = time.Now().UnixNano() % 2
		wildcard  = 0
		wildforce = 0
	)

	if (wildcard < 3 || wildforce > 3) && randInt == 0 {
		wildforce = 0
		wildcard++

		str = strings.ToUpper(str)

	} else {
		wildcard = 0
		wildforce++
	}

	return str
}

func CheckSubstrings(str string, subs ...string) int {
	m := 0
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			m += 1
		} else {
			m = 0
		}
	}

	return m
}

func RegexBetween(str, re string) string {
	/** Take a regex and check what's character(s) is/are inbetween
	 */

	var s string

	//Example take hostnames from url > `http?://*(.*?)*/`
	r := regexp.MustCompile(re)
	mS := r.FindAllStringSubmatch(str, -1)
	for _, v := range mS {
		s = v[1]
	}

	return s
}

/** Take a regex and check what's character(s) is/are inbetween*/
func RegexBetweenLst(str, re string) []string {
	var l []string

	//Example take hostnames from url > `http?://*(.*?)*/`
	r := regexp.MustCompile(re)
	mS := r.FindAllStringSubmatch(str, -1)

	for _, v := range mS {
		if !InLstAll(l, v[1]) {
			l = append(l, v[1])
		}
	}

	return l
}

func Err_Highlight(lE, l []string) string {
	var s string
	s = fmt.Sprintf("%v", l)
	s = s[1 : len(s)-1]

	for _, e := range lE {
		s = strings.ReplaceAll(s, e, fmt.Sprintf("\033[31m%s\033[0m", e))
	}

	return s
}

func LineIn(lb []byte) int { //OLD
	return len(strings.Split(string(lb[:]), "\n"))
}

func WordIn(lb []byte) int { //OLD
	return len(strings.Fields(string(lb[:])))
}

/** Count lines in a string*/
func WordCount(s string) int {
	return len(strings.Fields(s))
}

/** Count words in a string*/
func LineCount(s string) int {
	return len(strings.Split(s, "\n"))
}
