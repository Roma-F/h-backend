package rest

import (
	"net/http"
	"time"
)

func NewServer(handler http.Handler, port string) *http.Server {
	s := http.Server{
		Addr:         ":" + port,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}
	return &s
}
