package main

import (
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app"
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app/service"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	config := app.ReadAppConfig()

	if err := app.ConfigureLogging(config.LogLevel); err != nil {
		return err
	}

	db, err := app.ConnectToDatabase(config.DbConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	dataReceiverService := service.NewDataReceiverService(config.DataReceiverUrl)
	metService := service.NewMetService(config.MetProxyUrl)
	sensorReadingService := service.SensorReadingsService{Db: db}
	dataCollectorService := service.NewDataCollectorService(dataReceiverService, metService, sensorReadingService)

	return dataCollectorService.CollectData()
}
