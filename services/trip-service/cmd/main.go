package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	
	h "ride-sharing/services/trip-service/internal/infrastructure/http"
)


var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)



type Server struct {
	addr string	
	mux *http.ServeMux
}

func NewServer(addr string, mux *http.ServeMux) *Server{
	return &Server{addr: addr, mux: mux}
}

func (s *Server) Run() error {
	log.Printf("Server is running in port %s", s.addr)

	if err := http.ListenAndServe(s.addr, s.mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf(" %v", err)
	}
	return nil
}


func main(){
	inmemRepo := repository.NewInMemoryRepository()
	svc := service.NewTripUserService(inmemRepo)
	httphandler := h.HttpHandler{Service: svc}
	

	log.Println("Trip Service is running")
	server := NewServer(httpAddr, h.NewMux(&httphandler))
	if err := server.Run(); err != nil{
		log.Fatalf("Server stopped: %v", err)
	}
	
}



