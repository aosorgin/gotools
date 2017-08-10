/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Static (not random) data generators implementations
*/

package fglib

import (
	"fmt"
	"io"
)

/* Null data (consists on null values) generator implementation */

type nullGeneratorReader struct {
}

func (gen *nullGeneratorReader) Read(block []byte) (int, error) {
	for i, _ := range block {
		block[i] = byte(0)
	}
	return len(block), nil
}

type NullGenerator struct { // inherits DataGenerator
}

func (gen *NullGenerator) Seed(key []byte) error {
	return fmt.Errorf("Seed is not supported for null generator")
}

func (gen *NullGenerator) GetReader() (io.Reader, error) {
	return new(nullGeneratorReader), nil
}
