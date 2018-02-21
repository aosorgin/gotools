/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Static (not random) data generators implementations
*/

package fglib

/* Null data (consists on null values) generator implementation */

type nullGenerator struct { // inherits DataGenerator
}

func (gen *nullGenerator) Seed(key []byte) error {
	return ErrNotSupported
}

func (gen *nullGenerator) Read(block []byte) (int, error) {
	for i := range block {
		block[i] = byte(0)
	}
	return len(block), nil
}

func (gen *nullGenerator) Close() error {
	return nil
}

func (gen *nullGenerator) Clone() (DataGenerator, error) {
	return &nullGenerator{}, nil
}

func CreateNullDataGenerator() DataGenerator {
	return &nullGenerator{}
}
