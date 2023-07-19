package model

type Cart struct {
	ID     int `db:"id" json:"id"`
	UserID int `db:"user_id" json:"user_id"`
}
