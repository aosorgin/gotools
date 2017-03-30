/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Cryptographically secure pseudorandom data generator
*/

package fglib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
)

/* Seed tools */

func SeedFromUint64(s uint64) []byte {
	seed := make([]byte, 16)
	binary.PutUvarint(seed, s)
	return seed
}

/* Crypto data generator implementation */

type cryptoGeneratorReader struct {
}

func (gen *cryptoGeneratorReader) Read(block []byte) (int, error) {
	return rand.Read(block)
}

type CryptoGenerator struct { // inherits DataGenerator
}

func (gen *CryptoGenerator) Seed(key []byte) error {
	return fmt.Errorf("Seed is not supported for crypto random generator")
}

func (gen *CryptoGenerator) GetReader() (io.Reader, error) {
	return new(cryptoGeneratorReader), nil
}

/* Pseudo random data generator implementation */

type pseudoRandomDataReader struct {
	block   []byte
	encrypt cipher.Block
	index   int
}

func (gen *pseudoRandomDataReader) Read(block []byte) (int, error) {
	bitSize := len(gen.block)
	bitsCount := len(block) / bitSize
	if len(block)%bitSize > 0 {
		bitsCount++
	}

	/* increment one ob bytes of encrypted block */
	encrypted := make([]byte, len(gen.block))
	for i := 0; i < bitsCount; i++ {
		gen.block[gen.index] += 1
		gen.index += 1
		if gen.index == bitSize {
			gen.index = 0
		}

		gen.encrypt.Decrypt(encrypted, gen.block)
		copy(block[bitSize*i:], encrypted)
	}
	return len(block), nil
}

type PseudoRandomGenerator struct {
	seed []byte
}

func (gen *PseudoRandomGenerator) Seed(key []byte) error {
	if len(key) != 16 {
		return fmt.Errorf("Seed must have 16 bytes length. Got: %d", len(key))
	}

	gen.seed = key
	return nil
}

func (gen *PseudoRandomGenerator) GetReader() (io.Reader, error) {
	reader := new(pseudoRandomDataReader)
	var err error
	reader.encrypt, err = aes.NewCipher(gen.seed)
	if err != nil {
		return nil, err
	}
	reader.block = make([]byte, len(gen.seed))
	copy(reader.block, gen.seed)
	gen.seed[0]++
	return reader, nil
}
