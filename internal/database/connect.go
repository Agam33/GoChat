package database

import (
	"errors"
	"fmt"
	"go-chat/internal/model"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dbCfg *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", dbCfg.Host, dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Port, dbCfg.SslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
		TranslateError:                           true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		if err := sqlDB.Ping(); err == nil {
			break
		}
		log.Println("waiting for database")
		time.Sleep(2 * time.Second)

		if i == maxRetries {
			return nil, errors.New("database not ready after retries")
		}
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Room{},
		&model.UserRoom{},
		&model.Message{},
		&model.ReplyMessage{},
	); err != nil {
		return nil, err
	}
	return db, nil
}
