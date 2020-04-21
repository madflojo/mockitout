/*
Package app is the primary runtime for MockItOut.
*/
package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/madflojo/mockitout/config"
	"github.com/madflojo/mockitout/mocks"
	"github.com/madflojo/testcerts"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// Common errors returned by this app.
var (
	ErrShutdown = fmt.Errorf("Application shutdown gracefully")
)

// srv is the global reference for the HTTP Server.
var srv *server

// cfg is used across the app package to contain configuration.
var cfg config.Config

// log is used across the app package for logging.
var log *logrus.Logger

// mocked is the defined server mocks loaded from config.
var mocked mocks.Mocks

// Run starts the primary application. It handles starting background services,
// populating package globals & structures, and clean up tasks.
func Run(c config.Config) error {
	var err error
	// Apply config provided by main
	cfg = c

	// Initiate the logger
	log = logrus.New()
	if cfg.Debug {
		log.Level = logrus.DebugLevel
		log.Debug("Enabling Debug Logging")
	}
	if cfg.DisableLogging {
		log.Level = logrus.FatalLevel
	}

	// Setup the HTTP Server
	srv = &server{
		httpRouter: httprouter.New(),
	}
	srv.httpServer = &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: srv.httpRouter,
	}
	srv.httpRouter.SaveMatchedRoutePath = true

	// Setup TLS Configuration
	if cfg.EnableTLS {
		if cfg.GenCerts {
			// Create Test Certs
			cfg.CertFile = "/tmp/cert"
			cfg.KeyFile = "/tmp/key"
			log.Infof("Certificate Generation was enabled, creating new test certs at %s and %s", cfg.CertFile, cfg.KeyFile)
			err := testcerts.GenerateCertsToFile(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				return fmt.Errorf("Could not generate test certificates - %s", err)
			}
			defer os.Remove(cfg.CertFile)
			defer os.Remove(cfg.KeyFile)
		}

		srv.httpServer.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
	}

	// Register Health Check Handler
	srv.httpRouter.GET("/health", srv.middleware(srv.Health))

	// Start Registering Custom Mock HTTP Routes
	mocked, err = mocks.FromFile(cfg.MocksFile)
	if err != nil {
		return err
	}
	for m, r := range mocked.Routes {
		log.Infof("Registering mock %s with path %s", m, r.Path)
		srv.httpRouter.GET(r.Path, srv.middleware(srv.MockHandler))
		srv.httpRouter.POST(r.Path, srv.middleware(srv.MockHandler))
		srv.httpRouter.PUT(r.Path, srv.middleware(srv.MockHandler))
		srv.httpRouter.DELETE(r.Path, srv.middleware(srv.MockHandler))
	}

	// Start HTTP Listener
	log.Infof("Starting Listener on %s", cfg.ListenAddr)
	if cfg.EnableTLS {
		err := srv.httpServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			if err == http.ErrServerClosed {
				return ErrShutdown
			}
			return err
		}
	}
	err = srv.httpServer.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			return ErrShutdown
		}
		return err
	}

	return nil
}

// Stop is used to gracefully shutdown the server.
func Stop() {
	defer srv.httpServer.Shutdown(context.Background())
}
