package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/madflojo/mockitout/app"
	"github.com/madflojo/mockitout/config"
	"github.com/sirupsen/logrus"
	"os"
)

// options is used to process command line arguments. These are minimal, in lieu of using
// environment variables to control the server.
type options struct {
	Debug bool `long:"debug" description:"Enable debug logging"`
}

func main() {
	// Initiate a simple logger
	log := logrus.New()

	// Parse command line arguments
	var opts options
	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		log.Fatalf("Unable to parse command line options, shutting down - %s", err)
	}

	// Load Config from environment
	env, err := config.NewFromEnv()
	if err != nil {
		log.Fatalf("Unable to load config shutting down - %s", err)
	}

	// Override Debug with command-line
	if opts.Debug {
		env.Debug = opts.Debug
	}

	// Run application
	err = app.Run(env)
	if err != nil || err != app.ErrShutdown {
		log.Fatalf("MockItOut stopped abruptly - %s", err)
	}
	log.Infof("MockItOut shutdown gracefully - %s", err)
}
