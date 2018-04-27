package api

import (
	"context"

	"github.com/nanozuki/uexky/model"
)

// SliceInfoResolver ...
type SliceInfoResolver struct {
	SliceInfo *model.SliceInfo
}

// FirstCursor ...
func (si *SliceInfoResolver) FirstCursor(ctx context.Context) (string, error) {
	return si.SliceInfo.FirstCursor, nil
}

// LastCursor ...
func (si *SliceInfoResolver) LastCursor(ctx context.Context) (string, error) {
	return si.SliceInfo.LastCursor, nil
}
