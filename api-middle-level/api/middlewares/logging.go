package middlewares

import (
	"log"
	"net/http"
)

func LoggingMdw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := newTraceID()
		log.Printf("[%d]%s %s\n", traceID, r.RequestURI, r.Method)

		ctx := SetTraceID(r.Context(), traceID)
		r = r.WithContext(ctx)

		rlw := NewResLoggingWriter(w)
		next.ServeHTTP(rlw, r)
		log.Printf("[%d]res: %d\n", traceID, rlw.code)
	})
}

type resLoggingWriter struct {
	http.ResponseWriter
	code int
}

func NewResLoggingWriter(w http.ResponseWriter) *resLoggingWriter {
	return &resLoggingWriter{ResponseWriter: w, code: http.StatusOK}
}

func (rsw *resLoggingWriter) WriteHeader(code int) {
	rsw.code = code
	rsw.ResponseWriter.WriteHeader(code)
}
