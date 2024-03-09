package parameter

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// Global parameter variables:
var (
	SUPPORTED_PARAM_POSITIONS = []string{"url", "body", "cookie"}
	SUPPORTED_PARAM_METHODS   = []string{"replace", "append"}
)

type Parameter struct {
	// InsertKeyword is a shortcut to access the global insert keyword used for all params in all positions
	// Keyword presented for example in: URL, Body, Cookie
	InsertKeyword   string
	AutoQueryURL    bool
	AutoQueryBody   bool
	AutoQueryCookie bool
	URL             query
	Body            query
	Cookie          query
}

type query struct {
	// Auto detect and insert InsertKeyword in the RawQuery located in "query"
	// Holds the original raw query
	RawQueryOriginal string
	// Holds the built raw query and adapted raw query based on the given rules
	RawQueryInsertPoint string
	// Holds all the params and it's settings / properties, including it's values
	Params []paramSettings
	// Position holds the name of the position the request parameter was placed.
	// Example: (url, body, cookie).
	Position string
	// The rules to be used for the param and it's settings
	Rule QueryRules
}

type QueryRules struct {
	// Separators is optional and if no separators are set, the default separators in relation to the method will be used instead.
	Separators []rune
	// The method that the param will used for the "InsertKeyword" variable value when building adapted raw queries
	// Ex: (replace, append)
	Method string
	// The insert keyword that will used to method to be modified after the users request
	InsertKeyword string
}

// Holds the settings of a parameter
type paramSettings struct {
	// Name contains the full name of the parameter (not value).
	// In case the parameter name contains a name such as: "query[0]=1" the full name will be "query[0]".
	Name      string
	Separator string
	//Holds all the parameter including values as a raw string
	Raw string
	// The values holds all the values and in case the param is used multiple times then all values are added to the list:
	// When a parameter do have an equal sign but there is no value that representate that the param was defined with an empty stirng.
	// However when the value is "nil" the param did not have any equal sign and represent the param to only be declared and not defined.
	Value any
	// Index representate the order the given parameter was place in
	Index int
	// Type representate the value types that the param holds.
	// Param types can be string, int, bool, array, list, object etc...
	// Note : (By default if no parameter type is set, the default will be "string")
	Type reflect.Kind
	// Comment ...
	TypeValue reflect.Kind
}

// Make a param object that contain the param settings and rules for each position.
// The argument "paramRules" holds the supported *position* ("url", "cookie", "body") and the rules for the parameter placed in that position.
func NewParameter(rule map[string]QueryRules, insertKeyword string) (Parameter, error) {
	parameter := Parameter{}

	// Verify before returning the parameter [struct]ure:
	// !Note : (Rules should already be validated in the [struct]ure "NewParamRules")
	for position, rules := range rule {
		if r, err := NewRules(insertKeyword, rules.Method, rules.Separators); err == nil {
			q := query{
				Params: make([]paramSettings, 0),
				Rule:   r,
			}
			switch position {
			case "url":
				parameter.URL = q
				parameter.AutoQueryURL = true

			case "body":
				parameter.Body = q
				parameter.AutoQueryBody = true

			case "cookie":
				parameter.Cookie = q
				parameter.AutoQueryCookie = true

			default:
				return Parameter{}, errors.New("invalid method used, only support: " + strings.Join(SUPPORTED_PARAM_POSITIONS, ","))
			}
		} else {
			return parameter, err
		}
	}
	return parameter, nil
}

func NewRules(insertKeyword, method string, separators []rune) (QueryRules, error) {
	rules := QueryRules{
		Method:        method,
		Separators:    separators,
		InsertKeyword: insertKeyword,
	}
	//Verify that the given method is within the supported list of methods:
	if !slices.Contains(SUPPORTED_PARAM_METHODS, rules.Method) {
		return rules, errors.New("invalid method used, only support: " + strings.Join(SUPPORTED_PARAM_METHODS, ","))
	}
	return rules, nil
}

func (p *Parameter) GetURLParam(param string) []paramSettings {
	return getParam(param, p.URL.Params)
}

func (p *Parameter) GetBodyParam(param string) []paramSettings {
	return getParam(param, p.Body.Params)
}

func (p *Parameter) GetCookieParam(param string) []paramSettings {
	return getParam(param, p.Cookie.Params)
}

// Set the URL parameter from a rawQuery
func (p *Parameter) SetURLparams(rawQuery string) {
	p.URL.Position = "url"
	p.URL.RawQueryOriginal = rawQuery
	p.URL.Params = setParam(p.URL)
	p.URL.RawQueryInsertPoint = buildRawQueryInsertPoint(p.URL)
}

// Set the body parameter from a rawQuery
func (p *Parameter) SetBodyparams(rawQuery string) {
	p.Body.Position = "body"
	p.Body.RawQueryOriginal = rawQuery
	p.Body.Params = setParam(p.Body)
	p.Body.RawQueryInsertPoint = buildRawQueryInsertPoint(p.Body)
}

// Set the cookie parameter from a rawQuery
func (p *Parameter) SetCookieparams(rawQuery string) {
	p.Cookie.Position = "cookie"
	p.Cookie.RawQueryOriginal = rawQuery
	p.Cookie.Params = setParam(p.Cookie)
	p.Cookie.RawQueryInsertPoint = buildRawQueryInsertPoint(p.Cookie)
}

func getParam(param string, params []paramSettings) []paramSettings {
	var l []paramSettings
	for _, p := range params {
		//Add the parameter that matched the given, continue to check in case multiple is included:
		if p.Name == param {
			l = append(l, p)
		}
	}
	return l
}

// set param settings from the given rawQuery.
func setParam(q query) []paramSettings {
	q.Position = strings.ToLower(q.Position)
	params := []paramSettings{}

	if len(q.Rule.Separators) == 0 {
		q.Rule.Separators, _ = getSeparators(q.Position)
	}
	//Loop over all the discovered parameters:
	index := 0

	rawQuerySplit(q.Position, q.RawQueryOriginal, q.Rule.Separators)

	for _, lst := range rawQuerySplit(q.Position, q.RawQueryOriginal, q.Rule.Separators) {
		var (
			sep      = lst[0]
			rawParam = lst[1]
			name     string
			value    any //Note : (List in case the param is being found multiple times with different values)
		)
		//Split the first discovered equal sign (=) and define the parameter name:
		l := strings.SplitN(rawParam, "=", 2)
		name = l[0]

		//In case the list has the length of two. A value was presented in the parameter, then add it:
		//Note: (In case no equal sign (=) is presented within the parameter: v = nil)
		if len(l) == 2 {
			value = l[1]
		}
		//Add the parameter to the final map containing all parameters and values:
		param := paramSettings{
			Name:      name,
			Separator: sep,
			Value:     value,
			Raw:       rawParam,
			Index:     index,
		}
		param.setType()

		params = append(params, param)
		index++
	}
	return params
}

// Set the parameter type:
func (p *paramSettings) setType() {
	isArray := func(s string) bool {
		switch {
		case strings.Contains(s, "[") && strings.Contains(s, "]"):
			return true
		case strings.Contains(s, "%5B") && strings.Contains(s, "%5D"):
			return true
		default:
			return false
		}
	}

	switch p.Value.(type) {
	case string:
		v := p.Value.(string)

		if isArray(p.Name) { // List
			p.Type = reflect.Array

		} else if ok, err := strconv.ParseBool(v); ok && err == nil { //[Bool]ean
			p.Type = reflect.Bool

		} else if _, err = strconv.Atoi(v); err == nil { //[Int]eger
			p.Type = reflect.Int

		} else if _, err = strconv.ParseFloat(v, 64); err == nil { //Float (32/64)
			p.Type = reflect.Float64

		} else { //String
			p.Type = reflect.String
		}
	case nil: //Param value not defined = invalid / nil
		p.Type = reflect.Invalid
	}
}

// Build a raw query adapted to the given insertpoint keyword(s):
// !Note : (The function DO NOT verify the method in the "query" [struct]ure)
func buildRawQueryInsertPoint(q query) string {
	var rawQuery string

	// Loop over all params presented in the given query:
	for _, p := range q.Params {
		rawParam := (p.Separator + p.Name)

		// If the original value in the rawQuery was nil without an equal sign.
		// Then add directly and continue:
		if p.Value == nil {
			rawQuery += rawParam
			continue
		}

		// Check which method should be used to build the raw query together with the insertpoint keyword:
		var v string
		if q.Rule.Method == "replace" {
			v = q.Rule.InsertKeyword

		} else if q.Rule.Method == "append" {
			v = fmt.Sprintf("%v%s", p.Value, q.Rule.InsertKeyword)
		}

		rawQuery += (rawParam + "=" + v)
	}
	return rawQuery
}

// Get the default param separators based on the given param point (placed at)
func getSeparators(point string) ([]rune, error) {
	var defaultSeparators = map[string][]rune{
		"url":    {'&'},
		"body":   {'&'},
		"cookie": {';'},
	}
	if separators, ok := defaultSeparators[point]; ok {
		return separators, nil
	} else {
		return []rune{}, errors.New("invalid param point. The valid are: ")
	}
}

// Modified alias of "strings.FieldsFunc":
func rawQuerySplit(position, rawQuery string, separators []rune) [][2]string { // <----- BODY AND COOKIE DO NOT INCLUDE FIRST SEPARATOR
	var arryResult [][2]string

	f := func(r rune) bool {
		for _, i := range separators {
			if i == r {
				return true
			}
		}
		return false
	}

	// Add a value separator character at the end to make it possible to include all the raw parameters (name + value) and it's related separator.
	// The added separator will not be included and is only added to easily detect the last param properties at the end of the given raw query.
	if len(separators) > 0 {
		rawQuery = rawQuery + string(separators[0])
	}

	start := -1
	for end, r := range rawQuery {
		if f(r) {
			if start >= 0 {
				sep := ""
				if start > 0 {
					sep = string(rawQuery[start-1])
				}
				arryResult = append(arryResult, [2]string{sep, rawQuery[start:end]})
				// Set start to a negative value.
				// Note: using -1 here consistently and reproducibly
				// slows down this code by a several percent on amd64.
				start = ^start
			}
		} else {
			if start < 0 {
				start = end
			}
		}
	}
	return arryResult
}
