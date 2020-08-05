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
		if r.RequestURI != "/api/v1/met/location-forecast?lat=58.55288&lon=8.97572" {
			t.Errorf("Unexpected request url %s", r.RequestURI)
		}
		file, err := ioutil.ReadFile("testdata/locationforecast_2_0_complete_lat_58.55288_lon_8.97572.json")
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(file); err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	metService := service.NewMetService(ts.URL)
	forecast, err := metService.GetMostRecentImmediateForecast("58.55288", "8.97572")
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
