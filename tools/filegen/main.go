/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Tool to generate files
*/

package main

import (
	"log"

	"github.com/aosorgin/gotools/tools/filegen/fglib"
	"github.com/pkg/errors"
)

func getGenerator() (fglib.DataGenerator, error) {
	if fglib.Options.GeneratorType == fglib.GeneratorCrypto {
		return fglib.CreateMutliThreadGenerator(fglib.CreateCryptoDataGenerator(), fglib.CreateUnorderedQueue())
	} else if fglib.Options.GeneratorType == fglib.GeneratorPseudo {
		dataGen, err := fglib.CreatePseudoRandomDataGenerator(fglib.Options.Seed)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create pseudo-random generator")
		}
		return fglib.CreateMutliThreadGenerator(dataGen, fglib.CreateOrderedQueue())
	} else if fglib.Options.GeneratorType == fglib.GeneratorNull {
		return fglib.CreateMutliThreadGenerator(fglib.CreateNullDataGenerator(), fglib.CreateUnorderedQueue())
	}

	panic("Invalid generator type")
}

func generateFiles(options *fglib.CmdOptions) {
	gen, err := getGenerator()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to initialize generator"))
	}
	filesGen := fglib.CreateLinearFileGenerator(gen, options.Path, options.GenerateInMultipleThread,
		options.Generate.Folders, fglib.CreatePrefixNameGenerator("dir_"),
		options.Generate.Files, fglib.CreatePrefixNameGenerator("file_"), options.Generate.FileSize,
		options.QuietMode)

	defer func() {
		err = filesGen.Close()
		if err != nil {
			log.Print(errors.Wrap(err, "Failed to close generator"))
		}
	}()

	err = filesGen.Generate()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to generate files"))
	}
}

func changeFiles(options *fglib.CmdOptions) {
	gen, err := getGenerator()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to initialize generator"))
	}

	modifier := fglib.CreateFilesModifierWithInterval(gen, options.Path, options.GenerateInMultipleThread,
		options.Change.Ratio, options.Change.Interval, options.Change.Once, options.Change.Reverse,
		options.Change.Append, options.QuietMode)

	defer func() {
		err = modifier.Close()
		if err != nil {
			log.Print(errors.Wrap(err, "Failed to close generator"))
		}
	}()

	err = modifier.Modify()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to modify files"))
	}
}

func main() {
	options := fglib.ParseCmdOptions()
	switch fglib.Options.Command {
	case fglib.CommandGenerate:
		generateFiles(options)
	case fglib.CommandChange:
		changeFiles(options)
	}

}
