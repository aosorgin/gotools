/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Common data generator interfaces and tools
*/

package fglib

import (
	"fmt"
	"io"
	"os"
	//"runtime"
	"sync"
	"time"
)

type DataGenerator interface {
	Seed(key []byte) error
	Read(p []byte) (int, error)
}

type Generator struct {
	generator io.Reader
	data      chan []byte
	block     []byte
	complete  bool
	completed sync.WaitGroup
}

func generateRoutine(data chan []byte, generator io.Reader, complete *bool, completed *sync.WaitGroup) {
	defer completed.Done()
	blockSize := 1024 * 1024
	block := make([]byte, blockSize)

	for {
		read, err := generator.Read(block)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: Failed to generate data with error", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		for {
			if *complete == true {
				return
			}

			select {
			case data <- block[:read]:
				break // Go to generate of new block
			case <-time.After(time.Second):
				// Check if coplete flag is set
			}
		}
	}
}

func (gen *Generator) SetDataGenerator(generator io.Reader) {
	gen.generator = generator
}

func (gen *Generator) Init() error {
	if gen.generator == nil {
		return fmt.Errorf("Generator is not set")
	}

	// TODO: Use multimpe goroutines to generate pseudo random data (see PseudoRandomGenerator)
	// Issues:
	//   -- Block from different goroulines must be ordered with goroutine index to order the data flow
	//   -- Make its own encrypting block instance per goroutine. Could give DataGenerator factory instead it own

	cpuCount := /*runtime.NumCPU()*/ 1
	gen.data = make(chan []byte, cpuCount*2)
	gen.block = nil
	gen.complete = false
	for i := 0; i < cpuCount; i++ {
		gen.completed.Add(1)
		go generateRoutine(gen.data, gen.generator, &gen.complete, &gen.completed)
	}
	return nil
}

func (gen *Generator) Stop() error {
	gen.complete = true
	gen.completed.Wait()
	return nil
}

func (gen *Generator) Read(p []byte) (int, error) {
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
