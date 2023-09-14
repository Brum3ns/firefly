// supported encode format for Firfly
package encode

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strings"
)

var (
	encodeTo = map[string]func(string) string{
		"surl":   func(s string) string { return Url(s) },
		"sdurl":  func(s string) string { return DoubleUrl(s) },
		"url":    func(s string) string { return Url_smart(s) },
		"durl":   func(s string) string { return DoubleUrl_smart(s) },
		"base64": func(s string) string { return Base64(s) },
		"base32": func(s string) string { return Base32(s) },
		"html":   func(s string) string { return HTMLEsacpe(s) },
		"htmle":  func(s string) string { return HTMLEquivalent(s) },
		"hex":    func(s string) string { return Hex(s) },
		"json":   func(s string) string { sJson, _ := json.Marshal(s); return string(sJson) },
		"binary": func(s string) string {
			var b string
			for _, r := range s {
				b = fmt.Sprintf("%s%.8b", b, r)
			}
			return b
		},
	}
)

func Encode(payload string, encodes []string) string {
	for _, encode := range encodes {
		if _, ok := encodeTo[strings.ToLower(encode)]; ok {
			payload = encodeTo[strings.ToLower(encode)](payload)
		}
	}
	return payload
}

func Url(s string) string {
	return ("%" + hex.EncodeToString([]byte(s)))
}

func DoubleUrl(s string) string {
	return strings.ReplaceAll(Url(s), "%", "%25")
}

func Url_smart(s string) string {
	return url.QueryEscape(s)
}

func DoubleUrl_smart(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "%", "%25")
}

func HTMLEsacpe(s string) string {
	return html.EscapeString(s)
}

func HTMLEquivalent(s string) string {
	return strings.ReplaceAll(html.EscapeString(s), "&#34;", "&quot;")
}

func Base32(s string) string {
	return base32.StdEncoding.EncodeToString([]byte(s))
}

func Base64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Hex(s string) string {
	return hex.EncodeToString([]byte(s))
}
