package service

import (
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewDataCollectorService(dataReceiverService DataReceiverService, sensorReadingsService SensorReadingsService) DataCollectorService {
	return DataCollectorService{
		dataReceiverService:  dataReceiverService,
		sensorReadingService: sensorReadingsService,
	}
}

type DataCollectorService struct {
	dataReceiverService  DataReceiverService
	sensorReadingService SensorReadingsService
}

func (s DataCollectorService) CollectData() error {
	log.Info("Collecting data from DataReceiver")
	if err := s.CollectDataFromDataReceiver(); err != nil {
		return stacktrace.Propagate(err, "error while collecting data from DataReceiver")
	}
	log.Info("Collecting data from DataReceiver successfully completed")
	return nil
}

func (s DataCollectorService) CollectDataFromDataReceiver() error {
	loggerId := "bua"
	readings, err := s.dataReceiverService.ReadSensorData(loggerId)
	if err != nil {
		return stacktrace.Propagate(err, "unable to read sensor data for logger %s", loggerId)
	}
	log.Debug(readings)

	for _, reading := range readings {
		createdDate := time.Unix(reading.UnixTime, 0).UTC()
		err := s.sensorReadingService.RegisterValue(loggerId, reading.SensorName, createdDate, reading.Value)
		if err != nil {
			return stacktrace.Propagate(err, "error while registering value for sensor %s", reading.SensorName)
		}
	}
	return nil
}
