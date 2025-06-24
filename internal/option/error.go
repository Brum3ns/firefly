package option

import (
	"errors"
	"fmt"
)

// Failed messages:
var ERRORCODES = map[string]error{
	"OUT-FILE-001":      errors.New("invalid output file"),
	"OUT-LOG-001":       errors.New("invalid log output file"),
	"IN-REQ-001":        errors.New("invalid raw HTTP request"),
	"IN-URL-001":        errors.New("invalid target URL"),
	"PAY-WORDLIST-001":  errors.New("invalid payload wordlist"),
	"PAY-WORDLIST-002":  errors.New("invalid payload verify wordlist"),
	"HTTP-METHOD-001":   errors.New("invalid HTTP method "),
	"HTTP-HEADER-001":   errors.New("invalid HTTP header "),
	"HTTP-TIMEOUT-001":  errors.New("invalid HTTP timeout "),
	"HTTP-PATH-001":     errors.New("invalid HTTP path "),
	"HTTP-BODY-001":     errors.New("invalid HTTP body "),
	"HTTP-REDIRECT-001": errors.New("invalid redirect option"),
	"HTTP-PROXY-001":    errors.New("invalid HTTP proxy"),
	"HTTP-PROXY-002":    errors.New("invalid socks proxy"),
}

func getErrorByCode(errorcode string) error {
	if err, ok := ERRORCODES[errorcode]; ok {
		return err
	}
	return fmt.Errorf("could not find error message with given errorcode:[%s]", errorcode)
}
