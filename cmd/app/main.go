package main

import (
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app"
	"github.com/bentsolheim/kilsundvaeret-data-collector/internal/pkg/app/service"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	config := app.ReadAppConfig()

	if err := configureLogging(config.LogLevel); err != nil {
		return err
	}

	db, err := app.ConnectAndMigrateDatabase(config.DbConfig)
	if err != nil {
		return err
	}
	defer db.Close()

	dataReceiverService := service.NewDataReceiverService(config.DataReceiverUrl)
	sensorReadingService := service.SensorReadingsService{Db: db}
	dataCollectorService := service.NewDataCollectorService(dataReceiverService, sensorReadingService)

	return dataCollectorService.CollectData()
}

func configureLogging(logLevel string) error {
	log.SetFormatter(&log.TextFormatter{
		PadLevelText:    true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return stacktrace.Propagate(err, "error while parsing log level %s", logLevel)
	}
	log.SetLevel(level)
	return nil
}
