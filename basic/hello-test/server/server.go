package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

type Server struct {
	protocol string
	port     int
	mux      *http.ServeMux
	logger   *log.Logger
}

func New(protocol string, port int) *Server {
	mux := http.NewServeMux()
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

func (ws *Server) RegisterEndpoint(path string, handler http.Handler) {
	ws.mux.Handle(path, handler)
}

func (ws *Server) Run() error {
	srv := &http.Server{}
	srv.Addr = fmt.Sprintf(":%v", ws.port)
	srv.Handler = ws.mux

	listener, err := net.Listen(ws.protocol, srv.Addr)

	if err != nil {
		return err // unexpected error
	}

	defer listener.Close()

	return srv.Serve(listener) // server.Serve always returns an error
}
