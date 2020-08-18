package service

import (
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
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

type TimeSerie struct {
	Time string
	Data Data
}

type Data struct {
	Instant InstantForecast
}

type InstantForecast struct {
	Details InstantDetails
}

type InstantDetails struct {
	AirPressureAtSeaLevel    float32 `json:"air_pressure_at_sea_level"`
	AirTemperature           float32 `json:"air_temperature"`
	CloudAreaFraction        float32 `json:"cloud_area_fraction"`
	CloudAreaFractionHigh    float32 `json:"cloud_area_fraction_high"`
	CloudAreaFractionLow     float32 `json:"cloud_area_fraction_low"`
	CloudAreaFractionMedium  float32 `json:"cloud_area_fraction_medium"`
	DewPointTemperature      float32 `json:"dew_point_temperature"`
	FogAreaFraction          float32 `json:"fog_area_fraction"`
	RelativeHumidity         float32 `json:"relative_humidity"`
	UltravioletIndexClearSky float32 `json:"ultraviolet_index_clear_sky"`
	WindFromDirection        float32 `json:"wind_from_direction"`
	WindSpeed                float32 `json:"wind_speed"`
	WindSpeedOfGust          float32 `json:"wind_speed_of_gust"`
}

func NewMetService(metProxyUrl string) MetService {
	return MetService{metProxyUrl: metProxyUrl}
}

type MetService struct {
	metProxyUrl string
}

type InstantForecastWrapper struct {
	InstantForecast
	UpdatedAt string
}

func (s MetService) GetMostRecentImmediateForecast(lat string, lon string, alt string) (*InstantForecastWrapper, error) {
	forecast, err := s.LoadForecast(lat, lon, alt)
	if err != nil {
		return nil, err
	}

	instantForecast := forecast.Properties.Timeseries[0].Data.Instant
	return &InstantForecastWrapper{
		InstantForecast: instantForecast,
		UpdatedAt:       forecast.Properties.Meta.UpdatedAt,
	}, nil
}

func (s MetService) LoadForecast(lat string, lon string, alt string) (*LocationForecast20Response, error) {
	url := fmt.Sprintf("%s/weatherapi/locationforecast/2.0/compact?lat=%s&lon=%s&altitude=%s", s.metProxyUrl, lat, lon, alt)
	log.Debugf("Getting forecast from %s", url)
	res, err := http.Get(url)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error while reading forecast from met-proxy")
	}
	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error while reading forecast from met-proxy")
	}
	res.Body.Close()

	response := LocationForecast20Response{}
	if err := json.Unmarshal(bb, &response); err != nil {
		return nil, stacktrace.Propagate(err, "error while unmarshalling forecast response")
	}

	return &response, nil
}
