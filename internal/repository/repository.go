package repository

import (
	"time"

	"github.com/Toshiyana/BookingApp/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDates(roomID int, start, end time.Time) (bool, error)
}
