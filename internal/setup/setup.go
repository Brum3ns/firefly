package setup

import (
	"log"

	"github.com/Brum3ns/firefly/internal/option"
	"github.com/Brum3ns/firefly/internal/output"
)

func Setup(opt option.Option) error {

	// Check output options
	if output.FileExist(opt.Output.OutputFile) {
		if opt.Output.Overwrite {
			// Remove old output file
			output.RemoveFile(opt.Output.OutputFile)
		} else {
			log.Fatalln("the given output file already exists. To overwrite, use the flag: -oW")
		}
	}

	// Check analyze output options
	if output.FileExist(opt.Output.OutputFileAnalyze) {
		if opt.Output.Overwrite {
			// Remove old output file
			output.RemoveFile(opt.Output.OutputFileAnalyze)
		} else {
			log.Fatalln("the given output file already exists (analyze result). To overwrite, use the flag: -oW")
		}
	}
	return nil
}
