/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Command-line options parser and logic
*/

package fglib

import (
	"fmt"
	"os"
	"time"
	"io"
	"github.com/go2c/optparse"
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
	GeneratorNull
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
		Ratio    float64      // Change ratio
		Interval IntervalType // Interval to change files
		Once     bool         // Use once if true otherwise until the end of file
		Reverse  bool         // Change file from end if true
	}
}

var Options CmdOptions

func processFileSize(rawSize string) {
	var err error
	Options.Generate.FileSize, err = ParseSize(rawSize)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func processInterval(interval string) {
	if interval == "" {
		Options.Change.Interval = GetFullInterval()
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
		usage(os.Stderr)
		os.Exit(1)
	}
}

func processGenerateCommand() {
	if Options.Generate.Folders == 0 {
		fmt.Fprintf(os.Stderr, "Error: Use the --folders option to set files count to generate.\n")
		usage(os.Stderr)
		os.Exit(1)
	}

	if Options.Generate.Files == 0 {
		fmt.Fprintf(os.Stderr, "Error: Use the --files option to set files count to generate.\n")
		usage(os.Stderr)
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
	} else  if genType == "null" {
		Options.GeneratorType = GeneratorNull
		if seed != 0 {
			fmt.Fprintf(os.Stderr, "Warning: seed is not used with null generator.\n")
		}
	} else 	{
		fmt.Fprintf(os.Stderr, "Error: invalid generator type '%s'.\n", genType)
		usage(os.Stderr)
		os.Exit(1)
	}
}

func processCommand(cmd string) {
	if cmd == "gen" || cmd == "generate" {
		Options.Command = CommandGenerate
		processCommonCommand()
		processGenerateCommand()
	} else if cmd == "chg" || cmd == "change" {
		Options.Command = CommandChange
		processCommonCommand()
	} else {
		fmt.Fprintf(os.Stderr, "Error: Invalid command '%s'\n", cmd)
		usage(os.Stderr)
		os.Exit(1)
	}
}

func usage(f io.Writer) {
	fmt.Fprintln(f, "Usage:")
	fmt.Fprintf(f, "  %s [command] [options]\n\n", os.Args[0])

	fmt.Fprintln(f, "Commands:")
	fmt.Fprintln(f, "  gen, generate              Generate files")
	fmt.Fprintln(f, "  chg, change                Change files")
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Common options:")
	fmt.Fprintln(f, "  -p, --path                 Path to processing folder")
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Generate command options:")
	fmt.Fprintln(f, "  -d, --dirs                 Directories count to generate")
	fmt.Fprintln(f, "  -f, --files                Files count to generate")
	fmt.Fprintln(f, "  -s, --size                 File size to generate. Size format: [\\d{k,K,m,M,g,G}]")
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Change command options:")
	fmt.Fprintln(f, "  --scale                    Percent of files cont to change. Range: [0;1]. By default is equal to 1")
	fmt.Fprintln(f, "  -i, --interval             Interval to change file with. Format: ['data not to change', 'data to change',{'data not to change'}].")
	fmt.Fprintln(f, "                             Data format: [\\d{%, k,K,m,M,g,G}]. By default is [0,100%] and used until file ending")
	fmt.Fprintln(f, "  --once                     Using of interval only once. Used only with -i, --interval option.")
	fmt.Fprintln(f, "  --reverse                  Using interval from the file ending to begining. Used only with -i, --interval option.")
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Generator options:")
	fmt.Fprintln(f, "  -g, --generator            Type of generator to use")
	fmt.Fprintln(f, "     crypto                  Crypto random data generator. Used by default.")
	fmt.Fprintln(f, "     pseudo                  Pseudo random data generator")
	fmt.Fprintln(f, "     null                    Null contains data generator")
	fmt.Fprintln(f, "  --seed                     Initial seed for generated data. Can be used only with 'pseudo' generator")
	fmt.Fprintln(f)

	fmt.Fprintln(f, "Auxillary options:")
	fmt.Fprintln(f, "  -h, --help                 Print help and exit")
	fmt.Fprintln(f, "  -v, --version              Print version and exit")
}

func ParseCmdOptions() {
	/* Initializing flags for parsing command-line arguments */

	/* generate command options */
	optparse.UintVar(&Options.Generate.Files, "files", 'f', 0)
	optparse.UintVar(&Options.Generate.Folders, "dirs", 'd', 1)
	fileSize := optparse.String("size", 's', "0")

	/* change command option */
	optparse.FloatVar(&Options.Change.Ratio, "scale", 0, float64(1))
	optparse.BoolVar(&Options.Change.Once, "once", 0, false)
	optparse.BoolVar(&Options.Change.Reverse, "reverse", 0, false)
	interval := optparse.String("interval", 'i', "")

	/* common options */
	optparse.StringVar(&Options.Path, "path", 'p', "")

	/* generator options */
	genType := optparse.String("generator", 'g', "crypto")
	seed := optparse.Uint("seed", 0, 0)

	/* auxillary options */
	help := optparse.Bool("help", 'h', false)
	version := optparse.Bool("version", 'v', false)

	/* Parsing command-line */
	args, err := optparse.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n\n", os.Args[0], err.Error())
		usage(os.Stderr)
		os.Exit(1)
	}

	/* processing command line arguments */
	if *help {
		usage(os.Stdout)
		os.Exit(0)
	}

	if *version {
		fmt.Println("Version: 0.1.0")
		os.Exit(0)
	}

	if len(args) == 0{
		fmt.Fprintf(os.Stderr, "Error: Set command to use\n")
		usage(os.Stderr)
		os.Exit(1)
	}

	cmd := args[0]

	processFileSize(*fileSize)
	processInterval(*interval)
	processCommand(cmd)
	processGeneratorType(*genType, uint64(*seed))
}
