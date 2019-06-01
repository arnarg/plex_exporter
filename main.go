package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/arnarg/plex_exporter/collector"
	"github.com/arnarg/plex_exporter/config"
	"github.com/arnarg/plex_exporter/plex"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var Version = "0.1.0"

func Run(c *cli.Context) error {
	addr := c.String("listen-address")
	path := c.String("config-path")

	// Loading persisted config
	conf, isNew, err := config.Load(path)
	if err != nil {
		return fmt.Errorf("Could not load config file: %s", err)
	}

	// Save config if it was just generated
	// I mainly do this to run into config saving problems
	// before requesting an authentication PIN
	if isNew {
		log.Debugf("Saving new generated config to %s", path)
		err = config.Save(conf, path)
		if err != nil {
			return fmt.Errorf("Could not save config file: %s", err)
		}
	}

	// Create a Plex client
	clientLogger := log.WithFields(log.Fields{"context": "client"})
	client := plex.NewPlexClient(conf, Version, clientLogger)

	// Create a pin request for Plex authentication
	if client.Token == "" {
		pinRequest, err := client.GetPinRequest()
		if err != nil {
			return fmt.Errorf("Could not make a pin request: %s", err)
		}
		log.Infof("Got PIN code: %s", pinRequest.Code)
		log.Info("Go to https://plex.tv/pin and enter pin to authenticate.")

		// Repeatedly check pin request
		ticker := time.NewTicker(time.Second * 5)
		for t := range ticker.C {
			if pinRequest.Expiry.Before(t) {
				ticker.Stop()
				return fmt.Errorf("PIN expired, exiting.")
			}

			log.Debug("Checking PIN request")
			token, err := client.GetTokenFromPinRequest(pinRequest)
			if err != nil {
				ticker.Stop()
				return fmt.Errorf("Could not check PIN request: %s", err)
			}

			if token != "" {
				log.Info("Authenticated successfully")
				conf.Token = token
				ticker.Stop()
				break
			}
		}

		// Persist config
		err = config.Save(conf, path)
		if err != nil {
			return fmt.Errorf("Could not save config file: %s", err)
		}
	}

	// Initialize client
	err = client.Init()
	if err != nil {
		return fmt.Errorf("Could not initialize Plex client: %s", err)
	}

	// Create the Prometheus collector
	collectorLogger := log.WithFields(log.Fields{"context": "colletor"})
	pc := collector.NewPlexCollector(client, collectorLogger)
	prometheus.MustRegister(pc)

	// Start HTTP server
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Beginning to serve on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
	return nil
}

func Init(c *cli.Context) error {
	verbose := c.String("log-level")
	format := c.String("format")

	// Set verbosity level
	switch verbose {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "err":
		log.SetLevel(log.ErrorLevel)
	default:
		return fmt.Errorf("Available log levels are trace, debug, info, warn, err")
	}

	// Set log format
	switch format {
	case "text":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		return fmt.Errorf("Available log formats are text, json")
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "plex_exporter"
	app.Usage = "A Prometheus exporter that exports metrics on Plex Media Server."
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen-address, l",
			Value: ":9594",
			Usage: "Port for server",
		},
		cli.StringFlag{
			Name:  "config-path, c",
			Value: "/var/lib/plex_exporter/config.json",
			Usage: "Path to persistent authentication token",
		},
		cli.StringFlag{
			Name:  "log-level",
			Value: "info",
			Usage: "Verbosity level of logs",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "text",
			Usage: "Output format of logs",
		},
	}

	app.Action = Run
	app.Before = Init

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
