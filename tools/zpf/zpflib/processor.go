/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Compressing of files
 */

package zpflib

import
(
	"archive/zip"
	"bufio"
	"compress/flate"
	"io"
	"path/filepath"
	"os"
	"fmt"
	"runtime"
	"time"
)

func compress_file(srcPath string, destPath string, processedSignal *chan bool) (writen int64, err error) {
	defer func() { *processedSignal <- true }()

	src, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer src.Close()

	rawDst, err := os.Create(destPath)
	if err != nil {
		return
	}
	defer rawDst.Close()

	dst := bufio.NewWriter(rawDst)
	defer dst.Flush()

	zipWriter := zip.NewWriter(dst)
	defer  zipWriter.Close()
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		if Options.Compression == CompressionFast {
			return flate.NewWriter(out, flate.BestSpeed)
		} else if Options.Compression == CompressionBest {
			return flate.NewWriter(out, flate.BestCompression)
		} else {
			return flate.NewWriter(out, flate.DefaultCompression)
		}
	})

	_, fileName := filepath.Split(srcPath)
	zippedFile, err := zipWriter.Create(fileName)
	writen, err = io.Copy(zippedFile, src)
	if err != nil {
		return
	}
	return
}

func Compress(srcPath string, dstPath string) (err error) {
	processingQueue := make(chan bool, runtime.NumCPU() + 1)
	processedSignal := make(chan bool)
	filesCount, filesProcessed := 0, 0
	completedSignal := make(chan bool)
	completed := false
	go func () {
		filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			relPath, err := filepath.Rel(srcPath, path)
			dirPath, _ := filepath.Split(relPath)
			os.MkdirAll(filepath.Join(dstPath, dirPath), os.ModeDir)
			if err == nil {
				filesCount++
				processingQueue <- true
				go compress_file(path, filepath.Join(dstPath, relPath + ".zip"), &processedSignal)
			}
			return nil
		})
		completedSignal <- true
	}()

	timeout := time.Tick(time.Second)
	for {
		select {
		case <-processedSignal:
			filesProcessed++
			<-processingQueue
			if completed == true && filesCount == filesProcessed {
				fmt.Printf("\rProcessing/processed files: (%d/%d)          \nDone", filesCount, filesProcessed)
				return
			}
		case <- completedSignal:
			completed = true
			if filesCount == filesProcessed {
				fmt.Printf("\rProcessing/processed files: (%d/%d)          \nDone", filesCount, filesProcessed)
				return
			}
		case <-timeout:
			fmt.Printf("\rProcessing/processed files: (%d/%d)          ", filesCount, filesProcessed)
		}
	}
	return
}
