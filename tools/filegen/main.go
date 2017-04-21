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

func main() {
	fglib.ParseCmdOptions()

	var writer fglib.DataWriter
	var gen fglib.Generator
	gen.SetDataGenerator(new(fglib.CryptoGenerator))
	err := writer.Init(&gen)
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
