package api

import "context"

// SliceInfoResolver ...
type SliceInfoResolver struct {
	firstCursor string
	lastCursor  string
}

// FirstCursor ...
func (si *SliceInfoResolver) FirstCursor(ctx context.Context) (string, error) {
	return si.firstCursor, nil
}

// LastCursor ...
func (si *SliceInfoResolver) LastCursor(ctx context.Context) (string, error) {
	return si.lastCursor, nil
}
