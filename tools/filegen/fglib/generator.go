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
	"runtime"
	"sync"
	"time"
)

type DataGenerator interface {
	Seed(key []byte) error
	GetReader() (io.Reader, error)
}

/* Queues implementations */

type DataQueue interface {
	SetChanel(dest chan []byte, maxIndex int)
	Process(index int, data []byte)
}

/* Unordered queue implementation */

type UnorderedQueue struct {
	dest chan []byte
}

func (q *UnorderedQueue) SetChanel(dest chan []byte, maxIndex int) {
	q.dest = dest
}

func (q *UnorderedQueue) Process(index int, data []byte) {
	q.dest <- data
}

/* Ordered queue implementation */

type OrderedQueue struct {
	dest      chan []byte
	buffer    map[int][]byte
	guard     sync.Mutex
	nextIndex int
	maxIndex  int
}

func (q *OrderedQueue) SetChanel(dest chan []byte, maxIndex int) {
	q.dest = dest
	q.buffer = make(map[int][]byte)
	q.nextIndex = 0
	q.maxIndex = maxIndex
}

func (q *OrderedQueue) Process(index int, data []byte) {
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

/* Data generator implementation */

type Generator struct {
	generator DataGenerator
	data      chan []byte
	queue     DataQueue
	block     []byte
	complete  bool
	completed sync.WaitGroup
}

func generateRoutine(queue DataQueue, generator io.Reader, index int, complete *bool, completed *sync.WaitGroup) {
	defer completed.Done()
	blockSize := 1024 * 1024
	block := make([]byte, blockSize)
	processCompleted := make(chan bool)
	for {
		read, err := generator.Read(block)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: Failed to generate data with error", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		go func() {
			queue.Process(index, block[:read])
			processCompleted <- true
		}()

		processed := false
		for {
			if *complete == true {
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

func (gen *Generator) SetDataGenerator(generator DataGenerator, queue DataQueue) {
	gen.generator = generator
	gen.queue = queue
}

func (gen *Generator) Init() error {
	if gen.generator == nil {
		return fmt.Errorf("Generator is not set")
	}

	cpuCount := runtime.NumCPU()
	gen.data = make(chan []byte, cpuCount*2)
	gen.queue.SetChanel(gen.data, cpuCount)
	gen.block = nil
	gen.complete = false

	/* start generating goroutines */
	for i := 0; i < cpuCount; i++ {
		reader, err := gen.generator.GetReader()
		if err != nil {
			return err
		}
		gen.completed.Add(1)
		go generateRoutine(gen.queue, reader, i, &gen.complete, &gen.completed)
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
