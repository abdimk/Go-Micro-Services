package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8081")
)


type Server struct {
	addr       string	
	mux        *http.ServeMux
	httpServer *http.Server
}

func NewServer(addr string, mux *http.ServeMux) *Server{
	return &Server{
		addr:       addr,
		mux:        mux,
		httpServer: &http.Server{Addr: addr, Handler: mux},
	}
}

func (s *Server) Run() error {
	log.Printf("Server is running in port %s", s.addr)

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf(" %v", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}



func main() {
	log.Println("Starting API Gateway")
	server := NewServer(httpAddr, NewMux())
	serverErrors := make(chan error, 1)
	
	go func() {
		serverErrors <- server.Run()
	}()
	
	shutdown := make(chan os.Signal, 1)
	
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	
	
	select {
		case err := <- serverErrors:
			log.Printf("error starting the server: %v", err)
			
		case sig := <- shutdown:
			log.Printf("server is shutting down to %v signal",sig)
			
		
			ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
			
			defer cancel()
			
			if err := server.Shutdown(ctx); err != nil{
				log.Printf("could not stop the server gracefully:%v", err)
				server.Close()
			}
	}

}





