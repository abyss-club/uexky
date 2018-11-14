package uexky

import (
	"context"
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

func reqWithValue(req *http.Request, key, val interface{}) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), key, val))
}
