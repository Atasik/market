package app

import (
	"context"
	"log"
	"market/internal/config"
	ctrl "market/internal/controller/http"
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
	"time"

	_ "github.com/lib/pq"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

const timeout = 5 * time.Second

// @title Market API
// @version 1.0
// @description Simple market API

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Run(configDir string) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Error occurred while loading zapLogger: %s\n", err.Error())
		return
	}
	defer zapLogger.Sync() //nolint:errcheck
	logger := zapLogger.Sugar()

	cfg, err := config.InitConfig(configDir)
	if err != nil {
		logger.Errorf("Error occurred while loading config: %s\n", err.Error())
		return
	}

	db, err := postgres.NewPostgresqlDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User,
		cfg.Postgres.DBName, cfg.Postgres.Password, cfg.Postgres.SSLMode)
	if err != nil {
		logger.Errorf("Error occurred while loading DB: %s\n", err.Error())
		return
	}

	cld, err := cloud.NewCloudinary(cfg.Cloudinary.Cloud, cfg.Cloudinary.Key, cfg.Cloudinary.Secret)
	if err != nil {
		logger.Errorf("Error occurred while loading Cloudinary: %s\n", err.Error())
		return
	}

	hasher := hash.NewArgon2Hasher(cfg.Auth.Argon2.MemoryMegaBytes<<18, cfg.Auth.Argon2.Iterations, cfg.Auth.Argon2.SaltLength, //nolint:gomnd
		cfg.Auth.Argon2.KeyLength, cfg.Auth.Argon2.Parallelism)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Errorf("Error occurred while creating tokenManager: %s\n", err.Error())
		return
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, cld, hasher, tokenManager, cfg.Auth.JWT.AccessTokenTTL)

	validate := validator.New()
	if err = model.RegisterCustomValidations(validate); err != nil {
		logger.Errorf("Error occurred while registering validations: %s\n", err.Error())
		return
	}

	h := ctrl.NewHandler(services, validate, logger, tokenManager)

	mux := h.InitRoutes()

	srv := server.NewServer(cfg, mux)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := srv.Run(); err != nil {
			logger.Errorf("Failed to start server: %s\n", err.Error())
		}
	}()

	logger.Info("Application is running")

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	logger.Info("Application is shutting down")

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Error(err.Error())
	}
}
