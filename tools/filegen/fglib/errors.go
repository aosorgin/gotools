/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017-2018

Brief:     Errors declaration
*/

package fglib

import (
	"fmt"
)

var (
	ErrNotSupported = fmt.Errorf("Not supported")
)
