package app

import (
	"context"
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
	"market/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"market/docs"

	_ "github.com/lib/pq"

	"github.com/go-playground/validator/v10"
)

const (
	timeout = 5 * time.Second
	tunnel  = "qgr9d1cp-8080.euw.devtunnels.ms"
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
	// TODO: refactor logger

	// zapLogger, err := logger.NewZapLogger("zap", context.TODO())
	// if err != nil {
	// 	log.Fatalf("Error occurred while loading zapLogger: %v", err.Error())
	// 	return
	// }
	// defer func() {
	// 	if err = zapLogger.Sync(); err != nil {
	// 		zapLogger.Error("Error occurred while Sync", map[string]interface{}{"error": err.Error()})
	// 	}
	// }()

	docs.SwaggerInfo.Host = tunnel
	zapLogger := logger.NewBlobLogger()

	cfg, err := config.InitConfig(configDir)
	if err != nil {
		zapLogger.Error("Error occurred while loading config", map[string]interface{}{"error": err.Error()})
		return
	}

	db, err := postgres.NewPostgresqlDB(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User,
		cfg.Postgres.DBName, cfg.Postgres.Password, cfg.Postgres.SSLMode)
	if err != nil {
		zapLogger.Error("Error occurred while loading DB", map[string]interface{}{"error": err.Error()})
		return
	}

	cld, err := cloud.NewCloudinary(cfg.Cloudinary.Cloud, cfg.Cloudinary.Key, cfg.Cloudinary.Secret)
	if err != nil {
		zapLogger.Error("Error occurred while loading Cloudinary", map[string]interface{}{"error": err.Error()})
		return
	}

	hasher := hash.NewArgon2Hasher(cfg.Auth.Argon2.MemoryMegaBytes<<18, cfg.Auth.Argon2.Iterations, cfg.Auth.Argon2.SaltLength, //nolint:gomnd
		cfg.Auth.Argon2.KeyLength, cfg.Auth.Argon2.Parallelism)

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		zapLogger.Error("Error occurred while creating tokenManager", map[string]interface{}{"error": err.Error()})
		return
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, cld, hasher, tokenManager, cfg.Auth.JWT.AccessTokenTTL)

	validate := validator.New()
	if err = model.RegisterCustomValidations(validate); err != nil {
		zapLogger.Error("Error occurred while registering validations", map[string]interface{}{"error": err.Error()})
		return
	}

	h := ctrl.NewHandler(services, validate, zapLogger, tokenManager)

	mux := h.InitRoutes()

	srv := server.NewServer(cfg, mux)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				zapLogger.Error("Panic occurred", map[string]interface{}{"error": err})
			}
		}()
		if err := srv.Run(); err != nil {
			zapLogger.Error("Failed to start server", map[string]interface{}{"error": err.Error()})
		}
	}()

	zapLogger.Info("Application is running", nil)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	zapLogger.Info("Application is shutting down", nil)

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Error("Error occurred", map[string]interface{}{"error": err.Error()})
	}

	if err := db.Close(); err != nil {
		zapLogger.Error("Error occurred", map[string]interface{}{"error": err.Error()})
	}
}
