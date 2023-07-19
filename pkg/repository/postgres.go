package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

const (
	usersTable          = "users"
	productsTable       = "products"
	ordersTable         = "orders"
	reviewsTable        = "reviews"
	cartsTable          = "carts"
	productsUsersTable  = "products_users"
	productsOrdersTable = "products_orders"
	productsCartsTable  = "products_carts"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresqlDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ParsePostgresError(err error) error {
	if err == nil {
		return nil
	}

	pgErr, ok := err.(*pq.Error)
	if ok {
		if pgErr.Code == "23505" {
			return ErrAlreadyExists
		}
	}

	if err == sql.ErrNoRows {
		return ErrNotFound
	}

	return err
}
