package repo

import (
	"errors"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type AppleWhereAndOrderFunc func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error)
type DealSliceResultFunc func(i int, isFirst bool, isLast bool)

func applySliceQuery(fn AppleWhereAndOrderFunc, q *orm.Query, sq *entity.SliceQuery) (*orm.Query, error) {
	if (sq.After == nil && sq.Before == nil) || (sq.After != nil && sq.Before != nil) {
		return nil, errors.New("one and only one of before or after must be specified")
	}
	if sq.Limit == 0 {
		return nil, errors.New("limit must be specified")
	}
	nq := q.Limit(sq.Limit + 1) // limit+1 for "hasNext"
	if sq.After != nil {
		return fn(nq, true, *sq.After)
	}
	return fn(nq, true, *sq.Before)
}

func dealSliceResult(fn DealSliceResultFunc, sq *entity.SliceQuery, length int, reverse bool) {
	rstLen := sq.Limit
	if length < sq.Limit {
		rstLen = length
	}
	first := 0
	end := rstLen
	step := 1
	if reverse {
		first = length - 1
		end = first - rstLen
		step = -1
	}
	for i := first; i != end; i += step {
		fn(i, i == first, i == end-step)
	}
}