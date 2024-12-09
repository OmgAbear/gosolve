package http_interface

import (
	"fmt"
	"github.com/OmgAbear/gosolve/internal/config"
	"net/http"
	"time"
)

// LoggingMiddleware logs the details of each request and response
// There are probably some packages available already, to do this, I am just not familiar with them
// Generally, each project I worked on had some specific requirements for this and every time it was a custom imp
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := config.GetLogger()
		logger.Info(fmt.Sprintf("Request: %s %s %s", r.Method, r.RequestURI, r.RemoteAddr))

		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)
		logger.Info(fmt.Sprintf("Response: %d %s [%s]", recorder.statusCode, http.StatusText(recorder.statusCode), duration))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
