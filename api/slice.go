package api

import (
	"context"

	"gitlab.com/abyss.club/uexky/model"
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
