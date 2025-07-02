package tests

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Brum3ns/firefly/pkg/extract"
)

func Test_extract_prefix(t *testing.T) {

	e, err := extract.NewExtract(extract.Config{
		FilenameWordlist: "../wordlists/wordlist-pattern.txt",
		FilenameRegexes:  "../wordlists/wordlist-regex.txt",
	})
	if err != nil {
		log.Fatalln(err)
	}

	result := e.Run(dummySource)

	for k, v := range result.Keywords {
		fmt.Printf("%v => %+v\n", k, v)
	}
	fmt.Println(strings.Repeat("-", 22))
	for k, v := range result.Regexes {
		fmt.Printf("%v => %+v\n", k, v)
	}

}

const dummySource = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>PrefixMap Test</title>
</head>
<body>
  <h1>Test Word List</h1>
  <ul id="word-list">
    <li>friend</li>
    <li>fries</li>
    <li>friday</li>
    <li>train</li>
    <li>training</li>
    <li>transformers</li>
	<li>transformers</li>
	<li>transformers</li>
    <li>pasta321</li>
    <li>pantb</li>
    <li>alone616</li>
    <li>12345</li>
    <li>1532</li>
    <li>1267</li>
    <li>pasto</li>
  </ul>
</body>
</html>
`
