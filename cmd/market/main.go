package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"market/middleware"
	"market/pkg/handler"
	"market/pkg/repository"
	"market/pkg/service"
	"market/pkg/session"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	if err := initConfig(); err != nil {
		panic(err)
	}

	templates := template.Must(template.ParseGlob("./static/html/*"))

	db, err := repository.NewPostgresqlDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("db_password"),
	})
	if err != nil {
		panic(err)
	}

	cld, err := service.NewCloudinary(service.Config{
		Cloud:  os.Getenv("cloud"),
		Key:    os.Getenv("key"),
		Secret: os.Getenv("secret"),
	})
	if err != nil {
		panic(err)
	}

	newCloud := service.NewImageServiceCloudinary(cld)
	repos := repository.NewRepository(db)
	sessionManager := session.NewSessionsManager()
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	userHandler := &handler.UserHandler{
		Tmpl:       templates,
		Sessions:   sessionManager,
		Repository: repos,
		Logger:     logger,
	}

	productHandler := &handler.MarketHandler{
		Tmpl:         templates,
		Sessions:     sessionManager,
		Repository:   repos,
		Logger:       logger,
		ImageService: newCloud,
	}

	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./static")),
	)

	r := mux.NewRouter()

	r.HandleFunc("/", productHandler.Index).Methods("GET")
	r.HandleFunc("/about", productHandler.About).Methods("GET")
	r.HandleFunc("/privacy", productHandler.Privacy).Methods("GET")
	r.HandleFunc("/history", productHandler.History).Methods("GET")

	r.HandleFunc("/products/new", productHandler.AddProductForm).Methods("GET")
	r.HandleFunc("/products/new", productHandler.AddProduct).Methods("POST")
	r.HandleFunc("/products/update/{id}", productHandler.UpdateProductForm).Methods("GET")
	r.HandleFunc("/products/update/{id}", productHandler.UpdateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productHandler.Product).Methods("PUT")
	r.HandleFunc("/products/{id}", productHandler.Product).Methods("GET")
	r.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	r.HandleFunc("/basket/{id}", productHandler.AddProductToBasket).Methods("GET")
	r.HandleFunc("/basket/{id}", productHandler.DeleteProductFromBasket).Methods("DELETE")
	r.HandleFunc("/basket", productHandler.Basket).Methods("GET")
	r.HandleFunc("/register_order", productHandler.RegisterOrder).Methods("GET")

	r.HandleFunc("/register", userHandler.Register).Methods("GET")
	r.HandleFunc("/login", userHandler.Login).Methods("GET")
	r.HandleFunc("/logout", userHandler.Logout).Methods("GET")

	r.HandleFunc("/sign_up", userHandler.SignUp).Methods("POST")
	r.HandleFunc("/sign_in", userHandler.SignIn).Methods("POST")

	r.PathPrefix("/static/").Handler(staticHandler)

	mux := middleware.Auth(sessionManager, r)
	mux = middleware.AccessLog(logger, mux)
	mux = middleware.Panic(mux)

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
