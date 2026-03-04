package hotel

import (
	"errors"
	"time"
)

const (
	BookingStatusReserved  = "reserved"
	BookingStatusCancelled = "cancelled"
)

var (
	ErrHotelNotFound        = errors.New("hotel not found")
	ErrHotelRoomNotFound    = errors.New("hotel room not found")
	ErrHotelBookingNotFound = errors.New("hotel booking not found")
)

type Hotel struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	City      string    `db:"city" json:"city"`
	Address   string    `db:"address" json:"address"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CreateHotelInput struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}

type HotelRoom struct {
	ID              string `db:"id" json:"id"`
	HotelID         string `db:"hotel_id" json:"hotel_id"`
	RoomType        string `db:"room_type" json:"room_type"`
	RoomsTotal      int    `db:"rooms_total" json:"rooms_total"`
	RoomsAvailable  int    `db:"rooms_available" json:"rooms_available"`
	Price           int    `db:"price" json:"price"`
}

type CreateHotelRoomInput struct {
	HotelID        string `json:"hotel_id"`
	RoomType       string `json:"room_type"`
	RoomsTotal     int    `json:"rooms_total"`
	RoomsAvailable int    `json:"rooms_available"`
	Price          int    `json:"price"`
}

type HotelBooking struct {
	BookingID string    `db:"booking_id" json:"booking_id"`
	UserID    string    `db:"user_id" json:"user_id"`
	HotelID   string    `db:"hotel_id" json:"hotel_id"`
	RoomID    string    `db:"room_id" json:"room_id"`
	CheckIn   time.Time `db:"check_in" json:"check_in"`
	CheckOut  time.Time `db:"check_out" json:"check_out"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateHotelBookingInput struct {
	UserID   string    `json:"user_id"`
	HotelID  string    `json:"hotel_id"`
	RoomID   string    `json:"room_id"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
	Status   string    `json:"status"`
}

type ListHotelsParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListHotelRoomsParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListHotelBookingsParams struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
