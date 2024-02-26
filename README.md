# Prometheus-OpenWeatherMap-Exporter

This project allows for monitoring of OpenWeatherMap data via a Prometheus-compatible endpoint

## Usage

Clone the repository, build the application, run the executable

```bash
$ git clone git@github.com:/chazapp/prometheus-openweathermap-exporter && cd prometheus-openweathermap-exporter
...
$ go build main.go -o owm
$ ./owm
NAME:
   prometheus-openweathermap-exporter - A Prometheus Exporter for OpenWeatherMap

USAGE:
   prometheus-openweathermap-exporter [global options] command [command options]

COMMANDS:
   run      Run the Prometheus Exporter
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
$ ./owm-exporter --api-key 0xdeadbeef --cities paris,london,new-york
Exposing metrics on ::9001/metrics
```

## Production usage

The `deployment/` directory contains the Terraform code to apply to a Minikube cluster. The cluster will get configured with:

- Kube-Prometheus-Stack chart
- Grafana
- Prometheus-OpenWeatherMap-Exporter  
