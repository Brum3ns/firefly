package functions

import (
	"regexp"
	"strings"
)

// Extract the 'parameter typ' (get, post, cookie) and the chars to use within each
func GetParam(s string) map[string]string {
	m := make(map[string]string)
	m["post"] = "&"
	m["get"] = "?&"
	m["cookie"] = ";"

	//Check for custom param chars to replace with 'x' type:
	for _, i := range strings.Split(s, " ") {
		v := RegexBetween(i, `\[(.*?)\]`)
		t := strings.ToLower(RegexBetween(i, `^(.*?)\[`))

		if len(v) <= 0 || len(t) <= 0 {
			m["post"] = i
			m["get"] = i
			m["cookie"] = i
			break
		}

		m[t] = v
	}

	return m
}

func FuzzParam(s, typ, params, insert, au string) (string, error) {
	var (
		//inStr = 0
		sp   string
		l    []string
		m    = GetParam(params)
		chrs = m[typ]
	)

	if i, ok := m["all"]; ok {
		chrs = i

	}

	//Add '[' and ']' to solve regex syntax:
	if string(chrs[0]) != "[" && string(chrs[len(chrs)-1]) != "]" {
		chrs = ("[" + chrs + "]")
	}

	switch typ {
	case "get":
		u := s[strings.LastIndex(s, "?")+1:]
		l = regexp.MustCompile(chrs).Split(u, -1)

	case "post":
		l = regexp.MustCompile(chrs).Split(s, -1)
	case "cookie":
		//[TODO]
	}

	//Replace input 'au' to one value to easier use switch case:
	if au == "r" {
		au = "replace"
	} else if au == "a" {
		au = "append"
	}

	//Check input technique:
	switch au {
	case "replace":
		var param string
		sp = s

		if len(l) > 1 {
			for _, p := range l {

				v := p[strings.Index(p, "=")+1:]
				if v == "" {
					continue
				}
				param = strings.Replace(p, v, insert, 1)
				sp = strings.ReplaceAll(sp, p, param)
			}
		}

	case "append":
		sp = s
		for _, p := range l {
			param := p + insert
			sp = strings.ReplaceAll(sp, p, param)
		}
	}

	return sp, nil
}

/*
func DetectParam(url, Postdata string) ([]string, []string) {
	var (
		lst_paramsURL, lst_paramsData []string

		re_url, _  = regexp.MatchString(`(\?.*=|&.*=|;.*=)`, url)
		re_data, _ = regexp.MatchString(`(^.*=|&.*=)`, Postdata)
	)

	//Make sure either the URL or the post data contains atleast one param:
	if !re_url && !re_data {
		return lst_paramsURL, lst_paramsData
	}

	//Check URL for params:
	if re_url {

		//Clear the URL to only have the params left & add them to a list:
		u := url[strings.LastIndex(url, "?")+1:]
		lst_paramsURL = regexp.MustCompile("[?&;]").Split(u, -1)
	}

	//Check post data for params:
	if re_data {
		//-:-
		lst_paramsData = regexp.MustCompile("[&;]").Split(Postdata, -1)
	}

	return lst_paramsURL, lst_paramsData
}
*/
