package openweathermap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type WeatherMetrics struct {
	City string

	longitude float64
	latitude  float64

	CurrentTemperature float64
	CurrentRain        float64

	ForecastTemperature float64
	ForecastRain        float64
}

type OWMGeoResp struct {
	Name       string      `json:"name"`
	Latitude   float64     `json:"lat"`
	Longitude  float64     `json:"lon"`
	Country    string      `json:"country"`
	State      string      `json:"state"`
	LocalNames interface{} `json:"local_names"`
}

type OWMOnecallResp struct {
	Current OWMDataResp   `json:"current"`
	Hourly  []OWMDataResp `json:"hourly"`
}

type OWMDataResp struct {
	Timestamp   float64     `json:"dt"`
	Temperature float64     `json:"temp"`
	Rain        OWMRainResp `json:"rain"`
}

type OWMRainResp struct {
	OneHour float64 `json:"1h"`
}

func NewWeatherMetrics(city string, apiKey string) *WeatherMetrics {
	weatherMetrics := &WeatherMetrics{
		City: city,
	}
	if err := weatherMetrics.setLatitudeLongitudeFromCity(apiKey); err != nil {
		log.Panic().Err(err)
	}
	log.Debug().Msgf("City: %s, Latitude: %f, Longitude: %f", weatherMetrics.City, weatherMetrics.latitude, weatherMetrics.longitude)
	return weatherMetrics
}

func (wm *WeatherMetrics) GetWeatherMetrics(apiKey string) {

	log.Debug().Msgf("Collecting weather metrics for city: %s", wm.City)
	baseURL := "https://api.openweathermap.org/data/3.0/onecall"
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", wm.latitude))
	params.Add("lon", fmt.Sprintf("%f", wm.longitude))
	params.Add("exclude", "minutely,daily,alerts")
	params.Add("appid", apiKey)
	params.Add("units", "metric")

	u, _ := url.ParseRequestURI(baseURL)
	u.RawQuery = params.Encode()
	body, err := httpGetRequestToBody(fmt.Sprintf("%v", u))

	if err != nil {
		log.Error().Err(err).Msgf("OWM API Error with req: %s", u)
		return
	}

	var data OWMOnecallResp

	if err := json.Unmarshal(body, &data); err != nil {
		log.Error().Err(err)
		return
	}

	wm.CurrentTemperature = data.Current.Temperature
	wm.CurrentRain = data.Current.Rain.OneHour

	if len(data.Hourly) < 4 {
		log.Error().Msgf("OWM API Malformed forecast data. Count of hourly data %d < 4", len(data.Hourly))
		return
	}
	wm.ForecastTemperature = data.Hourly[4].Temperature
	wm.ForecastRain = data.Hourly[4].Rain.OneHour
}

func (wm *WeatherMetrics) setLatitudeLongitudeFromCity(apiKey string) error {
	baseURL := "https://api.openweathermap.org/geo/1.0/direct"
	params := url.Values{}
	params.Add("q", wm.City)
	params.Add("limit", "1")
	params.Add("appid", apiKey)

	u, _ := url.ParseRequestURI(baseURL)
	u.RawQuery = params.Encode()
	body, err := httpGetRequestToBody(fmt.Sprintf("%v", u))

	var data []OWMGeoResp

	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if len(data) != 1 {
		return fmt.Errorf("OpenWeatherMap API returned malformed data len %d =! 1", len(data))
	}
	wm.latitude = data[0].Latitude
	wm.longitude = data[0].Longitude
	return err
}

func httpGetRequestToBody(encodedURL string) ([]byte, error) {
	resp, err := http.Get(encodedURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenWeatherMap API returned status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}
