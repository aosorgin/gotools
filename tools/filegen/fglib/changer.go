/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Files changer for random generator
*/

package fglib

/* Interval tools */

type intervalValueType struct {
	value    uint64 // value
	obsolete bool   // if true, value stored in bytes otherwise in percents
}

// Intervals are present with 3 values [size not to modify; size to modify; size not to modify
type IntervalType struct {
	notModify      intervalValueType // not to modify
	modify         intervalValueType // modify
	notModifyUntil intervalValueType // not to modify
}

func GetFullInterval() IntervalType {
	return IntervalType{
		modify: intervalValueType{value: 100, obsolete: false},
	}
}

func GetEmptyInterval() IntervalType {
	return IntervalType{
		notModify: intervalValueType{value: 100, obsolete: false},
	}
}

func getObsoleteInterval(interval IntervalType, size uint64) IntervalType {
	result := interval
	correctProc := func(v intervalValueType) {
		if v.obsolete == false {
			v.value = uint64(float64(v.value*size) / float64(100))
			v.obsolete = true
		}
	}
	correctProc(result.notModify)
	correctProc(result.modify)
	correctProc(result.notModifyUntil)
	return result
}

/* Changer implementation */

type Changer struct {
	gen *Generator // data generator

	interval IntervalType
	once     bool // use once if true otherwise until file ending
	reverse  bool // use interval from the end of file if true
}

func (c *Changer) Init(gen *Generator, interval IntervalType, once bool, reverse bool) error {
	c.gen = gen
	c.interval = interval
	c.once = once
	c.reverse = reverse
	return c.gen.Init()
}

func (c *Changer) Close() error {
	return c.gen.Stop()
}

func (c *Changer) ModifyFiles() {
	// TODO: implement
}
