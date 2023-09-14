package functions

import (
	"regexp"
)

// Take a list of regex patterns and check if the string match any. Return true if match otherwise false.
func RegexInLst(pattern []string, str ...string) bool {
	for _, ptn := range pattern {
		for _, s := range str {
			//Check if any pattern match the any of the string then return 'true':
			if ok, _ := regexp.MatchString(ptn, s); ok {
				return true
			}
		}
	}
	return false
}

// Take a regex and check what's character(s) is/are inbetween
func RegexBetween(str, re string) string {
	//Example take hostnames from url: `http?://*(.*?)*/`
	var s string

	r := regexp.MustCompile(re)
	mS := r.FindAllStringSubmatch(str, -1)
	for _, v := range mS {
		s = v[1]
	}
	return s
}
