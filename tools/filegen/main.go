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

func generateFiles() {
	var writer fglib.DataWriter

	gen, err := getGenerator()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to initialize generator"))
	}
	writer.Init(gen)

	defer func() {
		err = writer.Close()
		if err != nil {
			log.Print(errors.Wrap(err, "Failed to close generator"))
		}
	}()

	err = writer.WriteFiles()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to generate files"))
	}
}

func changeFiles() {
	var changer fglib.Changer

	gen, err := getGenerator()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to initialize generator"))
	}

	changer.Init(gen, fglib.Options.Change.Interval, fglib.Options.Change.Once,
		fglib.Options.Change.Reverse)

	defer func() {
		err = changer.Close()
		if err != nil {
			log.Print(errors.Wrap(err, "Failed to close generator"))
		}
	}()

	err = changer.ModifyFiles()
	if err != nil {
		log.Print(errors.Wrap(err, "Failed to modify files"))
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
