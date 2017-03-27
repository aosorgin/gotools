/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Tool to compress the single files into separate zip files with storing file hierarchy
 */


package main

import (
	"github.com/aosorgin/gotools/tools/zpf/zpflib"
)

func main() {
	zpflib.PrepareCmdOptions()
	zpflib.Compress(zpflib.Options.SrcPath, zpflib.Options.DestPath)
}
