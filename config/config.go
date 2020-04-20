/*
Package config is used to provide configuration to the primary application. This
package works by parsing environment variables. Then provided access to those
values via a Config type.
*/
package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

// Config is the base configuration for MockItOut.
type Config struct {
	// Debug specifies if debug logging should be turned on or off.
	Debug bool `env:"DEBUG"`

	// DisableLogging will turn all logging off, use wisely.
	DisableLogging bool `env:"DISABLE_LOGGING" envDefault:"false"`

	// EnableTLS specifies if TLS is enabled for this service.
	EnableTLS bool `env:"ENABLE_TLS" envDefault:"true"`

	// ListenAddr specifies the HTTP Listener address used for this service.
	ListenAddr string `env:"LISTEN_ADDR" envDefault:"0.0.0.0:8443"`

	// CertFile specifies the location of the TLS certificate file. This is used only
	// if TLS is Enabled.
	CertFile string `env:"CERT_FILE"`

	// KeyFile specifies the location of the TLS key file. This is used only if TLS
	// is Enabled.
	KeyFile string `env:"KEY_FILE"`

	// GenCerts specifies if the service should generate a test Key and Cert to use
	// with TLS.
	GenCerts bool `env:"GEN_CERTS" envDefault:"false"`

	// MocksFile specifies the full path to the mocks configuration file. This value
	// must be set or the service will not start.
	MocksFile string `env:"MOCKS_FILE"`
}

// New will create a new Config instance with strong defaults.
func New() Config {
	c := Config{
		ListenAddr: "0.0.0.0:8443",
		EnableTLS:  true,
		Debug:      false,
		GenCerts:   true,
	}
	return c
}

// NewFromEnv will create a Config instance with data loaded from environment
// variables. When environment varaibles are not defined, defaults will be used.
func NewFromEnv() (Config, error) {
	c := New()
	err := env.Parse(&c)
	if err != nil {
		return New(), fmt.Errorf("could not load config from environment - %s", err)
	}
	return c, nil
}
