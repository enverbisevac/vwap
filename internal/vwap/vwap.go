package vwap

import (
	"sync/atomic"
)

const defaultMaxSize = 200

// DataPoint represents a single data point from coinbase.
type DataPoint struct {
	Price  float64
	Volume float64
}

// Buffer represents a queue of DataPoints.
type Buffer struct {
	DataPoints        []DataPoint
	SumVolumeWeighted float64
	SumVolume         float64
	VWAP              float64

	MaxSize int32
	Reader  int32 // Reader position
	Writer  int32 // Writer Position
	filled  bool
}

// NewBuffer creates a new queue.
func NewBuffer(maxSize uint) *Buffer {
	if maxSize == 0 {
		maxSize = defaultMaxSize
	}

	return &Buffer{
		DataPoints: make([]DataPoint, defaultMaxSize),
		MaxSize:    int32(maxSize),
	}
}

// Len returns the length of the queue.
func (r *Buffer) Len() int {
	return len(r.DataPoints)
}

// Push pushes an element onto the queue, drops the first one when MaxSize is reached.
func (r *Buffer) Push(d DataPoint) {
	current := atomic.LoadInt32(&r.Writer)

	if current == r.MaxSize || r.filled {
		d := r.DataPoints[atomic.LoadInt32(&r.Writer)]
		r.filled = true
		// remove oldest value
		r.SumVolumeWeighted = r.SumVolumeWeighted - d.Price*d.Volume
		r.SumVolume = r.SumVolume - d.Volume
		if r.SumVolume > 0 {
			r.VWAP = r.SumVolumeWeighted / r.SumVolume
		}
	}

	if r.VWAP > 0 {
		r.SumVolumeWeighted += d.Price * d.Volume
		r.SumVolume += d.Volume
		r.VWAP = r.SumVolumeWeighted / r.SumVolume
	} else {
		initialVW := d.Price * d.Volume
		r.SumVolumeWeighted = initialVW
		r.SumVolume = d.Volume
		r.VWAP = initialVW / d.Volume
	}

	r.DataPoints[current] = d
	next := (current + 1) % r.MaxSize
	atomic.StoreInt32(&r.Writer, next)
}
