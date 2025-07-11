package httpnode

import (
	"encoding/json"
	"net/http"
	"slices"
)

type HeaderMergeNode map[string]HeaderInfo

type HeaderInfo struct {
	Amount []int
	Values []string
}

func NewHeader() HeaderMergeNode {
	return make(HeaderMergeNode)
}

// Take HTTP headers and merge it within the prepared header
func (header HeaderMergeNode) Merge(httpheader http.Header) HeaderMergeNode {
	for h, values := range httpheader {

		// Note : The amount of values in the list "http.Header"
		// represent the amount of times the header was repeated in the HTTP response.
		amount := len(values)

		// Check existing header if we have the header
		if hasHeader, ok := header[h]; !ok {
			header[h] = HeaderInfo{
				Amount: []int{amount},
				Values: values,
			}
		} else {
			// Add unique amount of header that repeated in the HTTP response
			if !slices.Contains(hasHeader.Amount, amount) {
				hasHeader.Amount = append(hasHeader.Amount, amount)
			}

			// Add unique values from the HTTP header
			// The list is usually very short and major of case only contain one item
			for _, value := range values {
				// The item is unique, then add it
				if !slices.Contains(hasHeader.Values, value) {
					hasHeader.Values = append(hasHeader.Values, value)
				}
			}

			// Update the values in the current header from the new HTTP header
			header[h] = hasHeader
		}
	}
	return header
}

// Take a http.header and return a prepared header node
func GetHeaderNode(httpheader http.Header) HeaderMergeNode {
	header := NewHeader()
	// Note : We do not need to check for duplicates since http.Header is a map in it's core.
	// We also do return a fresh new Header type.
	for h, values := range httpheader {

		// Note : The amount of values in the list "http.Header"
		// represent the amount of times the header was repeated in the HTTP response.
		amount := len(values)

		header[h] = HeaderInfo{
			Amount: []int{amount},
			Values: values,
		}
	}
	return header
}

func HeaderNodeToJson(htmlnode HeaderMergeNode) ([]byte, error) {
	jsonBytes, err := json.Marshal(htmlnode)
	if err != nil {
		return jsonBytes, err
	}
	return jsonBytes, nil
}
