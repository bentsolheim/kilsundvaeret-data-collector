package service

import (
	"fmt"
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"time"
)

type LocationForecast20Response struct {
	Properties Properties
}

type Properties struct {
	Meta       Meta
	Timeseries []TimeSerie
}

type Meta struct {
	UpdatedAt string `json:"updated_at"`
}

func (m Meta) UpdatedAtParsed() *time.Time {
	return parseTime(m.UpdatedAt)
}

type TimeSerie struct {
	Time string
	Data Data
}

func (t TimeSerie) TimeParsed() *time.Time {
	return parseTime(t.Time)
}

type Data struct {
	Instant InstantForecast
}

type InstantForecast struct {
	Details InstantDetails
}

type InstantDetails struct {
	AirPressureAtSeaLevel float32 `json:"air_pressure_at_sea_level"`
	AirTemperature        float32 `json:"air_temperature"`
	/*
		CloudAreaFraction        float32 `json:"cloud_area_fraction"`
		CloudAreaFractionHigh    float32 `json:"cloud_area_fraction_high"`
		CloudAreaFractionLow     float32 `json:"cloud_area_fraction_low"`
		CloudAreaFractionMedium  float32 `json:"cloud_area_fraction_medium"`
		DewPointTemperature      float32 `json:"dew_point_temperature"`
		FogAreaFraction          float32 `json:"fog_area_fraction"`
	*/
	RelativeHumidity float32 `json:"relative_humidity"`
	//UltravioletIndexClearSky float32 `json:"ultraviolet_index_clear_sky"`
	WindFromDirection float32 `json:"wind_from_direction"`
	WindSpeed         float32 `json:"wind_speed"`
	//WindSpeedOfGust          float32 `json:"wind_speed_of_gust"`
}

func NewMetService(metProxyUrl string) MetService {
	return MetService{metProxyUrl: metProxyUrl}
}

type MetService struct {
	metProxyUrl string
}

type InstantForecastWrapper struct {
	TimeSerie
	UpdatedAt string
}

func (s MetService) GetMostRecentImmediateForecast(lat string, lon string, alt string) (*InstantForecastWrapper, error) {
	forecast, err := s.LoadForecast(lat, lon, alt)
	if err != nil {
		return nil, err
	}

	currentHourTimeSerie := forecast.Properties.Timeseries[0]
	for _, t := range forecast.Properties.Timeseries {
		if time.Now().After(*t.TimeParsed()) {
			currentHourTimeSerie = t
			continue
		}
		break
	}

	return &InstantForecastWrapper{
		TimeSerie: currentHourTimeSerie,
		UpdatedAt: forecast.Properties.Meta.UpdatedAt,
	}, nil
}

func (s MetService) LoadForecast(lat string, lon string, alt string) (*LocationForecast20Response, error) {
	url := fmt.Sprintf("%s/weatherapi/locationforecast/2.0/compact?lat=%s&lon=%s&altitude=%s", s.metProxyUrl, lat, lon, alt)
	log.Debugf("Getting forecast from %s", url)
	response := LocationForecast20Response{}
	if err := utils.HttpGetJson(url, &response); err != nil {
		return nil, stacktrace.Propagate(err, "error while reading forecast from met-proxy")
	}
	return &response, nil
}

func parseTime(timeString string) *time.Time {
	parsed, err := time.Parse("2006-01-02T15:04:05Z", timeString)
	if err != nil {
		return nil
	}
	return &parsed
}
