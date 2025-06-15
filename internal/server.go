package internal

import (
	"fmt"
	"music_catalog/internal/config"
	"music_catalog/internal/handler"
	"music_catalog/internal/helper"
	"music_catalog/internal/models"
	"music_catalog/internal/repository"
	"music_catalog/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	membershipsUc handler.MembershipsUc
	host          string
	app           *fiber.App
}

func (s *Server) initRoute() {
	rg := s.app.Group(helper.ApiGroup)

	//constructor handler
	membershipHandler := handler.NewMembershipsHandler(s.membershipsUc, rg)

	//setup routes
	membershipHandler.SetupRoutes()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.app.Listen(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, because of error %v", s.host, err.Error()))
	}
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

	err = db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		panic(fmt.Sprintf("auto migrate error: %v", err))
	}

	host := fmt.Sprintf(":%s", cfg.ApiPort)
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		// AllowCredentials: true,
	}))

	//constructor repo
	membershipRepo := repository.NewMembershipRepo(db)

	//constructor usecase
	membershipUc := usecase.NewMembershipUc(cfg, membershipRepo)

	return &Server{
		membershipsUc: membershipUc,
		host:          host,
		app:           app,
	}
}
