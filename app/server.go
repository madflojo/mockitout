package app

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/madflojo/mockitout/mocks"
	"github.com/madflojo/mockitout/variable"
	"github.com/sirupsen/logrus"
)

// server is used as an interface for managing the HTTP server.
type server struct {
	// httpServer is the primary HTTP server.
	httpServer *http.Server

	// httpRouter is used to store and access the HTTP Request Router.
	httpRouter *httprouter.Router
}

// Health is used to handle HTTP Health requests to this service.
func (s *server) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
}

// MockHandler is used to handle HTTP requests to the Mock Server.
func (s *server) MockHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var route mocks.Route
	var path string
	var ok bool
	if path, ok = mocked.Paths[ps.MatchedRoutePath()]; !ok {
		log.Errorf("Request URI %s not found within Mocks file - available paths %+v", r.RequestURI, mocked)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if route, ok = mocked.Routes[path]; !ok {
		log.Errorf("Request URI %s not found within Mocks file when looking for route named %s", r.RequestURI, path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := variable.NewVariableInstance(r, w, ps)

	// Verify Return Code is set if not default to 200
	if route.ReturnCode == 0 {
		route.ReturnCode = 200
	}

	log.WithFields(logrus.Fields{
		"return-code": route.ReturnCode,
		"path":        route.Path,
	}).Infof("Mocked end-point found for %s", r.RequestURI)

	// Add any user defined headers
	for k, v := range route.ResponseHeaders {
		v, err := ctx.ReplaceVariables(v)
		if err != nil {
			log.WithFields(logrus.Fields{
				"path": route.Path,
			}).Errorf("Error parsing header variable %s - %s", v, err)
			continue
		}
		w.Header().Set(k, v)
	}

	// Write out user defined response code
	w.WriteHeader(route.ReturnCode)

	// Write Body to caller
	varBody, err := ctx.ReplaceVariables(route.Body)
	if err != nil {
		log.WithFields(logrus.Fields{
			"path": route.Path,
		}).Errorf("Error parsing body variable %s - %s", route.Body, err)
	}
	fmt.Fprintf(w, "%s", varBody)
}

// middleware is used to intercept incoming HTTP calls and apply general functions upon
// them. e.g. Metrics, Logging...
func (s *server) middleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Log the basics
		log.WithFields(logrus.Fields{
			"method":         r.Method,
			"remote-addr":    r.RemoteAddr,
			"http-protocol":  r.Proto,
			"headers":        r.Header,
			"content-length": r.ContentLength,
		}).Debugf("HTTP Request to %s", r.URL)

		// Call registered handler
		n(w, r, ps)
	}
}
