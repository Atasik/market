package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	CartRepo
	OrderRepo
	ProductRepo
	UserRepo
	ReviewRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		CartRepo:    NewCartPostgresqlRepo(db),
		OrderRepo:   NewOrderPostgresqlRepo(db),
		ProductRepo: NewProductPostgresqlRepo(db),
		UserRepo:    NewUserPostgresqlRepo(db),
		ReviewRepo:  NewReviewPostgresqlRepo(db),
	}
}
