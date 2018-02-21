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

	"github.com/pkg/errors"
)

/* Seed tools */

func SeedFromUint64(s uint64) []byte {
	seed := make([]byte, 16)
	binary.PutUvarint(seed, s)
	return seed
}

/* Crypto data generator implementation */

type cryptoGenerator struct { // inherits DataGenerator
}

func (gen *cryptoGenerator) Seed(key []byte) error {
	return ErrNotSupported
}

func (gen *cryptoGenerator) Read(block []byte) (int, error) {
	return rand.Read(block)
}

func (gen *cryptoGenerator) Close() error {
	return nil
}

func (gen *cryptoGenerator) Clone() (DataGenerator, error) {
	return &cryptoGenerator{}, nil
}

func CreateCryptoDataGenerator() DataGenerator {
	return &cryptoGenerator{}
}

/* Pseudo random data generator implementation */

type pseudoRandomGenerator struct {
	seed []byte

	block   []byte
	encrypt cipher.Block
	index   int
}

func (gen *pseudoRandomGenerator) init() error {
	var err error
	gen.encrypt, err = aes.NewCipher(gen.seed)
	if err != nil {
		return errors.Wrap(err, "Failed to create AES cipher")
	}
	gen.block = make([]byte, len(gen.seed))
	copy(gen.block, gen.seed)
	return nil
}

func (gen *pseudoRandomGenerator) Seed(key []byte) error {
	if len(key) != 16 {
		return fmt.Errorf("Seed must have 16 bytes length. Got: %d", len(key))
	}

	gen.seed = make([]byte, len(key))
	copy(gen.seed, key)
	return gen.init()
}

func (gen *pseudoRandomGenerator) Read(block []byte) (int, error) {
	bitSize := len(gen.block)
	bitsCount := len(block) / bitSize
	if len(block)%bitSize > 0 {
		bitsCount++
	}

	/* increment one ob bytes of encrypted block */
	encrypted := make([]byte, len(gen.block))
	for i := 0; i < bitsCount; i++ {
		gen.block[gen.index]++
		gen.index++
		if gen.index == bitSize {
			gen.index = 0
		}

		gen.encrypt.Decrypt(encrypted, gen.block)
		copy(block[bitSize*i:], encrypted)
	}
	return len(block), nil
}

func (gen *pseudoRandomGenerator) Close() error {
	return nil
}

func (gen *pseudoRandomGenerator) Clone() (DataGenerator, error) {
	clone := &pseudoRandomGenerator{}
	clone.Seed(gen.seed)
	gen.seed[0]++
	return clone, nil
}

func CreatePseudoRandomDataGenerator(seed []byte) (DataGenerator, error) {
	gen := &pseudoRandomGenerator{}
	err := gen.Seed(seed)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to set seed")
	}
	return gen, nil
}
