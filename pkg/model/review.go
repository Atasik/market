package model

import "time"

type Review struct {
	ID           int       `db:"id"`
	CreationDate time.Time `schema:"creation_date" db:"creation_date"`
	ProductID    int       `db:"product_id"`
	UserID       int       `db:"user_id"`
	Username     string    `db:"username"`
	Text         string    `db:"review_text"`
	Rating       int       `db:"rating"`
}
