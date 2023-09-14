package functions

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// Take a list of []int and return the sum
func LstSum(l []int) int {
	sum := 0
	for i := 0; i < len(l); i++ {
		sum = sum + l[i]
	}
	return sum
}

// Take file input and return a list of it's collected data
func FileToList(filePath string) []string {
	var (
		lst     []string
		file, _ = os.Open(filePath)
		scanner = bufio.NewScanner(file)
	)
	defer file.Close()
	for scanner.Scan() {
		if scanner.Text() != "" {
			lst = AppendUniqueString(lst, scanner.Text())
		}
	}
	return lst
}

// Append a string to a list Works similar as append but do not append duplicates or empty strings
func AppendUniqueString(l []string, s string) []string {
	if !InLst(l, s) || s != "" {
		l = append(l, s)
	}
	return l
}

// Append a string to a list Works similar as append but do not append duplicates or empty strings
func AppendUniqueInt(l []int, i int) []int {
	for _, item := range l {
		if item == i {
			return l
		}
	}
	return append(l, i)
}

// String to list, take a string that contains a "char",splitt it and delete the duplicated values. Return the new list without duplicates
func ToLstSplit(str, char string) []string {
	var lst []string

	for _, item := range strings.Split(str, char) {
		if !InLst(lst, item) {
			lst = append(lst, item)
		}
	}
	return lst
}

// String to list. Use a regex with chars to splitt input string into a list
func ToLstSplitRe(s string, re string) []string { //<--- [TODO] - (Replace with ToLstSplit)
	var (
		start = 0
		reg   = regexp.MustCompile(re)
		idx   = reg.FindAllStringIndex(s, -1)
		l     = make([]string, len(idx)+1)
	)

	for i, e := range idx {
		l[i] = s[start:e[0]]
		start = e[1]
	}
	l[len(idx)] = s[start:] // <> l[len(idx)] = s[start:len(s)] //DELETE?

	return l
}

// Check if an item is wihin a list, (allow empty "no values")
func InLstAll(l []string, s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

// Check if an item is wihin a list (Not empty values)
func InLst(l []string, s string) bool {
	for _, i := range l {
		if i != "" && i == s {
			return true
		}
	}
	return false
}
