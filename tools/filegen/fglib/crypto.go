/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Cryptographically secure pseudorandom data generator
*/

package fglib

import (
	"crypto/rand"
	"runtime"
	"sync"
	"time"
)

func generateRoutine(data chan []byte, complete *bool, completed *sync.WaitGroup) {
	defer completed.Done()
	blockSize := 1024 * 1024
	block := make([]byte, blockSize)

	for {
		rand.Read(block)
		for {
			if *complete == true {
				return
			}

			select {
			case data <- block:
				break // Go to generate of new block
			case <-time.After(time.Second):
				// Check if coplete flag is set
			}
		}
	}
}

type CryptoGenerator struct {
	data      chan []byte
	block     []byte
	complete  bool
	completed sync.WaitGroup
}

func (gen *CryptoGenerator) Init() {
	cpuCount := runtime.NumCPU()
	gen.data = make(chan []byte, cpuCount*2)
	gen.block = nil
	gen.complete = false
	for i := 0; i < cpuCount; i++ {
		gen.completed.Add(1)
		go generateRoutine(gen.data, &gen.complete, &gen.completed)
	}
}

func (gen *CryptoGenerator) Stop() {
	gen.complete = true
	gen.completed.Wait()
}

func (gen *CryptoGenerator) Read(p []byte) (int, error) {
	offset := 0
	size := len(p)
	for size > 0 {
		if gen.block == nil {
			gen.block = <-gen.data
		}

		blockSize := len(gen.block)
		if blockSize <= size {
			copy(p[offset:], gen.block)
			offset += blockSize
			size -= blockSize
			gen.block = nil
			continue
		}

		copy(p[offset:], gen.block[:size])
		gen.block = gen.block[size:]
		break
	}

	return size, nil
}
