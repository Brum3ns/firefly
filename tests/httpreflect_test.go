package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Brum3ns/firefly/pkg/httpreflect"
)

var PAYLOADS = []string{
	"j89D!J9O1ck09c1",
	"x",
	"jd1j9D(!1j3987hf19SomeSuperLongPayload9318032182183021820309312Ij891j)DH!/(1d78)d1f10",
	"",
}

func Test_HttpReflect(t *testing.T) {
	// Variable values to test httpReflect against
	var (
		verifyPayload = "1333337"
		prefixLen     = 33
		suffixLen     = 33
		// Note : keep dummy values below on same level when testing (if possible)
		dummyText = getDummyText_level4()
		// Note : getDummyText_level5 is "impossible" to get surrounding from
		dummyExtractText = doExtract_level4()
	)

	// Configure and make new httpreflect
	httpReflect := httpreflect.NewReflect(httpreflect.Config{
		IndexEndLength:   prefixLen,
		IndexStartLength: suffixLen,
		Canary:           verifyPayload,
	})

	// Get all the reflected values from the dummy text data (prefix/suffix, Aka: The surroundings)
	fmt.Println("========[ SURROUNDING ]========") //DEBUG
	surrs := httpReflect.GetAllCanarySurroundings(
		strings.ReplaceAll(dummyText, "__REFLECT__", verifyPayload),
	)
	if len(surrs) == 0 {
		fmt.Println("No reflected value in the dummyText")
	}

	// Print the result for all surroundings found of the reflected values
	for _, surr := range surrs {
		fmt.Printf("\n[DEBUG] %+v\n", surr) // DEBUG
	}

	fmt.Println("==========[ EXTRACT ]==========")

	// Extract values from wihtin our surrounding(s)
	values := httpReflect.ExtractAll(
		dummyExtractText,
		surrs,
	)

	fmt.Println(strings.Join(values, "\n"))
}

func doExtract_level4() string {
	return PAYLOADS[0]
}

func doExtract_level3() string {
	return strings.TrimSpace(fmt.Sprintf(`
this here%sis some random%sstring that reflect
%sand this ends here,
right?
	`, PAYLOADS[0], PAYLOADS[1], PAYLOADS[2]))
}

// This is not possible to get surroundings from since the canaries are combined within any separator
// Very unlikely to appear in a HTTP response, but still need to be handled within the httpreflect package.
func getDummyText_level5() string {
	return strings.TrimSpace(`
	__REFLECT____REFLECT____REFLECT__
	`)
}

func getDummyText_level4() string {
	return "__REFLECT__"
}

func getDummyText_level3() string {
	return strings.TrimSpace(`
this here__REFLECT__is some random__REFLECT__string that reflect
__REFLECT__and this ends here,
right?
	`)
}

func getDummyText_level2() string {
	return strings.TrimSpace(`
a__REFLECT__b__REFLECT__c__REFLECT__d
	`)
}

func getDummyText_level1() string {
	return strings.TrimSpace(`
	__REFLECT__<!DOCTYPE html>
<html lang="en">
<body>
    <header>
        <h1 class="reflective">Welcome to the __REFLECT__ Page</h1>
    </header>
    <nav>
        <ul>
            <li><a href="#home" class="reflective-link">Home</a></li>
            <li><a href="#about" class="reflective-link">About</a></li>
            <li><a href="#contact" class="reflective-link">Contact</a></li>
        </ul>
    
	<div>
		<p>
		This is some random page that makes no sense, it's all about reflecting from left to right to mid to idk where.
		<p>
		We can skip the tags add a dnu81170__REFLECT__981nd1 with some random numbers or mix it as:__REFLECT__wazuuup!

		<pre>Let's see how many __REFLECT__ED values we can find here!<pre>
	</div>
    <footer class="reflective-footer">
        <p>&copy; 2024 __REFLECT__ Corporation. All rights reserved.</p>
    </footer>
</body>
</html>__REFLECT__
	`)
}
