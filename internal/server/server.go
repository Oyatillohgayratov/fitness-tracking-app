package server

import (
	"log/slog"
	"net/http"
)

type HTTPServer struct {
	srv *http.Server
	log slog.Logger
}

func New(addr string, hand http.Handler, log slog.Logger) HTTPServer {
	srv := &http.Server{
		Addr:    addr,
		Handler: hand,
	}

	return HTTPServer{
		srv: srv,
		log: log,
	}
}

func (s *HTTPServer) Run() error {
	s.log.Info("Server is running on " + s.srv.Addr)
	return s.srv.ListenAndServe()
}
