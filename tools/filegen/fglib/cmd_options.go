/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Command-line options parser and logic
*/

package fglib

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// CommandEnum
const (
	CommandGenerate = iota
	CommandChange
)

// GeneratorEnum
const (
	GeneratorCrypto = iota
	GeneratorPseudo
)

type CmdOptions struct {
	Command       int    // CommandEnum
	Path          string // Root path got processing files
	GeneratorType int    // GeneratorEnum
	Seed          []byte
	Generate      struct {
		Folders  uint   // Folders tree count. For example [3,4] - 3 folders in root, 4 folders in previous, etc.
		Files    uint   // Files count in each tree level
		FileSize uint64 // File size for each tree level
	}
	Change struct {
		Interval IntervalType // Interval to change files
		Once     bool         // Use once if true otherwise until the end of file
		Reverse  bool         // Change file from end if true
	}
}

var Options CmdOptions

func processInterval(interval string) {
	if interval == "" {
		Options.Change.Interval = GetEmptyInterval()
	} else {
		var err error
		Options.Change.Interval, err = ParseInterval(interval)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
func processCommonCommand() {
	/* Check options */

	if len(Options.Path) == 0 {
		fmt.Fprintf(os.Stderr, "Error: path is not set. Use the --path option.\n")
		flag.Usage()
		os.Exit(1)
	}
}

func processGenerateCommand() {
	if Options.Generate.Folders == 0 {
		fmt.Fprintf(os.Stderr, "Error: Use the --folders option to set files count to generate.\n")
		flag.Usage()
		os.Exit(1)
	}

	if Options.Generate.Files == 0 {
		fmt.Fprintf(os.Stderr, "Error: Use the --files option to set files count to generate.\n")
		flag.Usage()
		os.Exit(1)
	}
}

func processGeneratorType(genType string, seed uint64) {
	if genType == "crypto" {
		Options.GeneratorType = GeneratorCrypto
		if seed != 0 {
			fmt.Fprintf(os.Stderr, "Warning: seed is not used with crypto generator.\n")
		}
	} else if genType == "pseudo" {
		Options.GeneratorType = GeneratorPseudo
		if seed == 0 {
			seed = uint64(time.Now().UnixNano())
		}
		Options.Seed = SeedFromUint64(seed)
		fmt.Println("Using seed:", seed)
	} else {
		fmt.Fprintf(os.Stderr, "Error: invalid generator type '%s'.\n", genType)
		flag.Usage()
		os.Exit(1)
	}
}

func processCommand(cmd string) {
	if cmd == "gen" {
		Options.Command = CommandGenerate
		processCommonCommand()
		processGenerateCommand()
	} else if cmd == "change" {
		Options.Command = CommandChange
		processCommonCommand()
	} else {
		fmt.Fprintf(os.Stderr, "Error: Invalid command '%s'\n", cmd)
		flag.Usage()
		os.Exit(1)
	}
}

func ParseCmdOptions() {
	/* Initializing flags for parsing command-line arguments */

	flag.UintVar(&Options.Generate.Files, "files", 0, "Number of files to generate")
	flag.UintVar(&Options.Generate.Folders, "folders", 1, "Number of folders to generate")
	flag.Uint64Var(&Options.Generate.FileSize, "size", 0, "Size of files to generate")

	flag.BoolVar(&Options.Change.Once, "once", false, "Process interval only once. By default is false")
	flag.BoolVar(&Options.Change.Reverse, "reverse", false, "Process interval from the end of file. By default is false")
	interval := flag.String("interval", "", "Interval to change file. By default do not changes the file")

	flag.StringVar(&Options.Path, "path", "", "Path to root folder to generate files")

	genType := flag.String("gen-type", "crypto", "Generator type. [(random), pseudo]")
	seed := flag.Uint64("seed", 0, "Seed for data generator. Could be used with pseudo generator")

	cmd := flag.String("cmd", "gen", "Command to execute. Could be [(gen, change)]")

	/* Parsing command-line */
	flag.Parse()
	processInterval(*interval)
	processCommand(*cmd)
	processGeneratorType(*genType, *seed)
}
