package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// RingBuffer is a circular buffer for storing [Data].
// It allows for writing and emitting data. When the
// buffer is full, the oldest data is overwritten.
type RingBuffer struct {
	data []*Data
	// total size of buffer
	size int
	// last element that was written to in buffer
	lastInsert int
	// next element to read during emit
	nextRead int
	// time between emit cycles
	emitTime time.Time
}

// Data represents input received from sensors.
type Data struct {
	Stamp time.Time
	Value string
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data: make([]*Data, size),
		size: size,
		// initialize to -1 so that we can discern when
		// no insert has occured yet.
		lastInsert: -1,
	}
}

// Insert adds a new [Data] to the [RingBuffer].
// If the buffer is full, the oldest data is overwritten.
func (r *RingBuffer) Insert(input Data) {
	r.lastInsert = (r.lastInsert + 1) % r.size
	r.data[r.lastInsert] = &input

	if r.nextRead == r.lastInsert {
		r.nextRead = (r.nextRead + 1) % r.size
	}
}

// Emit returns all data in [RingBuffer] since the last call
// to Emit.  If no data has been written since the last call
// to Emit, an empty slice is returned.
func (r *RingBuffer) Emit() []*Data {
	output := []*Data{}
	for {
		if r.data[r.nextRead] != nil {
			output = append(output, r.data[r.nextRead])
			r.data[r.nextRead] = nil
		}
		if r.nextRead == r.lastInsert || r.lastInsert == -1 {
			break
		}
		r.nextRead = (r.nextRead + 1) % r.size
	}
	return output
}

func main() {
	rb := NewRingBuffer(5)
	currentRune := 'a' - 1
	fmt.Println("EMPTY TEST:")
	spew.Dump(rb.Emit())
	fmt.Println("FULL TEST:")
	for i := 0; i < 10; i++ {
		currentRune++
		rb.Insert(Data{
			Stamp: time.Now(),
			Value: string(currentRune),
		})
	}
	spew.Dump(rb.Emit())
}
