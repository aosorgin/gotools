package zpflib

import
(
	"archive/zip"
	"bufio"
	"compress/flate"
	"io"
	"path/filepath"
	"os"
	"sync"
	"fmt"
	"runtime"
	"time"
)

func compress_file(srcPath string, destPath string, sem *chan bool, wg *sync.WaitGroup, fc *int) (writen int64, err error) {
	defer wg.Done()
	defer func() { <-*sem }()
	defer func() { *fc -= 1	}()

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
		return flate.NewWriter(out, flate.DefaultCompression)
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
	var wg sync.WaitGroup
	sem := make (chan bool, runtime.NumCPU() + 1)
	filesCount, filesRemain := 0, 0
	fmt.Print("Processing files: 0")
	lastReportTime := time.Now()
	filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(srcPath, path)
		dirPath, _ := filepath.Split(relPath)
		os.MkdirAll(filepath.Join(dstPath, dirPath), os.ModeDir)
		if err == nil {
			wg.Add(1)
			filesCount += 1
			filesRemain += 1
			sem <- true
			go compress_file(path, filepath.Join(dstPath, relPath + ".zip"), &sem, &wg, &filesRemain)
			if time.Since(lastReportTime) > time.Second {
				lastReportTime = time.Now()
				fmt.Printf("\rProcessing files: %d          ", filesCount)
			}
		}
		return err
	})
	fmt.Println("")

	completed := false
	go func() {
		wg.Wait()
		completed = true
	}()
	for completed == false {
		fmt.Printf("\rRemain to process: %d         ", filesRemain)
		time.Sleep(time.Second)
	}
	fmt.Printf("\rRemain to process: 0         \nDone")
	return
}
