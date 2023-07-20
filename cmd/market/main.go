package main

import (
	"log"
	"net/http"
	"os"

	"market/pkg/handler"
	"market/pkg/model"
	"market/pkg/repository"
	"market/pkg/service"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal("Error occured while loading config: ", err.Error())
	}

	db, err := repository.NewPostgresqlDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("db_password"),
	})
	if err != nil {
		log.Fatal("Error occured while loading DB: ", err.Error())
	}

	cld, err := service.NewCloudinary(service.Config{
		Cloud:  os.Getenv("cloud"),
		Key:    os.Getenv("key"),
		Secret: os.Getenv("secret"),
	})
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
