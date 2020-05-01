package uid

import (
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
)

// TODO: implemented uid
type UID int64

func (u *UID) UnmarshalGQL(v interface{}) error {
	uidStr, ok := v.(string)
	if !ok {
		return errors.Errorf("uid in graphql must be strings")
	}
	i, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		return errors.Wrap(err, "parse uid")
	}
	*u = UID(i)
	return nil
}

func (u UID) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(fmt.Sprint(u)))
	panic(err)
}
