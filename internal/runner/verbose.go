package runner

import (
	"fmt"
	"time"
)

type Verbose struct {
	time   string
	result jobScanner
}

func PrintVerbose(v Verbose) {
	var (
		hdiff             = v.result.scanResult.HTTPDiff
		htmlTagsAppear    = hdiff.HTML.Appear.TagEndHits + hdiff.HTML.Appear.TagStartHits
		htmlTagsDisappear = hdiff.HTML.Disappear.TagEndHits + hdiff.HTML.Disappear.TagStartHits
	)

	fmt.Printf(
		"%s - Payload:[%s], Code:[%d], Time:[%v], Size:[%d], Line:[%d], Word:[%d], Headers:[%d], TagsA:[%d], TagsD:[%d], Extract:[%v]\n",
		time.Now().Format("15:04:05.00"),
		v.result.coreJob.payload,
		v.result.coreJob.httpResponse.StatusCode,
		v.result.coreJob.httpResponse.Time,
		v.result.coreJob.httpResponse.BodySize,
		v.result.coreJob.httpResponse.LineCount,
		v.result.coreJob.httpResponse.WordCount,
		v.result.coreJob.httpResponse.HeaderAmount,
		htmlTagsAppear,
		htmlTagsDisappear,
		v.result.scanResult.Extract.OK,
	)
}
