package functions

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Float(itf interface{}) (float64, error) {
	var floatType = reflect.TypeOf(float64(0))

	v := reflect.Indirect(reflect.ValueOf(itf))
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

// String to int (ignore errors)
func Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Int to string (ignore errors)
func Str(i int) string {
	s := strconv.Itoa(i)
	return s
}

// Remove duplicates from a map
func RmDup_m(m map[string][]string) []string {
	var (
		mp = make(map[string]int)
		l  []string
	)
	for _, lst := range m {
		for _, i := range lst {
			if mp[i] == 0 {
				l = append(l, i)
			}
		}
	}
	return l
}

// Take any string and encode to valid utf-8/ASCII. Return new string
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

func Err_Highlight(lE, l []string) string {
	var s string
	s = fmt.Sprintf("%v", l)
	s = s[1 : len(s)-1]

	for _, e := range lE {
		s = strings.ReplaceAll(s, e, fmt.Sprintf("\033[31m%s\033[0m", e))
	}

	return s
}

// Count words in a string
func WordCount(s string) int {
	return len(strings.Fields(s))
}

// Count lines in a string
func LineCount(s string) int {
	return len(strings.Split(s, "\n"))
}
