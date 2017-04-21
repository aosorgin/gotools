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
	"reflect"
	"strings"
	"strconv"
)

// CommandEnum
const (
	CommandGenerate = iota
)

type CmdOptions struct {
	Command  int    // CommandEnum
	Path     string // Root path got processing files
	Generate struct {
		Folders  []uint32 // Folders tree count. For example [3,4] - 3 folders in root, 4 folders in previous, etc.
		Files    []uint32 // Files count in each tree level
		FileSize []uint64 // File size for each tree level
	}
}

var Options CmdOptions

func collectGenerateOptions(folders string, files string, size string) {
	processSlice := func(itemsStr string, t reflect.Type) (res reflect.Value) {
		items := strings.Split(itemsStr, ",")
		res = reflect.MakeSlice(t, len(items), len(items))
		for i := range items {
			res[i] = strconv.ParseUint(items[i], 10, 64)
		}
	}
}

func checkOptions() {
	if Options.SrcPath == "" {
		fmt.Fprintf(os.Stderr, "Error: Use --src option to set source path for files to compress\n")
		flag.Usage()
		os.Exit(1)
	}

	if Options.DestPath == "" {
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
