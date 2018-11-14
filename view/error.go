package view

import (
	"fmt"
	"net/http"
)

func httpError(w http.ResponseWriter, statusCode int, a ...interface{}) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprint(a...)))
}

func httpErrorf(w http.ResponseWriter, statusCode int, format string, a ...interface{}) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf(format, a...)))
}
