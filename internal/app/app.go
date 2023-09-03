package app

import (
	"log"
	"market/internal/config"
	"market/internal/handler"
	"market/internal/model"
	"market/internal/repository"
	"market/internal/service"
	"market/pkg/cloud"
	"market/pkg/database/postgres"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title Market API
// @version 1.0
// @description Simple market API

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Run(configDir string) {
	cfg, err := config.InitConfig(configDir)
	if err != nil {
		log.Fatal("Error occured while loading config: ", err.Error())
	}

	db, err := postgres.NewPostgresqlDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username,
		cfg.Postgres.DBName, cfg.Postgres.Password, cfg.Postgres.SSLMode)
	if err != nil {
		log.Fatal("Error occured while loading DB: ", err.Error())
	}

	cld, err := cloud.NewCloudinary(cfg.Cloudinary.Cloud, cfg.Cloudinary.Key, cfg.Cloudinary.Secret)
	if err != nil {
		log.Fatal("Error occured while loading Cloudinary: ", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, cld)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Error occured while loading zapLogger: ", err.Error())
	}
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	validate := validator.New()
	model.RegisterCustomValidations(validate)

	hand := &handler.Handler{
		Services:  services,
		Logger:    logger,
		Validator: validate,
	}

	mux := hand.InitRoutes()

	//поменять
	addr := ":" + viper.GetString("port")
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	log.Fatalln(http.ListenAndServe(addr, mux))
}
