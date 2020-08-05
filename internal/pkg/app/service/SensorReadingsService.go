package service

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/palantir/stacktrace"
	"time"
)

type SensorReadingsService struct {
	Db *sql.DB
}

func (s SensorReadingsService) RegisterValue(loggerId string, sensorName string, createdDate time.Time, value float32) error {
	row := s.Db.QueryRow("select s.id from sensor s left join logger l on s.logger_id = l.id where l.name=? and s.name=?", loggerId, sensorName)
	var sensorId int32
	err := row.Scan(&sensorId)
	if err != nil {
		return stacktrace.Propagate(err, "error while finding sensorId for sensor %s for logger %s", sensorName, loggerId)
	}
	_, err = s.Db.Exec("insert into sensor_reading (sensor_id, created_date, value) values (?, ?, ?)", sensorId, createdDate, value)
	if err != nil {
		merr, ok := err.(*mysql.MySQLError)
		// Duplicate entry (error number 1062) indicates that we have already stored this reading, which is ok.
		if !ok || (ok && merr.Number != 1062) {
			return stacktrace.Propagate(err, "error while inserting value %f for sensor %d", value, sensorId)
		}
	}
	return nil
}