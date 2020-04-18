package app

import (
	"bytes"
	"crypto/tls"
	"github.com/madflojo/mockitout/config"
	"github.com/madflojo/testcerts"
	"net/http"
	"testing"
	"time"
)

func TestRunningServer(t *testing.T) {
	go func() {
		err := Run(config.Config{
			EnableTLS:      false,
			ListenAddr:     "localhost:9000",
			DisableLogging: true,
		})
		if err != nil && err != ErrShutdown {
			t.Errorf("Run unexpectedly stopped - %s", err)
		}
	}()

	// Wait for app to start
	time.Sleep(1 * time.Second)

	t.Run("Check Health HTTP Handler", func(t *testing.T) {
		r, err := http.Get("http://localhost:9000/health")
		if err != nil {
			t.Errorf("Unexpected error when requesting health status - %s", err)
		}
		if r.StatusCode != 200 {
			t.Errorf("Unexpected http status code when checking health - %d", r.StatusCode)
		}
	})

	// Clean up
	Stop()
}

func TestRunningTLSServer(t *testing.T) {
	// Create Test Certs
	err := testcerts.GenerateCertsToFile("/tmp/cert", "/tmp/key")
	if err != nil {
		t.Errorf("Failed to create certs - %s", err)
		t.FailNow()
	}

	// Disable Host Checking globally
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Start Server in goroutine
	go func() {
		err := Run(config.Config{
			EnableTLS:      true,
			ListenAddr:     "localhost:9000",
			CertFile:       "/tmp/cert",
			KeyFile:        "/tmp/key",
			DisableLogging: true,
		})
		if err != nil && err != ErrShutdown {
			t.Errorf("Run unexpectedly stopped - %s", err)
		}
	}()

	// Wait for app to start
	time.Sleep(1 * time.Second)

	t.Run("Check Health HTTP Handler", func(t *testing.T) {
		r, err := http.Get("https://localhost:9000/health")
		if err != nil {
			t.Errorf("Unexpected error when requesting health status - %s", err)
			t.FailNow()
		}
		if r.StatusCode != 200 {
			t.Errorf("Unexpected http status code when checking health - %d", r.StatusCode)
		}
	})

	// Clean up
	Stop()
}
