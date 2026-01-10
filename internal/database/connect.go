package database

import (
	"fmt"
	"go-chat/internal/model"

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

	if err := db.AutoMigrate(
		&model.User{},
		&model.Room{},
		&model.UserRoom{},
		&model.Message{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
