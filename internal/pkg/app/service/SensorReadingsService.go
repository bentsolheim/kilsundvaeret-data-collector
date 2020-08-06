package service

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/palantir/stacktrace"
	log "github.com/sirupsen/logrus"
	"time"
)

type SensorReadingsService struct {
	Db *sql.DB
}

func (s SensorReadingsService) RegisterValue(loggerId string, sensorName string, createdDate time.Time, value float32) error {
	row := s.Db.QueryRow("select s.id from sensor s left join logger l on s.logger_id = l.id where l.name=? and s.name=?", loggerId, sensorName)
	var sensorId int32
	if err := row.Scan(&sensorId); err != nil {
		return stacktrace.Propagate(err, "error while finding sensorId for sensor %s for logger %s", sensorName, loggerId)
	}
	_, err := s.Db.Exec("insert into sensor_reading (created_date, sensor_id, sensor_read_date, value) values (?, ?, ?, ?)", time.Now(), sensorId, createdDate, value)
	if err != nil {
		merr, isMysqlError := err.(*mysql.MySQLError)
		// Duplicate entry (error number 1062) indicates that we have already stored this reading, which is ok.
		if !isMysqlError || (isMysqlError && merr.Number != 1062) {
			return stacktrace.Propagate(err, "error while inserting value %f for sensor %d", value, sensorId)
		}
	} else {
		log.Debugf("Registered sensor value %f for sensor %s for logger %s at %s", value, sensorName, loggerId, createdDate)
	}
	return nil
}
