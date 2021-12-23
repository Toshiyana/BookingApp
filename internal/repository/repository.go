package repository

import "github.com/Toshiyana/BookingApp/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) error
}
