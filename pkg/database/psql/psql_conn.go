package psql

import (
	"fmt"
	"time"

	"github.com/parxyws/nego-gin/config"
	applogger "github.com/parxyws/nego-gin/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func NewDB(config *config.Config) (*gorm.DB, error) {

	username := config.MasterDB.User
	password := config.MasterDB.Password
	host := config.MasterDB.Host
	port := config.MasterDB.Port
	dbname := config.MasterDB.NameDB

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)

	customLogger := applogger.NewGormLogrusLogger()
	switch config.Logger.Level {
	case "trace":
		customLogger.LogLevel = logger.Info
	case "debug":
		customLogger.LogLevel = logger.Info
	case "error":
		customLogger.LogLevel = logger.Error
	case "warn":
		customLogger.LogLevel = logger.Warn
	default:
		customLogger.LogLevel = logger.Silent
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: customLogger,
	})

	if err != nil {
		return nil, err
	}

	readUsername := config.SlaveDB.User
	readPassword := config.SlaveDB.Password
	readHost := config.SlaveDB.Host
	readPort := config.SlaveDB.Port
	readDBName := config.SlaveDB.NameDB

	readDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", readUsername, readPassword, readHost, readPort, readDBName)
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(readDSN)},
		Policy:   dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(5 * time.Minute).SetConnMaxLifetime(time.Hour).SetMaxIdleConns(10).SetMaxOpenConns(100))

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// migrate -database "postgres://postgres:postgres@localhost:5540/nego_db?sslmode=disable" -path db/migrations up
// gentool -db "postgres" -dsn "postgres://postgres:postgres@localhost:5540/nego_db?sslmode=disable" -fieldWithTypeTag true -fieldWithIndexTag true -fieldNullable true
