package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ConfigureAndRun(ip, port string, debug bool, router *httprouter.Router) *Server {
	log.Printf("[Server] Bind application to host: %s and port: %s", ip, port)

	var listener net.Listener
	var err error
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Fatal(err)
	}

	c := s.cors(debug)

	s.httpServer = &http.Server{
		Handler:      c.Handler(router),
		WriteTimeout: 4 * time.Second,
		ReadTimeout:  4 * time.Second,
	}

	if err := s.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Fatal("[Server] Server shutdown")
		default:
			log.Fatal("[Server] " + err.Error())
		}
	}

	err = s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return &Server{}
}

func (s *Server) cors(debug bool) *cors.Cors {
	return cors.New(
		cors.Options{
			AllowedMethods:     []string{http.MethodGet, http.MethodPost},
			AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
			AllowCredentials:   true,
			AllowedHeaders:     []string{"Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Access"},
			OptionsPassthrough: true,
			ExposedHeaders:     []string{"Location", "Authorization", "Content-Disposition"},
			Debug:              false,
		},
	)
}
