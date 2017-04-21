/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Files writer for random generator
*/

package fglib

import ()
import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func writeFile(path string, size uint64, gen *CryptoGenerator) {
	rawFile, err := os.Create(path)
	if err != nil {
		return
	}
	defer rawFile.Close()

	file := bufio.NewWriter(rawFile)
	defer file.Flush()

	var bufferSize uint64 = 64 * 1024
	buffer := make([]byte, bufferSize)

	for size > 0 {
		if size < bufferSize {
			buffer = buffer[:size]
		}
		_, err = gen.Read(buffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to write '%'.\n", path)
		}
		file.Write(buffer)
		size -= uint64(len(buffer))
	}
}

func WriteFiles(gen *CryptoGenerator) {
	completeSignal := make(chan bool)
	filesGenerated := 0

	gen.Init()

	go func() {
		var i uint = 0
		for ; i < Options.Generate.Files; i++ {
			writeFile(filepath.Join(Options.Path, fmt.Sprintf("file%d", i)), Options.Generate.FileSize, gen)
			filesGenerated += 1
		}
		completeSignal <- true
	}()

	timeout := time.Tick(time.Second)
	for {
		select {
		case <-timeout:
			fmt.Printf("\rGenerated: (%d/%d)        ", filesGenerated, Options.Generate.Files)
		case <-completeSignal:
			fmt.Printf("\rGenerated: (%d/%d)        ", filesGenerated, Options.Generate.Files)
			gen.Stop()
			return
		}
	}
}
