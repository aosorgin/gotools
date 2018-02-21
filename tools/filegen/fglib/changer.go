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
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/* Interval tools */

type intervalValueType struct {
	value    int64 // value
	obsolete bool  // if true, value stored in bytes otherwise in percents
}

// Intervals are present with 3 values [size not to modify; size to modify; size not to modify
type IntervalType struct {
	notModify      intervalValueType // not to modify
	modify         intervalValueType // modify
	notModifyUntil intervalValueType // not to modify
}

func GetFullInterval() IntervalType {
	return IntervalType{
		notModify:      intervalValueType{obsolete: true},
		modify:         intervalValueType{value: 100, obsolete: false},
		notModifyUntil: intervalValueType{obsolete: true},
	}
}

func GetEmptyInterval() IntervalType {
	return IntervalType{
		notModify:      intervalValueType{value: 100, obsolete: false},
		modify:         intervalValueType{obsolete: true},
		notModifyUntil: intervalValueType{obsolete: true},
	}
}

func getObsoleteInterval(interval IntervalType, size int64) IntervalType {
	result := interval
	makeObsolete := func(v *intervalValueType) {
		if v.obsolete == false {
			v.value = int64(float64(v.value*size) / float64(100))
			v.obsolete = true
		}
	}
	makeObsolete(&result.notModify)
	makeObsolete(&result.modify)
	makeObsolete(&result.notModifyUntil)
	return result
}

// Interval format [digit{,kK,mM,gG,%},*3]. First value to seek without modification.
// The next one is to modify file. The third is for seeking.

func ParseInterval(data string) (result IntervalType, err error) {
	intervals := strings.Split(data, ",")
	if len(intervals) < 2 || len(intervals) > 3 {
		err = fmt.Errorf("Invalid interval format. Must be 2 values at least but not more than 3")
		return
	}

	processIntervalValue := func(serialized string, i *intervalValueType) error {
		r := regexp.MustCompile("([\\d]*)([%kKmMgG]{0,1})").FindStringSubmatch(serialized)
		if len(r) < 3 {
			return fmt.Errorf("Invalid interval format for '%s'", serialized)
		}
		i.obsolete = true
		value, err := strconv.Atoi(r[1])
		if err != nil {
			return errors.Wrapf(err, "Failed to convert string '%s' to integer", r[1])
		}
		if value < 0 {
			return fmt.Errorf("Invalid interval format for '%s'. Must be positive", serialized)
		}
		i.value = int64(value)

		switch r[2] {
		case "%":
			i.obsolete = false
			if i.value > 100 {
				return fmt.Errorf("Invalid interval format for '%s'. Must be in [0;100]", serialized)
			}
		case "k":
			i.value *= 1000
		case "K":
			i.value *= 1024
		case "m":
			i.value *= 1000 * 1000
		case "M":
			i.value *= 1024 * 1024
		case "g":
			i.value *= 1000 * 1000 * 1000
		case "G":
			i.value *= 1024 * 1024 * 1024
		case "":
		default:
			panic(fmt.Errorf("Invalid postfix in regex"))
		}
		return nil
	}

	err = processIntervalValue(intervals[0], &result.notModify)
	if err != nil {
		return
	}

	err = processIntervalValue(intervals[1], &result.modify)
	if err != nil {
		return
	}

	if len(intervals) > 2 {
		err = processIntervalValue(intervals[2], &result.notModifyUntil)
		if err != nil {
			return
		}
	}

	return
}

func ParseSize(rawSize string) (uint64, error) {
	r := regexp.MustCompile("([\\d]*)([kKmMgG]{0,1})").FindStringSubmatch(rawSize)

	if len(r) < 2 {
		return uint64(0), fmt.Errorf("Invalid size format for '%s'", rawSize)
	}

	size, err := strconv.ParseInt(r[1], 10, 0)
	if err != nil {
		return uint64(0), errors.Wrapf(err, "Failed to parse integer '%s'", r[1])
	}

	switch r[2] {
	case "k":
		size *= 1000
	case "K":
		size *= 1024
	case "m":
		size *= 1000 * 1000
	case "M":
		size *= 1024 * 1024
	case "g":
		size *= 1000 * 1000 * 1000
	case "G":
		size *= 1024 * 1024 * 1024
	case "":
	default:
		panic(fmt.Errorf("Invalid postfix in regex"))
	}

	return uint64(size), nil
}

/* Get random subset in sequence mode */

type SequenceChecker interface {
	IsFileIsSelected() (bool, error)
}

type randomSubSet struct { // implements SequenceChecker
	gen      DataGenerator
	weight   float64
	hitScore float64
	length   uint64
	count    uint64
	index    uint64
	hits     uint64
}

func (r *randomSubSet) init(gen DataGenerator, count uint64, length uint64) error {
	if length == 0 {
		return fmt.Errorf("length argument cannot be 0")
	}

	if count > length {
		return fmt.Errorf("count(%d) argument cannot be more than length(%d)", count, length)
	}

	r.gen = gen
	r.count = count
	r.length = length
	r.weight = float64(length-1) / float64(length)
	r.hitScore = float64(length-uint64((float64(count)/1.4))) / float64(length)
	return nil
}

func (r *randomSubSet) IsFileIsSelected() (bool, error) {
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

func (r *randomSubSet) getFloat64() (float64, error) {
	var v uint64
	err := binary.Read(r.gen, binary.LittleEndian, &v)
	if err != nil {
		return float64(0), errors.Wrap(err, "Failed toread float value from data generator")
	}
	return (float64(v) / float64(^uint64(0))), nil
}

/* Simple checker that returns always true */

type alwaysTrueChecker struct {
}

func (r *alwaysTrueChecker) IsFileIsSelected() (bool, error) {
	return true, nil
}

/* Changer implementation */

type Changer struct {
	gen DataGenerator

	interval IntervalType
	once     bool // use once if true otherwise until file ending
	reverse  bool // use interval from the end of file if true
}

func (c *Changer) Init(gen DataGenerator, interval IntervalType, once bool, reverse bool) {
	c.gen = gen
	c.interval = interval
	c.once = once
	c.reverse = reverse
}

func (c *Changer) Close() error {
	return c.gen.Close()
}

func min(a int64, b int64) int64 {
	if a <= b {
		return a
	}
	return b
}

func changeFile(path string, info os.FileInfo, c *Changer) error {
	i := getObsoleteInterval(c.interval, info.Size())
	if i.modify.value == 0 {
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return errors.Wrapf(err, "Failed to open file '%s'", path)
	}
	defer file.Close()

	size := info.Size()
	write_size := int64(0)
	offset, new_offset := int64(0), int64(0)
	for offset < size {
		if i.notModify.value > 0 {
			if (size - new_offset) < i.notModify.value {
				break
			}
			new_offset += i.notModify.value
		}

		if c.reverse == true || new_offset != offset {
			if c.reverse == true {
				if (size - new_offset) < i.modify.value {
					_, err := file.Seek(0, io.SeekStart)
					if err != nil {
						return errors.Wrap(err, "Failed to seek")
					}
				} else {
					_, err := file.Seek(size-new_offset-i.modify.value, io.SeekStart)
					if err != nil {
						return errors.Wrap(err, "Failed to seek")
					}
				}
				write_size = min(size-new_offset, i.modify.value)
			} else {
				_, err := file.Seek(new_offset, io.SeekStart)
				if err != nil {
					return errors.Wrap(err, "Failed to seek")
				}
				write_size = min(i.modify.value, size-new_offset)
			}
			offset = new_offset
		} else {
			write_size = min(i.modify.value, size-new_offset)
		}

		writen, err := io.CopyN(file, c.gen, write_size)
		new_offset += writen
		offset = new_offset
		if err != nil {
			return errors.Wrap(err, "Failed to copy data from data generator")
		}

		if c.once == true {
			break
		}

		if i.notModifyUntil.value > 0 {
			if (size - new_offset) < i.notModifyUntil.value {
				break
			}
			new_offset += i.notModifyUntil.value
		}
	}
	return nil
}

func getFilesCount(path string) (filesCount int64, err error) {
	filesCount = 0
	filepath.Walk(Options.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		filesCount++
		return nil
	})
	return
}

func (c *Changer) ModifyFiles() error {
	completeSignal := make(chan bool)
	filesProcessed := 0
	var totalFiles int64
	var fileSelector SequenceChecker

	if Options.Change.Ratio < 1 {
		var err error
		totalFiles, err = getFilesCount(Options.Path)
		if err != nil {
			return errors.Wrap(err, "Failed to get files count")
		}
		r := new(randomSubSet)
		r.init(c.gen, uint64(Options.Change.Ratio*float64(totalFiles)), uint64(totalFiles))
		fileSelector = r
	} else {
		go func() {
			totalFiles, _ = getFilesCount(Options.Path)
		}()
		fileSelector = new(alwaysTrueChecker)
	}

	go func() {
		filepath.Walk(Options.Path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			r, err := fileSelector.IsFileIsSelected()
			if err != nil {
				return errors.Wrap(err, "Failed to check if file is selected to change")
			}
			if r == true {
				err = changeFile(path, info, c)
				filesProcessed++
			}
			return err
		})
		completeSignal <- true
	}()

	report := func() {
		if totalFiles == 0 {
			fmt.Printf("\rProcessed: %d        ", filesProcessed)
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
			return nil
		}
	}
	return nil
}
