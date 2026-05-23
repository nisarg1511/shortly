package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code // Capture the status code
	w.ResponseWriter.WriteHeader(code)
}

func RequestLogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			method = r.Method
			path   = r.URL.Path
			status = http.StatusOK
			start  = time.Now()
			// user     = r.UserAgent()
			// ip       = r.RemoteAddr
		)
		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     status,
		}

		next.ServeHTTP(wrapper, r)
		// userAttr := slog.Group("", "user-agent", user, "ip", ip)
		log.Printf("[INFO] %s %s %d %vms\n", method, path, wrapper.statusCode, time.Since(start).Milliseconds())
	})
}
