package repo

import (
	"fmt"

	"gitlab.com/abyss.club/uexky/lib/uerr"
)

func dbErrWrap(err error, a ...interface{}) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprint(a...)
	return uerr.Errorf(uerr.DBError, "%s: %w", msg, err)
}

// func dbErrWrapf(err error, format string, a ...interface{}) error {
// 	if err == nil {
// 		return nil
// 	}
// 	msg := fmt.Sprintf(format, a...)
// 	return uerr.Errorf(uerr.DBError, "%s: %w", msg, err)
// }
