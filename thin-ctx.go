package main

import (
	"context"
	"time"
)

type thinContext struct{}

func ThinBackgroundContext() context.Context {
	return &thinContext{}
}

func (tc *thinContext) Deadline() (deadline time.Time, ok bool) {
	// TODO
	return deadline, ok
}

func (tc *thinContext) Done() <-chan struct{} {
	// TODO
	return nil
}

func (tc *thinContext) Err() error {
	// TODO
	return nil
}

func (tc *thinContext) Value(key any) any {
	// TODO
	return nil
}
