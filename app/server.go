package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

// server is used as an interface for for managing the HTTP server
type server struct {
	// httpServer is the primary HTTP server
	httpServer *http.Server

	// httpRouter is used to store and access the HTTP Request Router
	httpRouter *httprouter.Router
}

// Health is used to handle HTTP Health requests to this service
func (s *server) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

// MockHandler is used to handle HTTP requests to the Mock Server.
func (s *server) MockHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

// middleware is used to intercept incoming HTTP calls and apply general functions upon them. e.g. Metrics...
func (s *server) middleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.WithFields(logrus.Fields{
			"method":        r.Method,
			"remote-addr":   r.RemoteAddr,
			"http-protocol": r.Proto,
		}).Debugf("HTTP Request to %s", r.URL)
		n(w, r, ps)
	}
}
