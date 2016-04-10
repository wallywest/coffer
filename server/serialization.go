package server

import (
	"encoding/json"
	"net/http"

	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
)

func writeResponseWithBody(w http.ResponseWriter, code int, resp interface{}) {
	enc, err := json.Marshal(resp)
	if err != nil {
		logger.Logger.Errorf("failed json-encoding: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err = w.Write(enc); err != nil {
		logger.Logger.Errorf("failed to write http response: %v", err)
	}
}
