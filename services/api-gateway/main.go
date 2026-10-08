package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
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



func main() {
	log.Println("Starting API Gateway")
	server := NewServer(httpAddr, NewMux())
	
	if err := server.Run(); err != nil{
		log.Fatalf("Server stopped: %v", err)
	}

}





