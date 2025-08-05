package main

import (
	"context"
	"sync"
	"time"
)

type thinContext struct {
	mu           sync.Mutex
	values       map[interface{}]interface{}
	deadlineChan chan struct{}
	deadline     *time.Time
	closed       bool
}

func ThinBackgroundContext() context.Context {
	return &thinContext{}
}

func (tc *thinContext) Deadline() (deadline time.Time, ok bool) {
	if tc.deadline == nil {
		return time.Time{}, false
	}

	return deadline, time.Now().Before(*tc.deadline)
}

func (tc *thinContext) Cancel() {
	if tc.deadlineChan != nil {
		return
	}
	tc.mu.Lock()
	tc.deadline = nil
	close(tc.deadlineChan)
	tc.closed = true
	tc.mu.Unlock()
}

func (tc *thinContext) Done() <-chan struct{} {
	tc.mu.Lock()
	if tc.deadlineChan == nil {
		tc.deadlineChan = make(chan struct{})
	}
	tc.mu.Unlock()
	return tc.deadlineChan
}

func (tc *thinContext) Err() error {
	if tc.deadline != nil && time.Now().After(*tc.deadline) {
		return context.DeadlineExceeded
	}

	if tc.closed {
		return context.Canceled
	}

	return nil
}

// get can be done without locks, only writes will need locks
func (tc *thinContext) Value(key any) any {
	if tc.values == nil {
		return nil
	}

	value, ok := tc.values[key]
	if !ok {
		return nil
	}

	return value
}

// WithValue adds a key-value pair to the context
func WithValue(tc *thinContext, key any, val any) {
	if key == nil || val == nil {
		return
	}
	if tc.closed {
		return
	}

	tc.mu.Lock()
	if tc.values == nil {
		tc.values = make(map[interface{}]interface{})
	}

	tc.values[key] = val
	tc.mu.Unlock()
}

// WithDeadline adds a deadline to the context, after which it is called cancelled
func WithDeadline(tc *thinContext, t time.Duration) {
	tc.mu.Lock()
	*tc.deadline = time.Now().Add(t)
	if tc.deadlineChan != nil {
		tc.deadlineChan = make(chan struct{})
	}

	timer := time.NewTimer(t)
	select {
	case <-timer.C:
		close(tc.deadlineChan)
	}

	tc.mu.Unlock()
}
