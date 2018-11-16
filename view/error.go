package view

import (
	"fmt"
	"net/http"
)

func httpError(w http.ResponseWriter, code int, a ...interface{}) {
	err := fmt.Sprint(a...)
	http.Error(w, err, code)
}

func httpErrorf(w http.ResponseWriter, code int, format string, a ...interface{}) {
	err := fmt.Sprintf(format, a...)
	http.Error(w, err, code)
}
