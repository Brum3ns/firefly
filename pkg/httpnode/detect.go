package httpnode

import (
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// ResponseType is the type of response body content.
type ResponseType string

const (
	ResponseTypeUnknown ResponseType = "unknown"
	ResponseTypeHTML    ResponseType = "html"
	ResponseTypeXML     ResponseType = "xml"
	ResponseTypeJSON    ResponseType = "json"
)

// DetectResponseType inspects a HTTP response body and determines whether it's HTML, JSON, or unknown.
func DetectResponseType(body string) ResponseType {
	trimmed := strings.TrimSpace(body)

	// Try to detect JSON first
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		if json.Valid([]byte(trimmed)) {
			return ResponseTypeJSON
		}
	}

	// Try to detect HTML
	tokenizer := html.NewTokenizer(strings.NewReader(trimmed))
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		t := tokenizer.Token()
		if t.Type == html.StartTagToken || t.Type == html.SelfClosingTagToken {
			if t.Data == "html" || t.Data == "head" || t.Data == "body" || t.Data == "div" {
				return ResponseTypeHTML
			}
		}
	}

	return ResponseTypeUnknown
}

// DetectContentType inspects an HTTP response header and determines the content type.
func DetectContentType(header http.Header, xmlAsHTML bool) ResponseType {
	contentType := strings.ToLower(header.Get("Content-Type"))
	switch {
	case strings.Contains(contentType, "json"):
		return ResponseTypeJSON
	case strings.Contains(contentType, "html") || (xmlAsHTML && strings.Contains(contentType, "xml")):
		return ResponseTypeHTML
	case strings.Contains(contentType, "xml") && !xmlAsHTML:
		return ResponseTypeXML

	default:
		return ResponseTypeUnknown
	}
}
