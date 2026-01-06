package psql

import (
	"fmt"
	"time"

	"github.com/parxyws/nego-gin/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func NewDB(config *config.Config) (*gorm.DB, error) {

	username := config.WriteDB.User
	password := config.WriteDB.Password
	host := config.WriteDB.Host
	port := config.WriteDB.Port
	dbname := config.WriteDB.NameDB

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	readUsername := config.ReadDB.User
	readPassword := config.ReadDB.Password
	readHost := config.ReadDB.Host
	readPort := config.ReadDB.Port
	readDBName := config.ReadDB.NameDB

	readDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", readUsername, readPassword, readHost, readPort, readDBName)
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{postgres.Open(readDSN)},
		Policy:   dbresolver.RandomPolicy{},
	}).SetConnMaxIdleTime(5 * time.Minute).SetConnMaxLifetime(time.Hour).SetMaxIdleConns(10).SetMaxOpenConns(100))

	if err != nil {
		return nil, err
	}

	// Connection pool settings for master
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
