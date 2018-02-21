/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Files writer for random generator
*/

package fglib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

func writeFile(path string, size uint64, gen DataGenerator) {
	rawFile, err := os.Create(path)
	if err != nil {
		log.Print(errors.Wrapf(err, "Failed for create '%s'", path))
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
			log.Print(errors.Wrapf(err, "Failed to write to '%s'", path))
		}
		file.Write(buffer)
		size -= uint64(len(buffer))
	}
}

type DataWriter struct {
	gen DataGenerator
}

func (w *DataWriter) Init(gen DataGenerator) {
	w.gen = gen
}

func (w *DataWriter) Close() error {
	return w.gen.Close()
}

func (w *DataWriter) WriteFiles() error {
	completeSignal := make(chan bool)
	filesGenerated := 0

	go func() {
		for i := uint(0); i < Options.Generate.Folders; i++ {
			folderPath := filepath.Join(Options.Path, fmt.Sprintf("dir%04d", i))
			os.MkdirAll(folderPath, os.ModeDir|0755)
			for j := uint(0); j < Options.Generate.Files; j++ {
				writeFile(filepath.Join(folderPath, fmt.Sprintf("file%04d", j)), Options.Generate.FileSize, w.gen)
				filesGenerated += 1
			}
		}
		completeSignal <- true
	}()

	timeout := time.Tick(time.Second)
	filesTotal := Options.Generate.Folders * Options.Generate.Files
	for {
		select {
		case <-timeout:
			fmt.Printf("\rGenerated: (%d/%d)        ", filesGenerated, filesTotal)
		case <-completeSignal:
			fmt.Printf("\rGenerated: (%d/%d)        ", filesGenerated, filesTotal)
			return nil
		}
	}
	return nil
}
