package service

import (
	"encoding/json"
	"fmt"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type ReadingsResponse struct {
	Message string
	Items   []Reading
}

type Reading struct {
	SensorName string
	Value      float32
	LocalTime  int64
	UnixTime   int64
}

func NewDataReceiverService(url string) DataReceiverService {
	return DataReceiverService{url}
}

type DataReceiverService struct {
	url string
}

func (s DataReceiverService) ReadSensorData(loggerId string) ([]Reading, error) {
	readingsUrl := fmt.Sprintf("%s/api/v1/logger/%s/readings", s.url, loggerId)
	log.Debugf("Reading data from %s", readingsUrl)
	resp, err := http.Get(readingsUrl)
	if err != nil {
		return nil, stacktrace.Propagate(err, fmt.Sprintf("failed loading sensor data from [%s]", readingsUrl))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error reading sensor data body")
	}
	log.Debug(string(body))
	response := ReadingsResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, stacktrace.Propagate(err, "unable to deserialize response body")
	}

	return response.Items, nil
}
