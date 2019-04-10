package server

import (
	"context"
	"sync"
)

// The key used in oracle for locking user name.
type UserNameKey string

// Oracle is a data structure to solve complex concurrency problems.
type Oracle struct {
	sync.Mutex
	Map     map[interface{}]interface{}
	waiters []chan struct{}
}

// NewOracle creates an oracle.
func NewOracle() *Oracle {
	o := &Oracle{}
	o.Map = make(map[interface{}]interface{})
	return o
}

// Wait waits until either the context is cancelled or Broadcast() is called.
// If the context is cancelled, it returns the context's error and does NOT lock the oracle.
// Otherwise, it returns nil and locks the oracle.
// Required Mutex to be held.
func (o *Oracle) Wait(ctx context.Context) error {
	w := make(chan struct{}, 1)
	o.waiters = append(o.waiters, w)
	o.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w:
		o.Lock()
		return nil
	}
}

// Broadcast wakes up all waiting clients. Requires Mutex to be held.
func (o *Oracle) Broadcast() {
	for _, w := range o.waiters {
		select {
		case w <- struct{}{}:
		default:
		}
	}
	o.waiters = nil
}
