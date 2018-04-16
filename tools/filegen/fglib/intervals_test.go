/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Tests for intervals parsing
*/

package fglib

import (
	"fmt"
	"testing"
)

func checkIntervalValue(t *testing.T, valueName string, i IntervalValue, value int64, obsolete bool) {
	if i.Value != value || i.Obsolete != obsolete {
		t.Error(fmt.Errorf("Failed for process '%s' with values: (%d, %t). Must be (%d, %t)",
			valueName, i.Value, i.Obsolete, value, obsolete))
	}
}

func TestParseObsoleteIntervalShort(t *testing.T) {
	i, err := ParseInterval("1234,654")
	if err != nil {
		t.Error(err)
	}
	checkIntervalValue(t, "notModify", i.NotModify, 1234, true)
	checkIntervalValue(t, "modify", i.Modify, 654, true)
}

func TestParseObsoleteInterval(t *testing.T) {
	i, err := ParseInterval("8953,480,30902")
	if err != nil {
		t.Error(err)
	}
	checkIntervalValue(t, "notModify", i.NotModify, 8953, true)
	checkIntervalValue(t, "modify", i.Modify, 480, true)
	checkIntervalValue(t, "notModifyUntil", i.NotModifyUntil, 30902, true)
}

func TestFailOnNegativeInterval(t *testing.T) {
	_, err := ParseInterval("-8953,480,30902")
	if err == nil {
		t.Error(fmt.Errorf("Successfully parse negative interval"))
	}

	_, err = ParseInterval("8953,-480,30902")
	if err == nil {
		t.Error(fmt.Errorf("Successfully parse negative interval"))
	}

	_, err = ParseInterval("8953,480,-30902")
	if err == nil {
		t.Error(fmt.Errorf("Successfully parse negative interval"))
	}
}

func TestParseNotObsoleteInterval(t *testing.T) {
	i, err := ParseInterval("23%%,10%%,17%%")
	if err != nil {
		t.Error(err)
	}
	checkIntervalValue(t, "notModify", i.NotModify, 23, false)
	checkIntervalValue(t, "modify", i.Modify, 10, false)
	checkIntervalValue(t, "notModifyUntil", i.NotModifyUntil, 17, false)
}

func TestFailToParseNotObsoleteIntervalWithInvalidValues(t *testing.T) {
	i, err := ParseInterval("8953%%,80%%,3%%")
	if err != nil {
		t.Error(err)
	}

	if _, err = i.GetObsolete(int64(100), true); err == nil {
		t.Error(fmt.Errorf("Successfully get obsolete with more then 100%% interval"))
	}

	if _, err = i.GetObsolete(int64(100), false); err != nil {
		t.Error(err)
	}

	i, err = ParseInterval("89%%,840%%,3%%")
	if err != nil {
		t.Error(err)
	}

	if _, err = i.GetObsolete(int64(100), true); err == nil {
		t.Error(fmt.Errorf("Successfully get obsolete with more then 100%% interval"))
	}

	if _, err = i.GetObsolete(int64(100), false); err != nil {
		t.Error(err)
	}

	i, err = ParseInterval("89%%,40%%,378%%")
	if err != nil {
		t.Error(err)
	}

	if _, err = i.GetObsolete(int64(100), true); err == nil {
		t.Error(fmt.Errorf("Successfully get obsolete with more then 100%% interval"))
	}

	if _, err = i.GetObsolete(int64(100), false); err != nil {
		t.Error(err)
	}
}

func TestParseCombineInterval(t *testing.T) {
	i, err := ParseInterval("23K,96%%,17k")
	if err != nil {
		t.Error(err)
	}
	checkIntervalValue(t, "notModify", i.NotModify, 23*1024, true)
	checkIntervalValue(t, "modify", i.Modify, 96, false)
	checkIntervalValue(t, "notModifyUntil", i.NotModifyUntil, 17*1000, true)

	i, err = ParseInterval("23%%,96m,17M")
	if err != nil {
		t.Error(err)
	}
	checkIntervalValue(t, "notModify", i.NotModify, 23, false)
	checkIntervalValue(t, "modify", i.Modify, 96*1000*1000, true)
	checkIntervalValue(t, "notModifyUntil", i.NotModifyUntil, 17*1024*1024, true)

	i, err = ParseInterval("23g,96G,17%%")
	if err != nil {
		t.Error(err)
	}
	checkIntervalValue(t, "notModify", i.NotModify, 23*1000*1000*1000, true)
	checkIntervalValue(t, "modify", i.Modify, 96*1024*1024*1024, true)
	checkIntervalValue(t, "notModifyUntil", i.NotModifyUntil, 17, false)
}

func TestGetFullInterval(t *testing.T) {
	i := GetFullInterval()
	checkIntervalValue(t, "notModify", i.NotModify, 0, true)
	checkIntervalValue(t, "modify", i.Modify, 100, false)
	checkIntervalValue(t, "notModifyUntil", i.NotModifyUntil, 0, true)
}

func TestGetEmptyInterval(t *testing.T) {
	i := GetEmptyInterval()
	checkIntervalValue(t, "notModify", i.NotModify, 100, false)
	checkIntervalValue(t, "modify", i.Modify, 0, true)
	checkIntervalValue(t, "notModifyUntil", i.NotModifyUntil, 0, true)
}
