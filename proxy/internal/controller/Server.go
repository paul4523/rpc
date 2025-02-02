package controller

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	httpserver *http.Server
	port       string
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpserver: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		port: addr,
	}

}
func (s *Server) Serve() error {
	log.Println("Starting server...")
	if err := s.httpserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
		return err
	}
	return nil
}
func (s *Server) Shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Server stopped gracefully")
	s.httpserver.Shutdown(context.Background())
}
