// Package encode provides common encodings for pentesting payload manipulation.
package encode

import (
	"bytes"
	"compress/gzip"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strings"
)

type Encoder func(string) (string, error)

var encodeFuncs = map[string]Encoder{
	"url":       SingleUrlEncode,
	"urlfull":   UrlEscape,
	"urldouble": DoubleUrlEncode,
	"html":      HTMLEscape,
	"htmlhex":   HTMLHex,
	"base64":    Base64,
	"base32":    Base32,
	"hex":       Hex,
	"json":      JSONEncode,
	"binary":    BinaryEncode,
	"gzip":      Gzip,
	"unicode":   UnicodeEscape,
}

// Encode applies the selected encodings in order.
// Supports format: "xor:<key>" and others.
func Encode(payload string, parts []string, encodings []string) (string, error) {
	if len(encodings) == 0 {
		return payload, nil
	}

	var err error
	for _, enc := range encodings {
		encFunc, ok := encodeFuncs[strings.ToLower(enc)]
		if !ok {
			return "", fmt.Errorf("unsupported encoding: %s", enc)
		}

		if len(parts) > 0 {
			for _, part := range parts {
				rawPart := part
				part, err = encFunc(part)
				if err != nil {
					return "", fmt.Errorf("error in %s encoding: %w", enc, err)
				}
				payload = strings.ReplaceAll(payload, rawPart, part)
			}
		} else {
			payload, err = encFunc(payload)
			if err != nil {
				return "", fmt.Errorf("error in %s encoding: %w", enc, err)
			}
		}
	}
	return payload, nil
}

func GetEncoders() []string {
	var encoders []string
	for e := range encodeFuncs {
		encoders = append(encoders, e)
	}
	return encoders
}

// --- Encoding Implementations ---
func SingleUrlEncode(s string) (string, error) {
	return url.QueryEscape(s), nil
}

/* func DoubleUrlEncode(s string) (string, error) {
	encoded := url.QueryEscape(s)
	return url.QueryEscape(encoded), nil
} */

func UrlEscape(s string) (string, error) {
	var encoded strings.Builder
	for _, r := range s {
		encoded.WriteString(fmt.Sprintf("%%%02X", r))
	}
	return encoded.String(), nil
}

func DoubleUrlEncode(s string) (string, error) {
	escaped := url.QueryEscape(s)
	return strings.ReplaceAll(escaped, "%", "%25"), nil
}

func HTMLEscape(s string) (string, error) {
	return html.EscapeString(s), nil
}

func HTMLHex(s string) (string, error) {
	var out strings.Builder
	for _, r := range s {
		out.WriteString(fmt.Sprintf("&#x%X;", r))
	}
	return out.String(), nil
}

func Base32(s string) (string, error) {
	return base32.StdEncoding.EncodeToString([]byte(s)), nil
}

func Base64(s string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(s)), nil
}

func Hex(s string) (string, error) {
	return hex.EncodeToString([]byte(s)), nil
}

func JSONEncode(s string) (string, error) {
	j, err := json.Marshal(s)
	return string(j), err
}

func BinaryEncode(s string) (string, error) {
	var b strings.Builder
	for _, r := range s {
		_, err := fmt.Fprintf(&b, "%.8b", r)
		if err != nil {
			return "", err
		}
	}
	return b.String(), nil
}

func Gzip(s string) (string, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write([]byte(s))
	if err != nil {
		return "", err
	}
	gw.Close()
	return string(buf.Bytes()), nil
}

func UnicodeEscape(s string) (string, error) {
	var b strings.Builder
	for _, r := range s {
		fmt.Fprintf(&b, "\\u%04x", r)
	}
	return b.String(), nil
}

func XOREncode(input, key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("key cannot be empty")
	}
	result := make([]byte, len(input))
	for i := range input {
		result[i] = input[i] ^ key[i%len(key)]
	}
	return hex.EncodeToString(result), nil
}
