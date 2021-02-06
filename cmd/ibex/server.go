package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type Server struct {
	server *http.Server
	logger *log.Logger
}

type key int

const (
	requestIDKey key = 0
)

var (
	healthy int32
	spath   string
	logger  *log.Logger
)

func NewServer(address, path string) Server {

	logger = log.New(os.Stdout, "[Ibex]: ", log.LstdFlags)
	logger.Println("Server is starting...")

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	router := http.NewServeMux()
	router.Handle("/depot", depot())

	server := &http.Server{
		Addr:         address,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	spath = path
	return Server{
		server: server,
		logger: logger,
	}

}

func (s *Server) Start() {

	s.logger.Println("Server is ready to handle requests at", s.server.Addr)
	atomic.StoreInt32(&healthy, 1)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatalf("Could not listen on %s: %v\n", s.server.Addr, err)
	}
	defer s.Stop()

}

func (s *Server) Stop() {

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.logger.Println("Server is shutting down...")
	atomic.StoreInt32(&healthy, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)

	<-done
	s.logger.Println("Server stopped")

}

func health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})

}
func depot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/depot" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=filename")
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		content, err := ioutil.ReadFile(path)
		if err != nil {
			logger.Printf("Error: %v\n", err)
		}
		io.Copy(w, bytes.NewReader(content))
	})
}
