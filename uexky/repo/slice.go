package repo

import (
	"fmt"

	"github.com/go-pg/pg/v9/orm"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type sliceHelper struct {
	Column      string
	Desc        bool
	TransCursor func(string) (interface{}, error)
	SQ          *entity.SliceQuery
}

func (h *sliceHelper) Select(q *orm.Query) error {
	if (h.SQ.After == nil) == (h.SQ.Before == nil) {
		return uerr.New(uerr.ParamsError, "one and only one of before or after must be specified")
	}
	if h.SQ.Limit == 0 {
		return uerr.New(uerr.ParamsError, "limit must be specified")
	}

	var cursor string
	if h.SQ.After != nil {
		cursor = *h.SQ.After
	} else {
		cursor = *h.SQ.Before
	}
	if cursor != "" {
		value, err := h.TransCursor(cursor)
		if err != nil {
			return errors.Wrap(err, "invalid cursor")
		}
		var op string
		if (h.SQ.After != nil) == h.Desc {
			op = "<"
		} else {
			op = ">"
		}
		q = q.Where(fmt.Sprintf("%s %s ?", h.Column, op), value)
	}
	order := h.Column
	if (h.SQ.After != nil) == h.Desc {
		order += " DESC"
	}
	return q.Order(order).Limit(h.SQ.Limit + 1).Select()
}

func (h *sliceHelper) DealResults(length int, fn func(i int)) {
	rstLen := h.SQ.Limit
	if length < h.SQ.Limit {
		rstLen = length
	}
	first := 0
	end := rstLen
	step := 1
	if h.SQ.Before != nil { // reverse result
		first = rstLen - 1
		end = -1
		step = -1
	}
	for i := first; i != end; i += step {
		fn(i)
	}
}
