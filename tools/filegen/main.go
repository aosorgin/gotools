/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Tool to generate files
*/

package main

import (
	"fmt"
	"github.com/aosorgin/gotools/tools/filegen/fglib"
	"os"
)

func getGenerator() *fglib.Generator {
	var gen fglib.Generator
	if fglib.Options.GeneratorType == fglib.GeneratorCrypto {
		gen.SetDataGenerator(new(fglib.CryptoGenerator), new(fglib.UnorderedQueue))
	} else if fglib.Options.GeneratorType == fglib.GeneratorPseudo {
		var dataGen fglib.PseudoRandomGenerator
		dataGen.Seed(fglib.Options.Seed)
		gen.SetDataGenerator(&dataGen, new(fglib.OrderedQueue))
	} else if fglib.Options.GeneratorType == fglib.GeneratorNull {
		gen.SetDataGenerator(new(fglib.NullGenerator), new(fglib.UnorderedQueue))
	} else {
		panic("Invalid generator type")
	}
	return &gen
}

func generateFiles() {
	var writer fglib.DataWriter

	err := writer.Init(getGenerator())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Failed to initialize generator with error", err)
		return
	}
	defer func() {
		err = writer.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: Failed to close generator with error", err)
		}
	}()

	err = writer.WriteFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Failed to generate files with error", err)
	}
}

func changeFiles() {
	var changer fglib.Changer

	err := changer.Init(getGenerator(), fglib.Options.Change.Interval, fglib.Options.Change.Once,
		fglib.Options.Change.Reverse)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Failed to initialize generator with error", err)
		return
	}
	defer func() {
		err = changer.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: Failed to close generator with error", err)
		}
	}()

	err = changer.ModifyFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Failed to generate files with error", err)
	}
}

func main() {
	fglib.ParseCmdOptions()
	switch fglib.Options.Command {
	case fglib.CommandGenerate:
		generateFiles()
	case fglib.CommandChange:
		changeFiles()
	}

}
