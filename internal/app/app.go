package app

import (
	"context"
	"log"
	"market/internal/config"
	"market/internal/handler"
	"market/internal/model"
	"market/internal/repository"
	"market/internal/server"
	"market/internal/service"
	"market/pkg/auth"
	"market/pkg/cloud"
	"market/pkg/database/postgres"
	"market/pkg/hash"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/go-playground/validator/v10"
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

	hasher := hash.NewArgon2Hasher(cfg.Auth.Argon2.Memory, cfg.Auth.Argon2.Iterations, cfg.Auth.Argon2.SaltLength,
		cfg.Auth.Argon2.KeyLength, cfg.Auth.Argon2.Parallelism)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		log.Fatal("Error occured while creating tokenManager: ", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, cld, hasher, tokenManager)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Error occured while loading zapLogger: ", err.Error())
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	validate := validator.New()
	model.RegisterCustomValidations(validate)

	h := &handler.Handler{
		Services:     services,
		Logger:       logger,
		TokenManager: tokenManager,
		Validator:    validate,
	}

	mux := h.InitRoutes()

	srv := server.NewServer(cfg, mux)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := srv.Run(); err != nil {
			log.Println("error happened: ", err.Error())
		}
	}()

	log.Println("Application is running")

	<-quit

	log.Println("Application is shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error occurred on db connection close: %s", err.Error())
	}
}
