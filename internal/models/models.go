package models

import "time"

// Users is the user model
type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rooms is the room model
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions is the restriction model
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservations is the reservation model
type Reservations struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Rooms // RoomID associated with Rooms
	// RoomID associated with Rooms, so add Room.
	// By adding Room, you can include all of the room information in Reservations,
	// and use room information easily.
	// (you don't have to do, but you can do if you want to)
}

// RoomRestrictions is the room restriction model
type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Rooms        // RoomID associated with Rooms (easily access room information)
	Reservation   Reservations // ReservationID associated with Reservations (easily access reservation information)
	Restriction   Restrictions // RestrictionID associated with Restrictions (easily access restriction information)
}
