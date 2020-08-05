package service

import (
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"time"
)

func NewDataCollectorService(dataReceiverService DataReceiverService, metService MetService, sensorReadingsService SensorReadingsService) DataCollectorService {
	return DataCollectorService{
		dataReceiverService:  dataReceiverService,
		metService:           metService,
		sensorReadingService: sensorReadingsService,
	}
}

type DataCollectorService struct {
	dataReceiverService  DataReceiverService
	metService           MetService
	sensorReadingService SensorReadingsService
}

func (s DataCollectorService) CollectData() error {
	log.Info("Collecting data from DataReceiver")
	if err := s.CollectDataFromDataReceiver(); err != nil {
		return stacktrace.Propagate(err, "error while collecting data from DataReceiver")
	}
	log.Info("Collecting data from DataReceiver successfully completed")

	log.Info("Collecting data from met.no")
	if err := s.CollectDataFromMet(); err != nil {
		return stacktrace.Propagate(err, "error while collecting data from met.no")
	}
	log.Info("Collecting data from met.no successfully completed")
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
		if err := s.sensorReadingService.RegisterValue(loggerId, reading.SensorName, createdDate, reading.Value); err != nil {
			return stacktrace.Propagate(err, "error while registering value for sensor %s", reading.SensorName)
		}
	}
	return nil
}

func (s DataCollectorService) CollectDataFromMet() error {

	forecast, err := s.metService.GetMostRecentImmediateForecast("58.55288", "8.97572")
	if err != nil {
		return stacktrace.Propagate(err, "error while getting most recent immediate forecast")
	}
	updatedAt, err := time.Parse("2006-01-02T15:04:05Z", forecast.UpdatedAt)
	if err != nil {
		return stacktrace.Propagate(err, "unable to parse UpdatedAt date %s", forecast.UpdatedAt)
	}
	details := forecast.Details
	readings := map[string]float32{
		"air-temperature":           details.AirTemperature,
		"relative-humidity":         details.RelativeHumidity,
		"wind-from-direction":       details.WindFromDirection,
		"wind-speed":                details.WindSpeed,
		"air-pressure-at-sea-level": details.AirPressureAtSeaLevel,
	}
	for sensor, value := range readings {
		if err := s.sensorReadingService.RegisterValue("met", sensor, updatedAt, value); err != nil {
			return stacktrace.Propagate(err, "error while registering value for sensor %s for logger %s", sensor, "met")
		}
	}
	return nil
}
