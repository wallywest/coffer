package server

import (
	"net/http"

	"gitlab.vailsys.com/jerny/coffer/internal/recording"
)

const (
	serverError         = "server_error"
	invalidRequestError = "invalid_request"
	accessDeniedErrror  = "access_denied"
)

type apiError struct {
	Type        string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

func (e *apiError) Error() string {
	return e.Type
}

func newAPIError(typ, desc string) *apiError {
	return &apiError{Type: typ, Description: desc}
}

func writeAPIError(w http.ResponseWriter, code int, err error) {
	apierr, ok := err.(*apiError)
	if ok {
		writeResponseWithBody(w, code, apierr)
		return
	}

	repError, ok := err.(recording.RepoError)
	if ok {
		apierr := newAPIError(repError.Error(), "")
		writeResponseWithBody(w, code, apierr)
		return
	}

	apierr = newAPIError(serverError, "")
	if apierr.Type == "" {
		apierr.Type = serverError
	}

	if code == 0 {
		code = http.StatusInternalServerError
	}

	writeResponseWithBody(w, code, apierr)
}
