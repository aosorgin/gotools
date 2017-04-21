/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Command-line options parser and logic
 */

package zpflib

import (
	"flag"
	"fmt"
	"os"
)

// CompressionEnum
const (
	CompressionNormal = iota
	CompressionFast
	CompressionBest
)

type CmdOptions struct {
	SrcPath string
	DestPath string
	Compression int // used values from CompressionEnum
}

var Options CmdOptions

func checkOptions()  {
	if (Options.SrcPath == "") {
		fmt.Fprintf(os.Stderr, "Error: Use --src option to set source path for files to compress\n")
		flag.Usage()
		os.Exit(1)
	}

	if (Options.DestPath == "") {
		fmt.Fprintf(os.Stderr, "Error: Use --dest option to set destination path for compressed\n")
		flag.Usage()
		os.Exit(1)
	}
}

func parseCompression(compression string) {
	if compression == "fast" {
		Options.Compression = CompressionFast
	} else if compression == "normal" {
		Options.Compression = CompressionNormal
	} else if compression == "best" {
		Options.Compression = CompressionNormal
	} else {
		fmt.Fprintf(os.Stderr, "Error: Invalid compression options for --comp: '%s'\n", compression)
		flag.Usage()
		os.Exit(1)
	}
}

func PrepareCmdOptions() {
	/* Initializing flags for parsing command-line arguments */
	flag.StringVar(&Options.SrcPath, "src", "", "Source directory path")
	flag.StringVar(&Options.DestPath, "dest", "", "Destination directory path")
	comp := flag.String("comp", "normal", "Compression to use (normal is default, fast, best)")

	/* Parsing command-line */
	flag.Parse()
	checkOptions()
	parseCompression(*comp)
}