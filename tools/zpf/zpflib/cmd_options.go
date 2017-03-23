package zpflib

import "flag"

type CmdOptions struct {
	SrcPath string
	DestPath string
}

var Options CmdOptions

func PrepareCmdOptions() {
	/* Initializing flags for parsing command-line arguments */
	flag.StringVar(&Options.SrcPath, "src", "", "Source directory path")
	flag.StringVar(&Options.DestPath, "dest", "", "Destination directory path")

	/* Parsing command-line */
	flag.Parse()
}