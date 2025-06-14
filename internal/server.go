package internal

import (
	"fmt"
	"music_catalog/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	host string
}

func NewServer() *Server {
	cfg := config.NewConfig()
	//postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance")
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLife) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Minute)

	db = db.Debug()

	host := fmt.Sprintf(":%s", cfg.ApiPort)

	return &Server{
		host: host,
	}
}
