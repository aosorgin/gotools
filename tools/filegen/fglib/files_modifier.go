/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Files changer for random generator
*/

package fglib

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
)

/* Get random subset in sequence mode */

type FileSelector interface {
	IsFileIsSelected() (bool, error)
}

type randomFileSelector struct { // implements SequenceChecker
	gen      DataGenerator
	weight   float64
	hitScore float64
	length   uint64
	count    uint64
	index    uint64
	hits     uint64
}

func (r *randomFileSelector) IsFileIsSelected() (bool, error) {
	if r.hits >= r.count || r.index >= r.length {
		return false, nil
	}

	if (r.count - r.hits) == (r.length - r.index) {
		r.hitScore = 0
		r.index++
		r.hits++
		return true, nil
	}

	v, err := r.getFloat64()
	if err != nil {
		return false, errors.Wrap(err, "Failed to get random float")
	}

	r.index++
	if v > r.hitScore {
		r.hits++
		r.hitScore *= r.weight
		return true, nil
	}
	return false, nil
}

func (r *randomFileSelector) getFloat64() (float64, error) {
	var v uint64
	err := binary.Read(r.gen, binary.LittleEndian, &v)
	if err != nil {
		return float64(0), errors.Wrap(err, "Failed toread float value from data generator")
	}
	return (float64(v) / float64(^uint64(0))), nil
}

func CreateRundomFileSelector(gen DataGenerator, count uint64, total uint64) (FileSelector, error) {
	if total == 0 {
		return nil, fmt.Errorf("total argument cannot be 0")
	}

	if count > total {
		return nil, fmt.Errorf("count(%d) argument cannot be more than total(%d)", count, total)
	}

	return &randomFileSelector{
		gen:      gen,
		count:    count,
		length:   total,
		weight:   float64(total-1) / float64(total),
		hitScore: float64(total-uint64((float64(count)/1.4))) / float64(total),
	}, nil
}

/* Simple checker that returns always true */

type selectAllFiles struct {
}

func (r *selectAllFiles) IsFileIsSelected() (bool, error) {
	return true, nil
}

func CreateAllFilesSelector() FileSelector {
	return &selectAllFiles{}
}

/* Changer implementation */

type FilesModifier interface {
	io.Closer
	Modify() error
}

type modifyFilesWithIntervals struct {
	gen                       DataGenerator
	path                      string
	generateInMultipleThreads bool
	printProgress             bool

	changeRatio float64

	interval Interval
	once     bool // use once if true otherwise until file ending
	reverse  bool // use interval from the end of file if true
}

func (m *modifyFilesWithIntervals) Close() error {
	return m.gen.Close()
}

func min(a int64, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func (m *modifyFilesWithIntervals) changeFile(path string, info os.FileInfo) error {
	i := GetObsoleteInterval(m.interval, info.Size())
	if i.Modify.Value == 0 {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrapf(err, "Failed to open file '%s'", path)
	}
	defer file.Close()

	size := info.Size()
	writeSize := int64(0)
	offset, newOffset := int64(0), int64(0)
	for offset < size {
		if i.NotModify.Value > 0 {
			if (size - newOffset) < i.NotModify.Value {
				break
			}
			newOffset += i.NotModify.Value
		}

		if m.reverse == true || newOffset != offset {
			if m.reverse == true {
				if (size - newOffset) < i.Modify.Value {
					_, err := file.Seek(0, io.SeekStart)
					if err != nil {
						return errors.Wrap(err, "Failed to seek")
					}
				} else {
					_, err := file.Seek(size-newOffset-i.Modify.Value, io.SeekStart)
					if err != nil {
						return errors.Wrap(err, "Failed to seek")
					}
				}
				writeSize = min(size-newOffset, i.Modify.Value)
			} else {
				_, err := file.Seek(newOffset, io.SeekStart)
				if err != nil {
					return errors.Wrap(err, "Failed to seek")
				}
				writeSize = min(i.Modify.Value, size-newOffset)
			}
			offset = newOffset
		} else {
			writeSize = min(i.Modify.Value, size-newOffset)
		}

		writen, err := io.CopyN(file, m.gen, writeSize)
		newOffset += writen
		offset = newOffset
		if err != nil {
			return errors.Wrap(err, "Failed to copy data from data generator")
		}

		if m.once == true {
			break
		}

		if i.NotModifyUntil.Value > 0 {
			if (size - newOffset) < i.NotModifyUntil.Value {
				break
			}
			newOffset += i.NotModifyUntil.Value
		}
	}
	return nil
}

func (m *modifyFilesWithIntervals) getFilesCount() (filesCount int64, err error) {
	filesCount = 0
	err = filepath.Walk(m.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filesCount++
		return nil
	})
	return
}

func (m *modifyFilesWithIntervals) Modify() error {
	log.Print("Files modification is started")
	log.Printf("write in the single thread: %t", m.generateInMultipleThreads == false)

	completeSignal := make(chan bool)
	errorChannel := make(chan error)
	filesProcessed := 0
	failed := false
	var totalFiles int64
	var fileSelector FileSelector
	var wg sync.WaitGroup

	if m.changeRatio < 1 {
		var err error
		totalFiles, err = m.getFilesCount()
		if err != nil {
			return errors.Wrap(err, "Failed to get files count")
		}
		fileSelector, err = CreateRundomFileSelector(m.gen, uint64(m.changeRatio*float64(totalFiles)), uint64(totalFiles))
		if err != nil {
			return errors.Wrap(err, "Failed to create rundom file selector")
		}
	} else {
		go func() {
			var err error
			totalFiles, err = m.getFilesCount()
			if err != nil {
				errorChannel <- errors.Wrap(err, "Failed to get files count")
			}
		}()
		wg.Add(1)
		fileSelector = CreateAllFilesSelector()
	}

	go func() {
		if m.generateInMultipleThreads == false {
			// Attach goroutine to the single thread to wrile files on disk with no fragmentation
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
		}

		err := filepath.Walk(m.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if failed {
				return fmt.Errorf("Canceled")
			}
			if info.IsDir() {
				return nil
			}
			r, err := fileSelector.IsFileIsSelected()
			if err != nil {
				return errors.Wrap(err, "Failed to check if file is selected to change")
			}
			if r == true {
				err = m.changeFile(path, info)
				filesProcessed++
			}
			return err
		})
		if failed {
			return
		}

		if err != nil {
			errorChannel <- errors.Wrap(err, "Failed to modify files")
			return
		}
		completeSignal <- true
	}()
	wg.Add(1)

	report := func() {
		if totalFiles == 0 {
			if m.printProgress {
				fmt.Printf("\rProcessed: %d        ", filesProcessed)
			}
		} else {
			fmt.Printf("\rProcessed: (%d/%d)        ", filesProcessed, totalFiles)
		}
	}

	timeout := time.Tick(time.Second)
	for {
		select {
		case <-timeout:
			report()
		case <-completeSignal:
			report()
			fmt.Println("")
			log.Print("Files modification has completed successfully")
			return nil
		case err := <-errorChannel:
			failed = true
			wg.Wait()
			return err
		}
	}
}

func CreateFilesModifierWithInterval(gen DataGenerator, path string, generateInMultipleThreads bool,
	changeRatio float64, interval Interval, once, reverse bool, quietMode bool) FilesModifier {
	return &modifyFilesWithIntervals{
		gen:  gen,
		path: path,
		generateInMultipleThreads: generateInMultipleThreads,
		printProgress:             quietMode == false,
		changeRatio:               changeRatio,
		interval:                  interval,
		once:                      once,
		reverse:                   reverse,
	}
}
