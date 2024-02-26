package main

import (
	"net/http"
	"os"

	"github.com/chazapp/prometheus-openweathermap-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "prometheus-openweathermap-exporter",
		Usage: "A Prometheus Exporter for OpenWeatherMap",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run the Prometheus Exporter",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "apiKey",
						EnvVars:  []string{"OPENWEATHERMAP_APIKEY"},
						Usage:    "OpenWeatherMap API key",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "cities",
						Usage:    "Cities to scrape values",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "logLevel",
						Usage: "Set log level",
						Value: "debug",
					},
				},
				Action: func(c *cli.Context) error {
					setLogLevel(c.String("logLevel"))
					apiKey := c.String("apiKey")
					cities := c.StringSlice("cities")

					collector := collector.NewOWMCollector(apiKey, cities)
					prometheus.MustRegister(&collector)

					http.Handle("/metrics", promhttp.Handler())
					return http.ListenAndServe(":9001", nil)
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Err(err)
	}
}

func setLogLevel(logLevel string) {
	switch logLevel {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
