package model

import "time"

type Order struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	DeliveredAt time.Time `db:"delivered_at" json:"delivered_at"`
}
