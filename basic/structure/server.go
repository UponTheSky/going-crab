package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// define our own Server struct
type Server struct {
	protocol string
	port     int
	mux      *http.ServeMux
	logger   *log.Logger
}

func NewServer(protocol string, port int) *Server {
	// mux here
	mux := http.NewServeMux()

	// logger is for the server logging
	logger := log.New(
		os.Stderr,
		"[Going Crab] ",
		log.LstdFlags|log.Lshortfile|log.LUTC,
	)

	return &Server{
		protocol: protocol,
		port:     port,
		mux:      mux,
		logger:   logger,
	}
}

func (s *Server) Run(app *App) error {
	// register the application controllers
	app.RegisterController(s.mux)

	// create http.Server instance
	srv := &http.Server{}
	srv.Addr = fmt.Sprintf(":%v", s.port)
	srv.Handler = s.mux

	// listen to the specified IP address
	listener, err := net.Listen(s.protocol, srv.Addr)

	if err != nil {
		return err // unexpected error
	}

	defer listener.Close()

	// run the http.Server
	return srv.Serve(listener) // server.Serve always returns an error
}
