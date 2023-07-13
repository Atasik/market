package model

import "time"

type Review struct {
	ID           int       `db:"id" json:"id"`
	CreationDate time.Time `db:"creation_date" json:"creation_date"`
	ProductID    int       `db:"product_id"`
	UserID       int       `db:"user_id"`
	Username     string    `db:"username" json:"username"`
	Text         string    `db:"review_text" schema:"review_text" json:"review_text"`
	Rating       int       `db:"rating" schema:"rating" json:"rating"`
}
