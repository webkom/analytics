package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

var addr = flag.String("listen-address", ":8000", "The address to listen on for HTTP requests.")

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {

	raven.CapturePanic(func() {

		log.Warn("starting_server")

		flag.Parse()
		http.Handle("/metrics", promhttp.Handler())

		err := http.ListenAndServe(*addr, nil)

		if err != nil {
			log.Panic(err)
		}

	}, nil)
}
