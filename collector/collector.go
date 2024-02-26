package collector

import (
	"sync"

	"github.com/chazapp/prometheus-openweathermap-exporter/openweathermap"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type OWMCollector struct {
	owmApiKey string
	wg        sync.WaitGroup

	WeatherDataList []*openweathermap.WeatherMetrics

	tempCurrentCelsius  *prometheus.Desc
	tempForecastCelsius *prometheus.Desc
	rainCurrentMMH      *prometheus.Desc
	rainForecastMMH     *prometheus.Desc
}

func NewOWMCollector(apiKey string, cities []string) OWMCollector {

	weatherDataList := make([]*openweathermap.WeatherMetrics, len(cities))
	log.Debug().Msg("Building metrics for following cities")
	for idx, city := range cities {
		weatherDataList[idx] = openweathermap.NewWeatherMetrics(city, apiKey)
	}

	return OWMCollector{
		owmApiKey:           apiKey,
		WeatherDataList:     weatherDataList,
		tempCurrentCelsius:  prometheus.NewDesc("temperature_current_celsius", "Current recorded temperature in Celsius", []string{"city"}, nil),
		tempForecastCelsius: prometheus.NewDesc("temperature_forecast_celsius", "4-hour forecasted temperature in Celsius", []string{"city"}, nil),

		rainCurrentMMH:  prometheus.NewDesc("rain_current_mmh", "Current recorded rain volume in mm/h", []string{"city"}, nil),
		rainForecastMMH: prometheus.NewDesc("rain_forecast_mmh", "4-hour forecasted rain volume in mm/h", []string{"city"}, nil),
	}
}

func (c *OWMCollector) Collect(ch chan<- prometheus.Metric) {
	log.Debug().Msg("Collect called")
	for idx := range c.WeatherDataList {
		c.wg.Add(1)
		go func(idx int) {
			defer c.wg.Done()
			c.WeatherDataList[idx].GetWeatherMetrics(c.owmApiKey)
		}(idx)
	}

	c.wg.Wait()

	for _, item := range c.WeatherDataList {
		ch <- prometheus.MustNewConstMetric(c.tempCurrentCelsius, prometheus.GaugeValue, item.CurrentTemperature, item.City)
		ch <- prometheus.MustNewConstMetric(c.rainCurrentMMH, prometheus.GaugeValue, item.CurrentRain, item.City)
		ch <- prometheus.MustNewConstMetric(c.tempForecastCelsius, prometheus.GaugeValue, item.ForecastTemperature, item.City)
		ch <- prometheus.MustNewConstMetric(c.rainForecastMMH, prometheus.GaugeValue, item.ForecastRain, item.City)
	}

	return
}

func (c *OWMCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.tempCurrentCelsius
	ch <- c.tempForecastCelsius
	ch <- c.rainCurrentMMH
	ch <- c.rainForecastMMH
	return
}
