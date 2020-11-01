package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/errors"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errors.BadParams):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, errors.NoAuth):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, errors.Permission):
		w.WriteHeader(http.StatusForbidden)
	case errors.Is(err, errors.NotFound):
		w.WriteHeader(http.StatusNotFound)
	default: // InternalError and other errors
		w.WriteHeader(http.StatusInternalServerError)
	}
	if _, err := w.Write([]byte(err.Error())); err != nil {
		log.Error(errors.Internal.Handle(err, "write http error"))
	}
}
