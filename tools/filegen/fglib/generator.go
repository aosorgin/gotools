/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Common data generator interfaces and tools
*/

package fglib

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type DataGenerator interface {
	io.ReadCloser
	Clone() (DataGenerator, error)
	Seed(key []byte) error
}

/* Queues implementations */

type DataQueue interface {
	SetDestination(dest chan []byte, density int)
	ProcessBlock(index int, data []byte)
}

/* Unordered queue implementation */

type unorderedQueue struct {
	dest chan []byte
}

func (q *unorderedQueue) SetDestination(dest chan []byte, density int) {
	q.dest = dest
}

func (q *unorderedQueue) ProcessBlock(index int, data []byte) {
	q.dest <- data
}

func CreateUnorderedQueue() DataQueue {
	return &unorderedQueue{}
}

/* Ordered queue implementation */

type orderedQueue struct {
	dest      chan []byte
	buffer    map[int][]byte
	guard     sync.Mutex
	nextIndex int
	maxIndex  int
}

func (q *orderedQueue) SetDestination(dest chan []byte, density int) {
	q.dest = dest
	q.buffer = make(map[int][]byte)
	q.nextIndex = 0
	q.maxIndex = density
}

func (q *orderedQueue) ProcessBlock(index int, data []byte) {
	q.guard.Lock()
	defer q.guard.Unlock()
	if index != q.nextIndex {
		q.buffer[index] = data
		return
	}

	q.dest <- data
	nextIndex := func() {
		q.nextIndex = (q.nextIndex + 1) % q.maxIndex
	}
	nextIndex()
	for {
		var ok bool
		data, ok = q.buffer[q.nextIndex]
		if ok == false {
			return
		}
		q.dest <- data
		delete(q.buffer, q.nextIndex)
		nextIndex()
	}
}

func CreateOrderedQueue() DataQueue {
	return &orderedQueue{}
}

/* Data generator implementation */

type mutiThreadGenerator struct {
	generator DataGenerator
	data      chan []byte
	queue     DataQueue
	block     []byte
	stopping  bool
	completed sync.WaitGroup
}

func generateRoutine(queue DataQueue, generator io.ReadCloser, index int, stopping *bool, completed *sync.WaitGroup) {
	defer completed.Done()
	defer generator.Close()
	blockSize := 1024 * 1024
	block := make([]byte, blockSize)
	processCompleted := make(chan bool)
	for {
		read, err := generator.Read(block)
		if err != nil {
			log.Print(errors.Wrap(err, "Failed to generate data"))
			time.Sleep(100 * time.Millisecond)
			continue
		}

		go func() {
			queue.ProcessBlock(index, block[:read])
			processCompleted <- true
		}()

		processed := false
		for {
			if *stopping == true {
				return
			}

			if processed == true {
				break
			}

			select {
			case <-processCompleted:
				processed = true
				break // Go to generate of new block
			case <-time.After(time.Second):
				// Check if complete flag is set
			}
		}
	}
}

func (gen *mutiThreadGenerator) init() error {
	if gen.generator == nil {
		return fmt.Errorf("Generator is not set")
	}

	cpuCount := runtime.NumCPU()
	gen.data = make(chan []byte, cpuCount*2)
	gen.queue.SetDestination(gen.data, cpuCount)
	gen.block = nil
	gen.stopping = false

	/* start generating goroutines */
	for i := 0; i < cpuCount; i++ {
		reader, err := gen.generator.Clone()
		if err != nil {
			return errors.Wrap(err, "Failed to clone generator")
		}
		gen.completed.Add(1)
		go generateRoutine(gen.queue, reader, i, &gen.stopping, &gen.completed)
	}
	return nil
}

func (gen *mutiThreadGenerator) Seed(seed []byte) error {
	return ErrNotSupported
}

func (gen *mutiThreadGenerator) Close() error {
	gen.stopping = true
	gen.completed.Wait()
	return nil
}

func (gen *mutiThreadGenerator) Read(p []byte) (int, error) {
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

func (gen *mutiThreadGenerator) Clone() (DataGenerator, error) {
	return nil, ErrNotSupported
}

func CreateMutliThreadGenerator(generator DataGenerator, queue DataQueue) (DataGenerator, error) {
	gen := &mutiThreadGenerator{
		generator: generator,
		queue:     queue,
	}
	err := gen.init()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create multi-thread data generator")
	}
	return gen, nil
}
