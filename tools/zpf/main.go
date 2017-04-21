/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Tool to compress the single files into separate zip files with storing file hierarchy
 */


package main

import (
	"fmt"
	"github.com/aosorgin/gotools/tools/zpf/zpflib"
	"runtime"
)

func main() {
	zpflib.PrepareCmdOptions()
	maxThreads := runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Maximum threads:", maxThreads)
	zpflib.Compress(zpflib.Options.SrcPath, zpflib.Options.DestPath)
}
