package uexky

import (
	"context"

	"gitlab.com/abyss.club/uexky/lib/uerr"
)

type Limiter struct {
	Limit int
	Count int
}

type contextKey int

const (
	limiterKey contextKey = 1 + iota
)

func AttachLimiter(ctx context.Context, limit int) context.Context {
	return context.WithValue(ctx, limiterKey, &Limiter{Limit: limit})
}

func Cost(ctx context.Context, cost int) error {
	v := ctx.Value(limiterKey)
	limiter, ok := v.(*Limiter)
	if !ok || limiter == nil {
		return nil
	}
	limiter.Count += cost
	if limiter.Count > limiter.Limit {
		return uerr.Errorf(
			uerr.ComplexityError,
			"operation has complexity %v at least, which exceeds the limit of %v",
			limiter.Count, limiter.Limit,
		)
	}
	return nil
}
