package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	BasketRepo
	OrderRepo
	ProductRepo
	UserRepo
	ReviewRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		BasketRepo:  NewBasketPostgresqlRepo(db),
		OrderRepo:   NewOrderPostgresqlRepo(db),
		ProductRepo: NewProductPostgresqlRepo(db),
		UserRepo:    NewUserPostgresqlRepo(db),
		ReviewRepo:  NewReviewPostgresqlRepo(db),
	}
}
