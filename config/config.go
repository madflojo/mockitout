/*
Package config is used to provide configuration for the rest of MockItOut.
*/
package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

// Config is the base configuration for MockItOut.
type Config struct {
	// Debug determines if debug logging should be turned on or off
	Debug bool `env:"DEBUG"`

	// DisableLogging will turn all logging off, use wisely
	DisableLogging bool `env:"DISABLE_LOGGING" envDefault:"false"`

	// EnableTLS determines if TLS is enabled for this service
	EnableTLS bool `env:"ENABLE_TLS" envDefault:"true"`

	// ListenAddr is the HTTP Listener address used for this service
	ListenAddr string `env:"LISTEN_ADDR" envDefault:"0.0.0.0:443"`

	// CertFile is used to specify the location of the TLS certificate file
	CertFile string `env:"CERT_FILE"`

	// KeyFile is used to specify the location of the TLS key file
	KeyFile string `env:"KEY_FILE"`

	// MocksFile is used to specify the full path to the mocks configuration file
	MocksFile string `env:"ROUTES_FILE"`
}

// New will create a new Config instance with strong defaults.
func New() Config {
	c := Config{
		ListenAddr: "0.0.0.0:443",
		EnableTLS:  true,
		Debug:      false,
	}
	return c
}

// New will create a Config instance with data loaded from environment variables. When environment varaibles are not
// defined, defaults will be used.
func NewFromEnv() (Config, error) {
	c := New()
	err := env.Parse(&c)
	if err != nil {
		return New(), fmt.Errorf("could not load config from environment - %s", err)
	}
	return c, nil
}
