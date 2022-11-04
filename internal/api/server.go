package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/flash-cards-vocab/backend/pkg/application"
)

type Server struct {
	*http.Server
}

// NewServer creates a new Server instance
func NewServer() (*Server, error) {
	fmt.Println("Configuring server...")
	app, err := application.Get()
	if err != nil {
		return nil, err
	}

	router, err := NewRouter(app)
	if err != nil {
		log.Panicln("Failed to start new router:", err)
	}
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", app.Config.Host, app.Config.Port),
		Handler: router,
	}
	return &Server{&server}, nil
}

// Start - starts server
func (srv *Server) Start() {
	fmt.Println("Starting server..")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	fmt.Println("Listening on", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	fmt.Println("Shutting down server.. Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println("Server gracefully stopped")
}
