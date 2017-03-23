package main

import (
	"fmt"
	"tools/zpf/zpflib"
	"runtime"
)

func main() {
	zpflib.PrepareCmdOptions()
	maxThreads := runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Maximum threads:", maxThreads)
	zpflib.Compress(zpflib.Options.SrcPath, zpflib.Options.DestPath)
}
