package hotel

import (
	"context"
	"errors"
	"fmt"
)

const maxListLimit = 100

type Repository interface {
	CreateHotel(context.Context, CreateHotelInput) (*Hotel, error)
	GetHotelByID(context.Context, string) (*Hotel, error)
	ListHotels(context.Context, ListHotelsParams) ([]Hotel, error)
	DeleteHotel(context.Context, string) error

	CreateHotelRoom(context.Context, CreateHotelRoomInput) (*HotelRoom, error)
	GetHotelRoomByID(context.Context, string) (*HotelRoom, error)
	ListHotelRoomsByHotelID(context.Context, string, ListHotelRoomsParams) ([]HotelRoom, error)
	DeleteHotelRoom(context.Context, string) error

	CreateHotelBooking(context.Context, CreateHotelBookingInput) (*HotelBooking, error)
	GetHotelBookingByID(context.Context, string) (*HotelBooking, error)
	ListHotelBookings(context.Context, ListHotelBookingsParams) ([]HotelBooking, error)
	UpdateHotelBookingStatus(context.Context, string, string) (*HotelBooking, error)
	DeleteHotelBooking(context.Context, string) error
}

type Module interface {
	CreateHotel(context.Context, CreateHotelInput) (*Hotel, error)
	GetHotelByID(context.Context, string) (*Hotel, error)
	ListHotels(context.Context, ListHotelsParams) ([]Hotel, error)
	DeleteHotel(context.Context, string) error

	CreateHotelRoom(context.Context, CreateHotelRoomInput) (*HotelRoom, error)
	GetHotelRoomByID(context.Context, string) (*HotelRoom, error)
	ListHotelRoomsByHotelID(context.Context, string, ListHotelRoomsParams) ([]HotelRoom, error)
	DeleteHotelRoom(context.Context, string) error

	CreateHotelBooking(context.Context, CreateHotelBookingInput) (*HotelBooking, error)
	GetHotelBookingByID(context.Context, string) (*HotelBooking, error)
	ListHotelBookings(context.Context, ListHotelBookingsParams) ([]HotelBooking, error)
	UpdateHotelBookingStatus(context.Context, string, string) (*HotelBooking, error)
	DeleteHotelBooking(context.Context, string) error
}

type module struct {
	repository Repository
}

func NewModule(repository Repository) Module {
	return &module{repository: repository}
}

func (m *module) CreateHotel(ctx context.Context, input CreateHotelInput) (*Hotel, error) {
	if input.Name == "" || input.City == "" || input.Address == "" {
		return nil, errors.New("name, city and address are required")
	}
	created, err := m.repository.CreateHotel(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("creating hotel: %w", err)
	}
	return created, nil
}

func (m *module) GetHotelByID(ctx context.Context, id string) (*Hotel, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	item, err := m.repository.GetHotelByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting hotel by id: %w", err)
	}
	return item, nil
}

func (m *module) ListHotels(ctx context.Context, params ListHotelsParams) ([]Hotel, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > maxListLimit {
		params.Limit = maxListLimit
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	items, err := m.repository.ListHotels(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing hotels: %w", err)
	}
	return items, nil
}

func (m *module) DeleteHotel(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if err := m.repository.DeleteHotel(ctx, id); err != nil {
		return fmt.Errorf("deleting hotel: %w", err)
	}
	return nil
}

func (m *module) CreateHotelRoom(ctx context.Context, input CreateHotelRoomInput) (*HotelRoom, error) {
	if input.HotelID == "" || input.RoomType == "" {
		return nil, errors.New("hotel_id and room_type are required")
	}
	if input.RoomsTotal <= 0 || input.RoomsAvailable < 0 || input.RoomsAvailable > input.RoomsTotal {
		return nil, errors.New("rooms_total and rooms_available must be valid")
	}
	if input.Price < 0 {
		return nil, errors.New("price must be non-negative")
	}
	created, err := m.repository.CreateHotelRoom(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("creating hotel room: %w", err)
	}
	return created, nil
}

func (m *module) GetHotelRoomByID(ctx context.Context, id string) (*HotelRoom, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	item, err := m.repository.GetHotelRoomByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting hotel room by id: %w", err)
	}
	return item, nil
}

func (m *module) ListHotelRoomsByHotelID(ctx context.Context, hotelID string, params ListHotelRoomsParams) ([]HotelRoom, error) {
	if hotelID == "" {
		return nil, errors.New("hotel_id is required")
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > maxListLimit {
		params.Limit = maxListLimit
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	items, err := m.repository.ListHotelRoomsByHotelID(ctx, hotelID, params)
	if err != nil {
		return nil, fmt.Errorf("listing hotel rooms: %w", err)
	}
	return items, nil
}

func (m *module) DeleteHotelRoom(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if err := m.repository.DeleteHotelRoom(ctx, id); err != nil {
		return fmt.Errorf("deleting hotel room: %w", err)
	}
	return nil
}

func (m *module) CreateHotelBooking(ctx context.Context, input CreateHotelBookingInput) (*HotelBooking, error) {
	if input.UserID == "" || input.HotelID == "" || input.RoomID == "" {
		return nil, errors.New("user_id, hotel_id and room_id are required")
	}
	if !input.CheckOut.After(input.CheckIn) {
		return nil, errors.New("check_out must be after check_in")
	}
	if input.Status == "" {
		input.Status = BookingStatusReserved
	}
	if input.Status != BookingStatusReserved && input.Status != BookingStatusCancelled {
		return nil, errors.New("status must be reserved or cancelled")
	}
	created, err := m.repository.CreateHotelBooking(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("creating hotel booking: %w", err)
	}
	return created, nil
}

func (m *module) GetHotelBookingByID(ctx context.Context, id string) (*HotelBooking, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	item, err := m.repository.GetHotelBookingByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting hotel booking by id: %w", err)
	}
	return item, nil
}

func (m *module) ListHotelBookings(ctx context.Context, params ListHotelBookingsParams) ([]HotelBooking, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > maxListLimit {
		params.Limit = maxListLimit
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	items, err := m.repository.ListHotelBookings(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing hotel bookings: %w", err)
	}
	return items, nil
}

func (m *module) DeleteHotelBooking(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if err := m.repository.DeleteHotelBooking(ctx, id); err != nil {
		return fmt.Errorf("deleting hotel booking: %w", err)
	}
	return nil
}

func (m *module) UpdateHotelBookingStatus(ctx context.Context, id, status string) (*HotelBooking, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if status != BookingStatusReserved && status != BookingStatusCancelled {
		return nil, errors.New("status must be reserved or cancelled")
	}
	item, err := m.repository.UpdateHotelBookingStatus(ctx, id, status)
	if err != nil {
		return nil, fmt.Errorf("updating hotel booking status: %w", err)
	}
	return item, nil
}
