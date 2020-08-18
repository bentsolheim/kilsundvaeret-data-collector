package service_test

import (
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app"
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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

	metService := service.NewMetService(ts.URL)
	forecast, err := metService.GetMostRecentImmediateForecast("58.55322", "8.97692", "26")
	if err != nil {
		t.Errorf("Should not fail: %s", err)
	}
	if forecast.UpdatedAt != "2020-07-30T08:54:10Z" {
		t.Errorf("Unexpected UpdatedAt %s", forecast.UpdatedAt)
	}
	if forecast.Details.AirTemperature != 18.7 {
		t.Errorf("Unexpected AirTemperature %f", forecast.Details.AirTemperature)
	}
	//println(fmt.Sprintf("%+v", forecast))
}
