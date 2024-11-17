package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Brum3ns/firefly/pkg/httpreflect"
)

func Test_HttpReflectSimple(t *testing.T) {
	// Variable values to test httpReflect against
	var (
		verifyPayload        = "1333337"
		prefixLen            = 3
		suffixLen            = 3
		dummySinglelineTexts = []string{
			"__REFLECT__",
			"__REFLECT__ ",
			" __REFLECT__ ",
			"__REFLECT__l\nf",
			"l\nf__REFLECT__",
			"\n__REFLECT__\n",
			"\n__REFLECT__",
			"__REFLECT__\n",
			"\t__REFLECT__\t",
			"s__REFLECT__e",
			"s__REFLECT__",
			"__REFLECT__e",
			"abc__REFLECT__",
			"__REFLECT__end",
			"12__REFLECT__12356",
			"123456__REFLECT__12",
			"testing__REFLECT__thisstuff",
			"I do not reflect",
		}
	)
	var lst_surr []httpreflect.Surrounding

	// Configure and make new httpreflect
	httpReflect := httpreflect.NewReflect(httpreflect.Config{
		IndexEndLength:   prefixLen,
		IndexStartLength: suffixLen,
		Canary:           verifyPayload,
	})

	// Test with dummy text data to get prefix/suffix (surrounding)
	for _, txt := range dummySinglelineTexts {
		surr, ok := httpReflect.GetCanarySurrounding(
			strings.ReplaceAll(txt, "__REFLECT__", verifyPayload),
		)
		if ok {
			lst_surr = append(lst_surr, surr)
			fmt.Printf("\n[DEBUG] %+v\n", surr) // DEBUG
		}
	}

}
