package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// NewRouter make router with all apis
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", index)
	accountApis(router)
	return router
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	jsonRes(w, map[string]string{"hello": "world"})
}

func jsonRes(w http.ResponseWriter, res interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func errRes(w http.ResponseWriter, err error) {
	var status int
	var message string
	if httpError, ok := err.(HTTPError); ok {
		status = httpError.Status
		message = httpError.Message
	} else {
		status = 502
		message = err.Error()
	}
	res := map[string]string{"error": message}
	jsonRes(w, res)
	w.WriteHeader(errCode)
}

// HTTPError is a error with http status code
type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}
