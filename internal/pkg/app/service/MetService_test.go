package service_test

import (
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app"
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMetService_GetMostRecentImmediateForecast(t *testing.T) {
	_ = app.ConfigureLogging("debug")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/weatherapi/locationforecast/2.0/compact?lat=58.55322&lon=8.97692&altitude=26" {
			t.Errorf("Unexpected request url %s", r.RequestURI)
		}
		file, err := ioutil.ReadFile("testdata/locationforecast_2_0_compact_altitude_26_lat_58.55322_lon_8.97692.json")
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(file); err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	metService := service.NewMetService(ts.URL, func() time.Time {
		return time.Date(2020, time.August, 19, 0, 15, 0, 0, nil)
	})
	forecast, err := metService.GetMostRecentImmediateForecast("58.55322", "8.97692", "26")
	if err != nil {
		t.Errorf("Should not fail: %s", err)
	}
	if forecast.UpdatedAt != "2020-08-18T20:58:25Z" {
		t.Errorf("Unexpected UpdatedAt %s", forecast.UpdatedAt)
	}
	temperature := forecast.Data.Instant.Details.AirTemperature
	if temperature != 19.6 {
		t.Errorf("Unexpected AirTemperature %f", temperature)
	}
	//println(fmt.Sprintf("%+v", forecast))
}
