package faker

import (
	"fmt"
	"strings"

	"github.com/Brum3ns/firefly/pkg/design"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/charmbracelet/lipgloss"
)

var placeholders = map[string]string{
	"?": "__QUESTION_MARK__",
	"#": "__HASHTAG__",
	"*": "__ASTERISK__",
}

// Generate custom insert functions
/* func init() {
	gofakeit.AddFuncLookup("random", gofakeit.Info{
		Category:    "custom",
		Description: "Random name in lowercase",
		Example:     "john doe",
		Output:      "string",
		Generate: func(f *gofakeit.Faker, m *gofakeit.MapParams, info *gofakeit.Info) (any, error) {
			return "", nil
		},
	})
}
*/
// Same as gofakeit (v7) Generate, but ignore wildcard such as ? # to be replaced
func Generate(dataVal string) (string, error) {
	dataVal = replace(dataVal)
	dataVal, err := gofakeit.Generate(dataVal)
	if err != nil {
		return dataVal, err
	}
	dataVal = replaceReverse(dataVal)

	return dataVal, nil
}

func PrintList() {
	// Style for normal text
	normal := lipgloss.NewStyle().Foreground(lipgloss.Color(design.COLOR_GREY))
	// Style just for the colored word
	highlight := lipgloss.NewStyle().Foreground(lipgloss.Color(design.COLOR_ORANGE)).Bold(true)
	example := lipgloss.NewStyle().Foreground(lipgloss.Color(design.COLOR_RED))

	for name, info := range gofakeit.FuncLookups {
		fmt.Printf(
			"> %s, %s,  Ex: %s\n",
			highlight.Render(name),
			normal.Render(info.Description),
			example.Render(info.Example),
		)
	}
	fmt.Printf("Available placeholders: [%d]\n", len(gofakeit.FuncLookups))
}

func replace(dataVal string) string {
	for repl, val := range placeholders {
		dataVal = strings.ReplaceAll(dataVal, repl, val)
	}
	return dataVal
}

func replaceReverse(dataVal string) string {
	for repl, val := range placeholders {
		dataVal = strings.ReplaceAll(dataVal, val, repl)
	}
	return dataVal
}

/* type Insert struct {
	Keyword string
	Payload string
}

// Take the "Insert" structure and the "Request" structure
// Return the "Request" structure with the insert points replaced by the current payload.
func NewInsert(keyword, payload string) Insert {
	return Insert{
		Keyword: keyword,
		Payload: payload,
	}
}

// Insert the payload based on the insert point (default=FUZZ) from user options to a string
func (ist Insert) addKeyword(s string) string {
	return strings.ReplaceAll(random.RandomInsert(s), ist.Keyword, ist.Payload)
}

func (ist Insert) SetHeaders(sliceArry [][2]string) http.Header {
	var headers = http.Header{}
	for _, h := range sliceArry {
		hName := random.RandomInsert(h[0])
		hValue := random.RandomInsert(h[1])
		headers.Add(ist.addKeyword(hName), ist.addKeyword(hValue))
	}
	return headers
}

func (ist Insert) SetURL(s string) string {
	return strings.ReplaceAll(random.RandomInsert(s), ist.Keyword, normalizeURLstring(ist.Payload))
}

func (ist Insert) SetPostBody(s string) string {
	return ist.addKeyword(s)
}

func (ist Insert) SetMethod(s string) string {
	return ist.addKeyword(s)
}

// Normalize common characters in the URL into URL-encode:
func normalizeURLstring(s string) string {
	var (
		l_find        = []string{" ", "\t", "\n", "#", "&", "?"}
		l_URLEncodeTo = []string{"%20", "%09", "%0a", "%23", "%26", "%3F"}
	)
	for i := 0; i < len(l_URLEncodeTo); i++ {
		if strings.Contains(s, l_find[i]) {
			s = strings.ReplaceAll(s, l_find[i], l_URLEncodeTo[i])
		}
	}
	return s
}
*/
