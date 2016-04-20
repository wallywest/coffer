package server

import (
	"net/http"
	"time"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/logger"

	"github.com/codegangsta/negroni"
)

func loggerMiddleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		header := w.Header()
		requestId := header.Get("X-Pulse-Request-Id")

		path := r.URL.Path
		logger.Logger.Infof("requestid=%s method=%s path=%s agent=%s host=%s request=%s", requestId, r.Method, r.URL.Path, r.UserAgent(), r.Host, r.RequestURI)

		start := time.Now()

		next(w, r)

		end := time.Now()
		latency := end.Sub(start)
		res := w.(negroni.ResponseWriter)

		logger.Logger.Infof("requestid=%s status=%d latency=%v method=%s path=%s", requestId, res.Status(), latency, r.Method, path)
	}
}
