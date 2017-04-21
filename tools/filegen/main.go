/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Tool to generate files
*/

package main

import (
	"github.com/aosorgin/gotools/tools/filegen/fglib"
)

func main() {
	fglib.ParseCmdOptions()

	var gen fglib.CryptoGenerator
	fglib.WriteFiles(&gen)
}
