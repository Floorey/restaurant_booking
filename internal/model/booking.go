package model

import (
	"time"
)

type Booking struct {
	ID          string    `db:"id"`
	TableID     int       `db:"table_id"`
	BookingTime time.Time `db:"booking_time"`
	Persons     int       `db:"persons"`
	Phone       string    `db:"guest_phone"`
	Email       string    `db:"guest_email"`
	Comment     string    `db:"comment"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
