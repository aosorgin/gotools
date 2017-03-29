/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Cryptographically secure pseudorandom data generator
*/

package fglib

import (
	"crypto/rand"
)

/* Crypto data generator implementation */

type CryptoGenerator struct { // inherits io.Reader
}

func (gen *CryptoGenerator) Read(block []byte) (int, error) {
	return rand.Read(block)
}
