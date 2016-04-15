package server

import (
	"net/http"
	"time"

	"github.com/codegangsta/negroni"

	"gitlab.vailsys.com/jerny/coffer/internal/logger"
)

func loggerMiddleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		header := w.Header()
		requestId := header.Get("X-Pulse-Request-Id")

		path := r.URL.Path
		logger.Logger.Infof("requestid=%s method=%s path=%s agent=%s host=%s request=%s", requestId, r.Method, r.URL.Path, r.UserAgent(), r.Host, r.RequestURI)

		start := time.Now()

		next(w, r)
		//h.ServeHTTP(w, r)

		end := time.Now()
		latency := end.Sub(start)
		res := w.(negroni.ResponseWriter)

		logger.Logger.Infof("requestid=%s status=%d latency=%v method=%s path=%s", requestId, res.Status(), latency, r.Method, path)
	}
}
