package main


import (
	"log"
	"net/http"
	"time"
)


type statusRecorder struct {
	http.ResponseWriter
	status  int
	written bool
}

func (s *statusRecorder) WriteHeader(code int) {
	if !s.written {
		s.status = code
		s.written = true
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	s.written = true
	return s.ResponseWriter.Write(b)
}

func logMiddleware(next apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		err := next(rec, r)

		log.Printf("%s %s -> %d (%v)", r.Method, r.URL.Path, rec.status, time.Since(start))

		return err
	}
}
