package main

import (
	"context"
	// "errors"
	// "fmt"
	"log"
	"net"
	// "net/http"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"syscall"

	// h "ride-sharing/services/trip-service/internal/infrastructure/http"

	grpcserver "google.golang.org/grpc"
	grpc "ride-sharing/services/trip-service/internal/infrastructure/grpc"
)


var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
	GrpcAddr = ":9093"
)



// type Server struct {
// 	addr string	
// 	mux *http.ServeMux
// }

// func NewServer(addr string, mux *http.ServeMux) *Server{
// 	return &Server{addr: addr, mux: mux}
// }

// func (s *Server) Run() error {
// 	log.Printf("Server is running in port %s", s.addr)

// 	if err := http.ListenAndServe(s.addr, s.mux); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 		return fmt.Errorf(" %v", err)
// 	}
// 	return nil
// }


func main(){
	inmemRepo := repository.NewInMemoryRepository()
	svc := service.NewTripUserService(inmemRepo)
	// httphandler := h.HttpHandler{Service: svc}
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	go func(){
		sigCh := make(chan os.Signal,1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
	}()

	log.Println("Trip Service is running")
	lis, err := net.Listen("tcp", GrpcAddr)
	
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}
	
	// Starting the grpc Server
	grpcServer := grpcserver.NewServer()
	
	
	grpc.NewGRPCHandler(grpcServer, svc)
	// TODO initilize out grpc hadler implementaion
	
	
	log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())
	
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()
	
	
	// wait for the shutdown signal
	<- ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
	
	
	// server := NewServer(httpAddr, h.NewMux(&httphandler))
	// if err := server.Run(); err != nil{
	// 	log.Fatalf("Server stopped: %v", err)
	// }
	
}



