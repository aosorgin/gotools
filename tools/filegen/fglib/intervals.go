/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Intervals implementation
*/

package fglib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

/* Interval tools */

type IntervalValue struct {
	Value    int64 // value
	Obsolete bool  // if true, value stored in bytes otherwise in percents
}

// Intervals are present with 3 values [size not to modify; size to modify; size not to modify
type Interval struct {
	NotModify      IntervalValue // not to modify
	Modify         IntervalValue // modify
	NotModifyUntil IntervalValue // not to modify
}

func GetFullInterval() Interval {
	return Interval{
		NotModify:      IntervalValue{Obsolete: true},
		Modify:         IntervalValue{Value: 100, Obsolete: false},
		NotModifyUntil: IntervalValue{Obsolete: true},
	}
}

func GetEmptyInterval() Interval {
	return Interval{
		NotModify:      IntervalValue{Value: 100, Obsolete: false},
		Modify:         IntervalValue{Obsolete: true},
		NotModifyUntil: IntervalValue{Obsolete: true},
	}
}

func GetObsoleteInterval(interval Interval, size int64) Interval {
	result := interval
	makeObsolete := func(v *IntervalValue) {
		if v.Obsolete == false {
			v.Value = int64(float64(v.Value*size) / float64(100))
			v.Obsolete = true
		}
	}
	makeObsolete(&result.NotModify)
	makeObsolete(&result.Modify)
	makeObsolete(&result.NotModifyUntil)
	return result
}

// Interval format [digit{,kK,mM,gG,%},*3]. First value to seek without modification.
// The next one is to modify file. The third is for seeking.

func ParseInterval(data string) (result Interval, err error) {
	intervals := strings.Split(data, ",")
	if len(intervals) < 2 || len(intervals) > 3 {
		err = fmt.Errorf("Invalid interval format. Must be 2 values at least but not more than 3")
		return
	}

	processIntervalValue := func(serialized string, i *IntervalValue) error {
		r := regexp.MustCompile("([\\d]*)([%kKmMgG]{0,1})").FindStringSubmatch(serialized)
		if len(r) < 3 {
			return fmt.Errorf("Invalid interval format for '%s'", serialized)
		}
		i.Obsolete = true
		value, err := strconv.Atoi(r[1])
		if err != nil {
			return errors.Wrapf(err, "Failed to convert string '%s' to integer", r[1])
		}
		if value < 0 {
			return fmt.Errorf("Invalid interval format for '%s'. Must be positive", serialized)
		}
		i.Value = int64(value)

		switch r[2] {
		case "%":
			i.Obsolete = false
			if i.Value > 100 {
				return fmt.Errorf("Invalid interval format for '%s'. Must be in [0;100]", serialized)
			}
		case "k":
			i.Value *= 1000
		case "K":
			i.Value *= 1024
		case "m":
			i.Value *= 1000 * 1000
		case "M":
			i.Value *= 1024 * 1024
		case "g":
			i.Value *= 1000 * 1000 * 1000
		case "G":
			i.Value *= 1024 * 1024 * 1024
		case "":
		default:
			panic(fmt.Errorf("Invalid postfix in regex"))
		}
		return nil
	}

	err = processIntervalValue(intervals[0], &result.NotModify)
	if err != nil {
		return
	}

	err = processIntervalValue(intervals[1], &result.Modify)
	if err != nil {
		return
	}

	if len(intervals) > 2 {
		err = processIntervalValue(intervals[2], &result.NotModifyUntil)
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
