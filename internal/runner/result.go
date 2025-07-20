package runner

import (
	"context"
	"log"

	"github.com/Brum3ns/firefly/internal/output"
	"github.com/Brum3ns/firefly/pkg/extract"
)

func (r *Runner) handlerResult(ctx context.Context) {
	defer r.wg.process.Done()
	for {
		select {
		case scannerResult := <-r.channel.result:
			r.wg.scanner.Done()
			// Analyze result - Check if the scanner got a hit/miss in the result
			if scannerResult.scanResult.OK {
				// Make the final output result to be saved to output file
				outputJson, err := output.MakeOutputFileJSON(
					output.Output{
						Date: getTime(),
						Input: output.Input{
							Payload: scannerResult.coreJob.payload,
						},
						Scan: output.Scan{
							Httpdiff: scannerResult.scanResult.HTTPDiff,
							Extract: extract.Result{
								Regexes:  scannerResult.scanResult.Extract.Body.Regexes,
								Keywords: scannerResult.scanResult.Extract.Body.Keywords,
							},
						},
						Http: output.Http{
							Response:   scannerResult.coreJob.httpResponse,
							RawRequest: string(scannerResult.coreJob.httpRawRequest.TemplateRawRequest),
						},
					},
					r.option.Output.OutputHTTPResponseBody,
				)
				if err != nil {
					log.Println(err)
				}

				// Make output to file
				if r.option.Output.OutputFile != "" {
					err := output.AppendOutputJsonToFile(outputJson, r.option.Output.OutputFile)
					if err != nil {
						log.Fatalln(err)
					}
				}
				// Check analyze output
				if r.option.Output.OutputFileAnalyze != "" {
					r.appendAnalyzeResult(outputJson)
				}
			}

			r.wg.verbose.Add(1)
			r.channel.statistic <- scannerResult

		case <-ctx.Done():
			return
		}
	}
}
