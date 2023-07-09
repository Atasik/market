package model

import "time"

type Order struct {
	ID           int       `db:"id"`
	CreationDate time.Time `db:"creation_date"`
	DeliveryDate time.Time `db:"delivery_date"`
}
