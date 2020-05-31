package server

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.com/abyss.club/uexky/lib/uerr"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, uerr.New(uerr.ParamsError)):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, uerr.New(uerr.AuthError)):
		w.WriteHeader(http.StatusUnauthorized)
	case errors.Is(err, uerr.New(uerr.PermissionError)):
		w.WriteHeader(http.StatusForbidden)
	case errors.Is(err, uerr.New(uerr.NotFoundError)):
		w.WriteHeader(http.StatusNotFound)
	default: // InternalError and other errors
		w.WriteHeader(http.StatusInternalServerError)
	}
	if _, err := w.Write([]byte(err.Error())); err != nil {
		log.Error(uerr.Errorf(uerr.InternalError, "write http error: %w", err))
	}
}
