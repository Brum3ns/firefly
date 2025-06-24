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
	"surl":    SingleUrlEncode,
	"sdurl":   DoubleUrlEncode,
	"url":     UrlEscape,
	"durl":    UrlDoubleEscape,
	"html":    HTMLEscape,
	"htmle":   HTMLEquivalent,
	"base64":  Base64,
	"base32":  Base32,
	"hex":     Hex,
	"json":    JSONEncode,
	"binary":  BinaryEncode,
	"gzipb64": GzipBase64,
	"unicode": UnicodeEscape,
}

// Encode applies the selected encodings in order.
// Supports format: "xor:<key>" and others.
func Encode(payload string, encodings []string) (string, error) {
	var err error
	for _, enc := range encodings {
		if strings.HasPrefix(enc, "xor:") {
			key := strings.TrimPrefix(enc, "xor:")
			payload, err = XOREncode(payload, key)
			if err != nil {
				return "", fmt.Errorf("xor error: %w", err)
			}
			continue
		}
		encFunc, ok := encodeFuncs[strings.ToLower(enc)]
		if !ok {
			return "", fmt.Errorf("unsupported encoding: %s", enc)
		}
		payload, err = encFunc(payload)
		if err != nil {
			return "", fmt.Errorf("error in %s encoding: %w", enc, err)
		}
	}
	return payload, nil
}

// --- Encoding Implementations ---

func SingleUrlEncode(s string) (string, error) {
	return "%" + hex.EncodeToString([]byte(s)), nil
}

func DoubleUrlEncode(s string) (string, error) {
	encoded := "%" + hex.EncodeToString([]byte(s))
	return strings.ReplaceAll(encoded, "%", "%25"), nil
}

func UrlEscape(s string) (string, error) {
	return url.QueryEscape(s), nil
}

func UrlDoubleEscape(s string) (string, error) {
	escaped := url.QueryEscape(s)
	return strings.ReplaceAll(escaped, "%", "%25"), nil
}

func HTMLEscape(s string) (string, error) {
	return html.EscapeString(s), nil
}

func HTMLEquivalent(s string) (string, error) {
	escaped := html.EscapeString(s)
	escaped = strings.ReplaceAll(escaped, "&#34;", "&quot;")
	escaped = strings.ReplaceAll(escaped, "&#39;", "&apos;")
	return escaped, nil
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

func GzipBase64(s string) (string, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write([]byte(s))
	if err != nil {
		return "", err
	}
	gw.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
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
