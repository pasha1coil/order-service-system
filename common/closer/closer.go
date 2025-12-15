package closer

import (
	"context"
)

type Closer interface {
	Close(ctx context.Context) error
}

type CloserFunc func(ctx context.Context) error

func (cf CloserFunc) Close(ctx context.Context) error {
	return cf(ctx)
}

type CloserGroup struct {
	closers []Closer
}

func NewCloserGroup() *CloserGroup {
	return &CloserGroup{}
}

func (cg *CloserGroup) Add(c Closer) {
	cg.closers = append(cg.closers, c)
}

func (cg *CloserGroup) Call(ctx context.Context) error {
	var closeErr error
	for i := len(cg.closers) - 1; i >= 0; i-- {
		if err := cg.closers[i].Close(ctx); err != nil && closeErr == nil {
			closeErr = err
		}
	}
	return closeErr
}
