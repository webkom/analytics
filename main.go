package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	_ "github.com/lib/pq"
	"os"
)

func init() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	raven.CapturePanic(func() {
		migrate := flag.Bool("migrate", false, "migrate database")
		flag.Parse()

		listenAddress := os.Getenv("LISTEN_ADDRESS")

		app := App{}
		app.Initialize(
			os.Getenv("POSTGRES_URL"),
		)

		if *migrate {
			log.Info("migrating_database")
			app.Migrate()
		} else {
			log.Info("starting_server")
			app.Run(listenAddress)
		}
	}, nil)
}
