package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"market/middleware"
	"market/pkg/basket"
	"market/pkg/handlers"
	"market/pkg/order"
	"market/pkg/product"
	"market/pkg/services"
	"market/pkg/session"
	"market/pkg/user"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	if err := initConfig(); err != nil {
		panic(err)
	}

	cloud := os.Getenv("cloud")
	key := os.Getenv("key")
	secret := os.Getenv("secret")
	password := os.Getenv("db_password")

	templates := template.Must(template.ParseGlob("./static/html/*"))

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.username"), password, viper.GetString("db.dbname"))
	db, err := sqlx.Connect("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	//проверка на то, что произошёл коннект
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	cld, err := cloudinary.NewFromParams(cloud, key, secret)
	if err != nil {
		panic(err)
	}

	_, err = cld.Admin.Ping(context.TODO())
	if err != nil {
		panic(err)
	}

	newCloud := services.NewImageServiceCloudinary(cld)
	userRepo := user.NewPostgresqlRepo(db)
	productRepo := product.NewPostgresqlRepo(db)
	orderRepo := order.NewPostgresqlRepo(db)
	basketRepo := basket.NewPostgresqlRepo(db)
	sessionManager := session.NewSessionsManager()
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	userHandler := &handlers.UserHandler{
		Tmpl:     templates,
		Sessions: sessionManager,
		UserRepo: userRepo,
		Logger:   logger,
	}

	productHandler := &handlers.MarketHandler{
		Tmpl:         templates,
		Sessions:     sessionManager,
		ProductRepo:  productRepo,
		OrderRepo:    orderRepo,
		BasketRepo:   basketRepo,
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
	//r.HandleFunc("/products/update/{id}", productHandler.AddProduct).Methods("POST")
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

	addr := ":" + viper.GetString("port")
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	http.ListenAndServe(addr, mux)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
