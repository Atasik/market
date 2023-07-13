package model

type Basket struct {
	ID     int `db:"id"`
	UserID int `db:"user_id"`
}
