package http

import (
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"
	"time"
)


type APIError struct {
	Error string
}

type APIResponse struct {
	Message string `json:"message"`
}

type apiFunc func (w http.ResponseWriter, r *http.Request) error


type HttpHandler struct {
	Service domain.TripService
}


type previewTripRequest struct {
	UserID       string          `json:"userID"`  
	Pickup       types.Coordinate `json:"pickup"`
	Destination  types.Coordinate  `json:"destination"`
}




func (s *HttpHandler) HandleTripPreview(w http.ResponseWriter, r *http.Request) error{
	var reqBody previewTripRequest
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil{
		return WriteJSON(w, http.StatusBadRequest,"failed to parse JSON data")
	}
	
	ctx := r.Context()
	fare := &domain.RideFareModel{
		UserId: "43",
	}
	
	t, err := s.Service.CreateTrip(ctx, fare)
	if err != nil{
		log.Println(err)
	}
	return WriteJSON(w, http.StatusOK,t)
}


func NewMux(handler *HttpHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /preview", makeHTTPHandlerFunc(LogMiddleware(handler.HandleTripPreview)))
	return mux
}



func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			log.Printf("handler error on %s %s: %v", r.Method, r.URL.Path, err)
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}





// Logger 
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

func LogMiddleware(next apiFunc) apiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		err := next(rec, r)

		log.Printf("%s %s -> %d (%v)", r.Method, r.URL.Path, rec.status, time.Since(start))

		return err
	}
}
