package app

import (
	"database/sql"
	"github.com/bentsolheim/go-app-utils/db"
	"github.com/palantir/stacktrace"
	"github.com/sirupsen/logrus"
)

func ConnectToDatabase(config db.DbConfig) (*sql.DB, error) {
	logrus.Info("Connecting to db: ", config.ConnectString("***"))
	db, err := sql.Open("mysql", config.ConnectString(""))
	if err != nil {
		return nil, stacktrace.Propagate(err, "error while connecting to database")
	}
	if err := db.Ping(); err != nil {
		return nil, stacktrace.Propagate(err, "unable to ping database")
	}
	return db, nil
}
