package functions

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	G "github.com/Brum3ns/FireFly/pkg/functions/globalVariables"
	"github.com/Brum3ns/FireFly/pkg/storage"
)

/**Take a list of []int and return the sum*/
func LstSum(l []int) int {
	sum := 0

	for i := 0; i < len(l); i++ {
		sum = sum + l[i]
	}
	return sum
}

func LstAppendSame(a, b []string) []string {
	/**Compare two lists and see if they share atleast one item in common
	* Return true/false if it does or dosen't
	 */
	l := []string{}
	for _, i := range a {
		for _, j := range b {
			if i == j {
				l = append(l, i)
			}
		}
	}
	return l
}

/** Make a ASCII list containing [A-Za-z0-9] Return the lst*/
func LstAscii() []string {
	aci := ""
	for v := 65; v < 98; v += 32 { // 65++[A-Z] - 97++[a-z] -> (2 loops)
		for i := 0; i < 26; i++ {
			aci += string(rune(v + i))
		}
	}
	l := strings.Split(aci, "")
	return l
}

func FileToLst(filePath string) []string {
	/** Insert list input to a file
	 */

	var lst []string

	//Setup static list:
	file, _ := os.Open(filePath)
	scanner := bufio.NewScanner(file)
	lbar := []string{"|", "/", "-", "\\"}
	i := 0
	c := 0

	for scanner.Scan() {

		//Loading bar verbose:
		if i >= len(lbar) {
			i = 0
		}
		if c%2 == 0 {
			print("\r", lbar[i])
			i++
		}

		if !InLst(lst, scanner.Text()) {
			//Check if it's from a resource file. Then count the amount of items "Grep_":
			if strings.Contains(filePath, "Grep_") {
				G.Amount_LstG += 1
			}

			lst = ToLst_s(lst, scanner.Text())
		}
	}
	return lst
}

/**Append none existing item to given list. Return given list with or without the new item added*/
/* func AppendUniq_s(l []string, s string) []string {
	for _, i := range l {
		if i == s {
			return l
		}
	}
	return append(l, s)
}

func LstClear(lst []string, char []string) []string {

	var l []string
	for _, i := range lst {
		for _, c := range char {
			if i == c {
				break
			}
			l = append(l, i)
		}
	}

	return l
}*/

func FileToArray(filePath string, lst []string) []string {
	/** Take a file and insert the items into an array
	 */
	var item string

	file, _ := os.Open(filePath)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		item = scanner.Text()

		if !InLst(lst, item) {
			lst = append(lst, item)
		}
	}

	return lst
}

func Payload_FileToArray(replace, encode, tamper, filePath string) []string {
	var (
		item string
		lst  []string
	)

	G.Amount_Lst++

	file, _ := os.Open(filePath)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if !InLst(lst, scanner.Text()) {
			item = scanner.Text()
			//Calculate the amount of items all wordlists have together (Not dublicatres). Then add them to it's separete list:
			G.Amount_Item++

			//Regex to replace within each payload:
			if replace != "" {
				var lst_replace = make(map[int][]string)

				lst := strings.Split(replace, "\n")
				for id, i := range lst {
					if !strings.Contains(i, " => ") {
						fmt.Println("[\033[31mFA\033[0m] \"payload-replace\" do not contain the \" => \" (space included). FireFly dosen't know what to replace the regex/string with ")
						os.Exit(0)
					}
					l := strings.Split(i, " => ")
					lst_replace[id] = []string{l[0], l[1]}
				}

				for _, i := range lst_replace {
					re := regexp.MustCompile(i[0])
					item = re.ReplaceAllString(item, i[1])
				}
			}

			//Tamper to use for the payloads:
			if tamper != "" {
				item = PayloadTamper(item, tamper)
			}

			//Encode each payload:
			if encode != "" {
				item = PayloadEncode(item, encode)
			}

			lst = append(lst, item)
		}
	}

	return lst
}

func WordlistToUse(AttackM string, wl *storage.Wordlists) (string, []string) {
	/** Select attack list type
	 */

	var (
		lst []string
		tag string
	)

	switch AttackM {
	case "fuzz":
		lst = wl.Fuzz
		tag = "fuzz"
	case "transformation":
		lst = wl.TransformationPayloads
		tag = "transformation"
	}
	return tag, lst
}

func StringToArray_Attack(replace string, encode string, tamper string, str string, wl *storage.Wordlists, sys bool) {
	/** String to an array
	 */

	var tag, wordlist string

	if sys {
		split_str := strings.Split(str, "_")
		split_toTag := strings.Split(split_str[0], "/")
		tag = split_toTag[len(split_toTag)-1]
		wordlist = str

	} else {
		split_str := strings.Split(str, ":")
		tag = split_str[1]
		wordlist = split_str[0]
	}

	//Checking type of wordlist:
	switch tag {
	case "fuzz":
		wl.Fuzz = Payload_FileToArray(replace, encode, tamper, wordlist)

	case "transformation":
		wl.TransformationPayloads = append(wl.TransformationPayloads)

	case "reflect":
		wl.Reflect = Payload_FileToArray(replace, encode, tamper, wordlist)

	case "cp":
		wl.CachePoisoning = Payload_FileToArray(replace, encode, tamper, wordlist)

	case "headers":
		wl.Headers = Payload_FileToArray(replace, encode, tamper, wordlist)

	case "dir":
		wl.Directories = Payload_FileToArray(replace, encode, tamper, wordlist)

	case "time":
		wl.Timebased = Payload_FileToArray(replace, encode, tamper, wordlist)
	}
}

/**Check if two slices are the same*/
/* func SlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
} */

/** Count how many times a string repeat inside a list  -> Return a map that only contained the unique items within the list*/
/* func OnlyUnique(l []string) map[int][]string {
	m := make(map[int][]string)
	l_count := make(map[string]int)
	lst := []string{}

	//Check how many times each item repeat inside the input "lst":
	for _, str := range l {
		l_count[str] += 1
	}
	//Check the items that only repeater one single time in the "lc" (items from lst it's amount):
	for str, amount := range l_count {
		if amount == 1 {
			lst = append(lst, str)
		}
	}
	//Amount of line diff, diff list:
	m[len(lst)] = lst

	return m
} */

func OnlyUnique(l []string) []string {
	m_checkDups := make(map[string]int)
	lst := []string{}

	//Check how many times each item repeat inside the input "lst":
	for _, str := range l {
		m_checkDups[str] += 1
	}
	//Check the items that only repeater one single time in the "lc" (items from lst it's amount):
	for str, amount := range m_checkDups {
		if amount == 1 {
			lst = append(lst, str)
		}
	}
	return lst
}

/** Remove duplicates from a list. Return the same list that was imported with no duplicates*/
func LstRmDups(l []string) []string {
	lAmount := make(map[string]bool)
	lst := []string{}
	for _, i := range l {
		if _, v := lAmount[i]; !v {
			lAmount[i] = true
			lst = append(lst, i)
		}
	}
	return lst
}

/** Append 'x' amount of lists and return a new single list. This function ignore duplicates and nil strings*/
/* func ToLst_l(ll [][]string) []string {
	var lst []string
	for _, l := range ll {
		for _, i := range l {
			if !InLst(lst, i) || i != "" {
				lst = append(lst, i)
			}
		}
	}

	return lst
} */

/** Append a string to a list Works similar as append but do not append duplicates or nil strings */
func ToLst_s(l []string, s string) []string {
	if !InLst(l, s) || s != "" {
		l = append(l, s)
	}

	return l
}

/**Append unique strings to a list that only repeat 'one time' -> Return to same list */
/* func ToLst_u(l []string) []string {
	lc := make(map[string]int)

	//Check how many times each item repeat inside the input "lst":
	for _, i := range l {
		for _, j := range strings.Fields(i) {
			lc[j] += 1
		}
	}

	//Check the items that only repeater one single time in the "lc" (items from lst it's amount):
	for i, nr := range lc {
		if nr == 1 {
			l = append(l, i)
		}
	}
	fmt.Println(">", l) //DEBUG

	return l
} */

/**Take HTTP headers. Return a list with the headers*/
func ToLstHeaders(hs http.Header) []string {
	l := []string{}

	for h, v := range hs {
		fmt.Println(h, v)
	}
	return l
}

/**String to list, take a string that contains a "char",splitt it and delete the duplicated values. Return the new list without duplicates.*/
func ToLstSplit(str, char string) []string {
	var (
		lst     []string
		new_str = strings.Split(str, char)
	)

	for _, item := range new_str {
		if !InLst(lst, item) {
			lst = append(lst, item)
		}
	}

	return lst
}

/**String to list. Use a regex with chars to splitt input string into a list */
func ToLstSplitRe(s string, re string) []string { //<--- [TODO] - (Replace with ToLstSplit)
	reg := regexp.MustCompile(re)
	idx := reg.FindAllStringIndex(s, -1)
	start := 0
	l := make([]string, len(idx)+1)
	for i, e := range idx {
		l[i] = s[start:e[0]]
		start = e[1]
	}
	l[len(idx)] = s[start:len(s)]

	return l
}

/* func RmLstItem(l []string, nr int) []string {
	return append(l[:nr], l[nr+1:]...)
} */

/**Check if an item is wihin a list, (allow empty "no values")*/
func InLstAll(l []string, s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}

/** Check if an item is wihin a list (Not empty values)*/
func InLst(l []string, s string) bool {
	for _, i := range l {
		if i == "" {
			continue
		}
		if i == s {
			return true
		}
	}
	return false
}
