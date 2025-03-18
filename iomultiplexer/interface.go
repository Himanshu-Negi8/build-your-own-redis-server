package iomultiplexer

import "time"

type IOMultiplexer interface {
	Subscribe(event Event) error
	Poll(timeout time.Duration) ([]Event, error)
	Close() error
}

type Event struct {
	Fd int
	Op Operations
}

const (
	// OpRead represents the read operation
	OpRead Operations = 1 << iota
	// OpWrite represents the write operation
	OpWrite
)

type Operations uint32
