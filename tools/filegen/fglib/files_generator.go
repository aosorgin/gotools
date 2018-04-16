/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Files writer for random generator
*/

package fglib

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

func generateToFile(file *os.File, gen DataGenerator, size int64) (int64, error) {
	bufFile := bufio.NewWriter(file)
	defer bufFile.Flush()

	writen, err := io.CopyN(bufFile, gen, size)
	if err != nil {
		return writen, errors.Wrap(err, "Failed to generate data to file")
	}

	return writen, nil
}

func writeFile(path string, size uint64, gen DataGenerator) error {
	rawFile, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "Failed for create '%s'", path)
	}
	defer rawFile.Close()

	if _, err = generateToFile(rawFile, gen, int64(size)); err != nil {
		return errors.Wrapf(err, "Failed to write to '%s'", path)
	}

	return nil
}

type NameGenerator interface {
	GetName(index uint) (string, error)
}

type prefixNameGenerator struct {
	prefix string
}

func (ng *prefixNameGenerator) GetName(index uint) (string, error) {
	return fmt.Sprintf("%s%d", ng.prefix, index), nil
}

func CreatePrefixNameGenerator(prefix string) NameGenerator {
	return &prefixNameGenerator{prefix: prefix}
}

type FilesGenerator interface {
	io.Closer
	Generate() error
}

type linearFilesGenerator struct {
	gen                       DataGenerator
	path                      string
	generateInMultipleThreads bool
	printProgress             bool

	dirsCount uint
	dirNames  NameGenerator

	filesCount uint
	fileNames  NameGenerator
	fileSize   uint64
}

func (g *linearFilesGenerator) Close() error {
	return g.gen.Close()
}

func (g *linearFilesGenerator) Generate() error {
	log.Print("Files generation is started")
	log.Printf("directories: %d, files: %d, file size: %d", g.dirsCount, g.filesCount, g.fileSize)
	log.Printf("write in the single thread: %t", g.generateInMultipleThreads == false)
	completeSignal := make(chan bool)
	errorChannel := make(chan error)
	filesGenerated := 0

	go func() {
		if g.generateInMultipleThreads == false {
			// Attach goroutine to the single thread to wrile files on disk with no fragmentation
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
		}

		for i := uint(0); i < g.dirsCount; i++ {
			dirName, err := g.dirNames.GetName(i)
			if err != nil {
				errorChannel <- errors.Wrap(err, "Failed to generate directory name")
				return
			}
			folderPath := filepath.Join(g.path, dirName)
			os.MkdirAll(folderPath, os.ModeDir|0755)
			for j := uint(0); j < g.filesCount; j++ {
				fileName, err := g.fileNames.GetName(j)
				if err != nil {
					errorChannel <- errors.Wrap(err, "Failed to generate file name")
				}
				filePath := filepath.Join(folderPath, fileName)
				err = writeFile(filePath, g.fileSize, g.gen)
				if err != nil {
					errorChannel <- errors.Wrapf(err, "Failed to generate file '%s'", filePath)
				}
				filesGenerated++
			}
		}
		completeSignal <- true
	}()

	timeout := time.Tick(time.Second)
	filesTotal := Options.Generate.Folders * Options.Generate.Files
	for {
		select {
		case <-timeout:
			if g.printProgress {
				fmt.Printf("\rGenerated: (%d/%d)        ", filesGenerated, filesTotal)
			}
		case <-completeSignal:
			fmt.Printf("\rGenerated: (%d/%d)        \n", filesGenerated, filesTotal)
			log.Print("Generation has completed successfully")
			return nil
		case err := <-errorChannel:
			return errors.Wrapf(err, "Failed to generate files")
		}
	}
}

func CreateLinearFileGenerator(gen DataGenerator, path string, generateInMultipleThreads bool,
	dirsCount uint, dirNames NameGenerator,
	filesCount uint, fileNames NameGenerator, fileSize uint64,
	quietMode bool) FilesGenerator {
	return &linearFilesGenerator{
		gen:  gen,
		path: path,
		generateInMultipleThreads: generateInMultipleThreads,
		printProgress:             quietMode == false,
		dirsCount:                 dirsCount,
		dirNames:                  dirNames,
		filesCount:                filesCount,
		fileNames:                 fileNames,
		fileSize:                  fileSize,
	}
}
