package httpfilter

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	DEFAULT_MODE     = "or"
	OPERATOR_RANGE   = "-"
	OPERATOR_LESS    = "--"
	OPERATOR_GREATER = "++"
	OPERATOR_EQUAL   = "=="
	VALID_MODES      = map[string]struct{}{
		"and": {},
		"or":  {},
	}
)

type Filter struct {
	// Reperesent the amount of different filter types that has been set
	// Note : Not the amount of filter values for each section
	amount int
	// Contains all the configurations
	Config
}

// Response is lightway and adapted to the filter support.
// This makes the filter process faster when a lot of HTTP responses are being analyzed.
type Response struct {
	Body         []byte
	StatusCode   int
	ResponseSize int
	WordCount    int
	LineCount    int
	ResponseTime float64
	Headers      http.Header
}

type Config struct {
	// Mode represent the mode when the filter is running
	// Valid modes: "or", "and", if no mode is set the "or" mode will be used
	Mode                  string
	HeaderRegex           string
	BodyRegex             string
	StatusCodes           []string
	WordCounts            []string
	LineCounts            []string
	ResponseSizes         []string
	ResponseTimesMillisec []string
	Header                http.Header

	// Local configuration for faster lookup and greater preformance
	headerRegex          *regexp.Regexp
	bodyRegex            *regexp.Regexp
	statusCode           map[string][]string
	wordCount            map[string][]string
	lineCount            map[string][]string
	responseSize         map[string][]string
	responseTimeMillisec map[string][]string
	headers              http.Header
}

func NewFilter(config Config) (Filter, error) {
	conf, err := makeConfig(config)
	return Filter{
		Config: conf,
		amount: getFilterAmount(config),
	}, err
}

// Run the filter
// The "typ" represent the filter type: Filter / Match
func (f Filter) Run(resp Response) bool {
	if f.amount == 0 {
		return false
	}

	count := 0
	if f.Config.HeaderRegex != "" && f.HeaderRegex(makeHeaderToBytes(resp.Headers)) {
		count++
	}
	if f.Config.BodyRegex != "" && f.BodyRegex(resp.Body) {
		count++
	}
	if len(f.Config.StatusCodes) > 0 && f.StatusCode(resp.StatusCode) {
		count++
	}
	if len(f.Config.ResponseSizes) > 0 && f.ResponseSize(resp.ResponseSize) {
		count++
	}
	if len(f.Config.WordCounts) > 0 && f.WordCount(resp.WordCount) {
		count++
	}
	if len(f.Config.LineCounts) > 0 && f.LineCount(resp.LineCount) {
		count++
	}
	if len(f.Config.ResponseTimesMillisec) > 0 && f.ResponseTime(resp.ResponseTime) {
		count++
	}
	if len(f.Config.Header) > 0 && f.Header(resp.Headers) {
		count++
	}

	// Check the result in relation to the filter mode
	if f.Mode == "or" && count > 0 {
		return true

	} else if f.Mode == "and" && count == f.amount {
		return true
	}
	return false
}

func (f Filter) IsSet() bool {
	return f.amount > 0
}

// Set the filter mode
func (f Filter) SetMode(mode string) {
	f.Mode = strings.ToLower(mode)
}

func (f Filter) doFilter(valueCompare float64, m map[string][]string) bool {
	mapLength := len(m)

	if mapLength == 0 {
		return false
	}

	countModeAnd := 0
	countValues := 0
	for operator, values := range m {
		countValues += len(values)
		for _, valueStr := range values {
			ok := filter(valueCompare, operator, valueStr)

			// If the mode is true and only a single filter is equal to true, return true on the whole filter process
			if f.Mode == "or" && ok {
				return true

			} else if f.Mode == "and" {
				if !ok {
					return false
				}
				countModeAnd++
			}
		}
	}

	if f.Mode == "and" && countModeAnd == countValues {
		return true
	}
	return false
}

// Child function of: Filter.doFilter()
func filter(valueCompare float64, operator string, value string) bool {
	// Make sure range operator is used first since it require the string value
	// To be splitted and converted to int values
	if operator == OPERATOR_RANGE {
		valueArry, _ := getIntRange(value)
		if valueCompare >= valueArry[0] && valueCompare <= valueArry[1] {
			return true
		}
	} else {
		// The value must be converted to an int and must work (makeConfig has responsibility for this)
		valueFloat64 := mustToFloat64(value)

		if (operator == OPERATOR_EQUAL) && (valueCompare == valueFloat64) {
			return true

		} else if (operator == OPERATOR_GREATER) && (valueCompare > valueFloat64) {
			return true

		} else if (operator == OPERATOR_LESS) && (valueCompare < valueFloat64) {
			return true
		}
	}
	return false
}

func (f Filter) Header(headers http.Header) bool {
	for headerName, _ := range f.headers {
		if _, ok := headers[headerName]; ok {
			return true
		}
	}
	return false
}

func (f Filter) HeaderRegex(headers []byte) bool {
	return len(headers) > 0 && f.headerRegex.Match(headers)
}

func (f Filter) BodyRegex(body []byte) bool {
	return len(body) > 0 && f.bodyRegex.Match(body)
}

func (f Filter) ResponseTime(time float64) bool {
	return f.doFilter(time, f.responseTimeMillisec)
}

func (f Filter) StatusCode(statuscode int) bool {
	return f.doFilter(float64(statuscode), f.statusCode)
}

func (f Filter) ResponseSize(size int) bool {
	return f.doFilter(float64(size), f.responseSize)
}

func (f Filter) WordCount(count int) bool {
	return f.doFilter(float64(count), f.wordCount)
}

func (f Filter) LineCount(count int) bool {
	return f.doFilter(float64(count), f.lineCount)
}

// Set the regex for the header
func (config *Config) SetBodyRegex(regexStr string) error {
	var err error
	if regexStr != "" {
		config.bodyRegex, err = regexp.Compile(regexStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set the regex for the body
func (config *Config) SetHeaderRegex(regexStr string) error {
	var err error
	if regexStr != "" {
		config.headerRegex, err = regexp.Compile(regexStr)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate the configuration of the Filter structure
func makeConfig(config Config) (Config, error) {
	var err error

	config.Mode = strings.ToLower(strings.TrimSpace(config.Mode))

	// Validate the filter mode if it's set to an empty string, set it to its default value
	if config.Mode == "" {
		config.Mode = DEFAULT_MODE

	} else if _, ok := VALID_MODES[config.Mode]; !ok {
		return config, errors.New("invalid filter mode. Valid modes are 'or', 'and'")
	}
	// Validate the regex for the body
	if err = config.SetBodyRegex(config.BodyRegex); err != nil {
		return config, err
	}
	// Validate the regex for the header
	if err = config.SetHeaderRegex(config.HeaderRegex); err != nil {
		return config, err
	}

	// Make operator maps
	if config.statusCode, err = makeOperatorMap(config.StatusCodes); err != nil {
		return config, err
	}
	if config.lineCount, err = makeOperatorMap(config.LineCounts); err != nil {
		return config, err
	}
	if config.wordCount, err = makeOperatorMap(config.WordCounts); err != nil {
		return config, err
	}
	if config.responseSize, err = makeOperatorMap(config.ResponseSizes); err != nil {
		return config, err
	}
	if config.responseTimeMillisec, err = makeOperatorMap(config.ResponseTimesMillisec); err != nil {
		return config, err
	}
	// Data is copied only to keep the organization correct for future called methods (It will keep the same data as the original)
	config.headers = config.Header

	return config, nil
}

// Calculate the amount of filter categories that has been set
func getFilterAmount(config Config) int {
	var count = 0
	if len(config.HeaderRegex) > 0 {
		count++
	}
	if len(config.BodyRegex) > 0 {
		count++
	}
	if len(config.StatusCodes) > 0 {
		count++
	}
	if len(config.WordCounts) > 0 {
		count++
	}
	if len(config.LineCounts) > 0 {
		count++
	}
	if len(config.ResponseSizes) > 0 {
		count++
	}
	if len(config.ResponseTimesMillisec) > 0 {
		count++
	}
	if len(config.Header) > 0 {
		count++
	}
	return count
}

func makeHeaderToBytes(headers http.Header) []byte {
	var headerString strings.Builder
	for key, values := range headers {
		for _, value := range values {
			headerString.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		}
	}
	return []byte(headerString.String())
}

// Take a list of int and convert it into a lookup map.
// In case there are any space characters as prefix and/or suffix, it will be trimmed with 'strings.TrimSpace'.
func makeOperatorMap(lst []string) (map[string][]string, error) {
	var m = make(map[string][]string)

	for _, i := range lst {
		i = strings.TrimSpace(i)
		value, operator, err := getOperator(i)
		if err != nil {
			return m, err
		}

		switch operator {
		case OPERATOR_EQUAL:
			operator = OPERATOR_EQUAL
		case OPERATOR_RANGE:
			_, err := getIntRange(value)
			if err != nil {
				return m, err
			}
		}
		m[operator] = append(m[operator], value)
	}
	return m, nil
}

// Return the operator that was set in the string.
// If an empty value is returned, The operator is invalid or none are set
// !WARNING! : If the operator is the range operator (global variable: OPERATOR_RANGE) the argument string value is returned instead.
func getOperator(s string) (string, string, error) {
	// Verify the given value
	ok_range, _ := regexp.MatchString(`^(-|)(\d+\.|)\d+-(-|)(\d+\.|)\d+$`, s)

	ok_digit, _ := regexp.MatchString(`^(\+\+|)(-{0,3}|)(\d+\.|)\d+$`, s)

	if ok_range && ok_digit || !ok_range && !ok_digit {
		return "", "", fmt.Errorf("httpfilter - invalid operator format given: %s", s)
	}

	if v, ok := strings.CutPrefix(s, OPERATOR_GREATER); ok {
		return v, OPERATOR_GREATER, nil

	} else if v, ok := strings.CutPrefix(s, OPERATOR_LESS); ok {
		return v, OPERATOR_LESS, nil

	} else if ok_range {
		return s, OPERATOR_RANGE, nil

	} else {
		return s, OPERATOR_EQUAL, nil
	}
}

// Get the range from a string value and return the two values within the range as an int array
// If the range of the string value is invalid an error will be triggered
func getIntRange(value string) ([2]float64, error) {
	var (
		errMsg = "invalid range value when trying to get range between two values"
		arry   [2]float64
	)

	lastHyphenIndex := strings.LastIndex(value, "-")
	if lastHyphenIndex == -1 {
		return arry, errors.New("invalid range format, no hyphen (-) found")
	}

	// Handle edge cases like "-100--200", "-100-200", "100--200", "100-200"
	firstPart := value[:lastHyphenIndex]
	secondPart := value[lastHyphenIndex:]

	// Check if the values are negative then a modification is needed
	// Note : Special case when the second part starts with a double hyphen
	if (len(firstPart) > 1 && len(secondPart) > 1) && (strings.HasSuffix(firstPart, "-") && strings.HasPrefix(secondPart, "-")) {
		firstPart, _ = strings.CutSuffix(firstPart, "-")
	} else {
		secondPart, _ = strings.CutPrefix(secondPart, "-")
	}

	v1, err := strconv.ParseFloat(firstPart, 64)
	if err != nil {
		return arry, fmt.Errorf("%s - %s", errMsg, err)
	}
	v2, err := strconv.ParseFloat(secondPart, 64)
	if err != nil {
		return arry, fmt.Errorf("%s - %s", errMsg, err)
	}
	return [2]float64{v1, v2}, nil
}

// Take value of type string and convert it into a int type
func mustToFloat64(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Panic("Invalid filter, can't be converted to float64:", s, err)
	}
	return v
}
