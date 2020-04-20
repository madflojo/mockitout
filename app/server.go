package app

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/madflojo/mockitout/mocks"
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
	var route mocks.Route
	var path string
	var ok bool
	if path, ok = mocked.Paths[r.RequestURI]; !ok {
		log.Errorf("Request URI %s not found within Mocks file - available paths %+v", r.RequestURI, mocked)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if route, ok = mocked.Routes[path]; !ok {
		log.Errorf("Request URI %s not found within Mocks file when looking for route named %s", r.RequestURI, path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.WithFields(logrus.Fields{
		"return-code": route.ReturnCode,
		"path":        route.Path,
	}).Infof("Mocked end-point found for %s", r.RequestURI)

	// Write out user defined response code or a 404
	w.WriteHeader(route.ReturnCode)
	// Write Body to caller
	fmt.Fprintf(w, "%s", route.Body)
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