package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/SosisterRapStar/hotels/internal/domain/hotel"
)

type HotelRepository struct {
	db *sqlx.DB
}

func NewHotelRepository(db *sqlx.DB, manager *Manager) *HotelRepository {
	return &HotelRepository{
		db: db,
	}
}

func (r *HotelRepository) queryer(ctx context.Context) sqlx.QueryerContext {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *HotelRepository) execer(ctx context.Context) sqlx.ExecerContext {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *HotelRepository) CreateHotel(ctx context.Context, input hotel.CreateHotelInput) (*hotel.Hotel, error) {
	id := uuid.New().String()
	execer := r.execer(ctx)
	if _, err := execer.ExecContext(ctx, createHotelQuery, id, input.Name, input.City, input.Address); err != nil {
		return nil, fmt.Errorf("executing create hotel query: %w", err)
	}
	row := hotel.Hotel{}
	queryer := r.queryer(ctx)
	if err := sqlx.GetContext(ctx, queryer, &row, getHotelByIDQuery, id); err != nil {
		return nil, fmt.Errorf("selecting created hotel: %w", err)
	}
	return &row, nil
}

func (r *HotelRepository) GetHotelByID(ctx context.Context, id string) (*hotel.Hotel, error) {
	row := hotel.Hotel{}
	queryer := r.queryer(ctx)
	if err := sqlx.GetContext(ctx, queryer, &row, getHotelByIDQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, hotel.ErrHotelNotFound
		}
		return nil, fmt.Errorf("executing get hotel by id query: %w", err)
	}
	return &row, nil
}

func (r *HotelRepository) ListHotels(ctx context.Context, params hotel.ListHotelsParams) ([]hotel.Hotel, error) {
	rows := make([]hotel.Hotel, 0, params.Limit)
	queryer := r.queryer(ctx)
	if err := sqlx.SelectContext(ctx, queryer, &rows, listHotelsQuery, params.Limit, params.Offset); err != nil {
		return nil, fmt.Errorf("executing list hotels query: %w", err)
	}
	return rows, nil
}

func (r *HotelRepository) CreateHotelRoom(ctx context.Context, input hotel.CreateHotelRoomInput) (*hotel.HotelRoom, error) {
	id := uuid.New().String()
	execer := r.execer(ctx)
	if _, err := execer.ExecContext(ctx, createHotelRoomQuery,
		id, input.HotelID, input.RoomType, input.RoomsTotal, input.RoomsAvailable, input.Price); err != nil {
		return nil, fmt.Errorf("executing create hotel room query: %w", err)
	}
	row := hotel.HotelRoom{}
	queryer := r.queryer(ctx)
	if err := sqlx.GetContext(ctx, queryer, &row, getHotelRoomByIDQuery, id); err != nil {
		return nil, fmt.Errorf("selecting created hotel room: %w", err)
	}
	return &row, nil
}

func (r *HotelRepository) GetHotelRoomByID(ctx context.Context, id string) (*hotel.HotelRoom, error) {
	row := hotel.HotelRoom{}
	queryer := r.queryer(ctx)
	if err := sqlx.GetContext(ctx, queryer, &row, getHotelRoomByIDQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, hotel.ErrHotelRoomNotFound
		}
		return nil, fmt.Errorf("executing get hotel room by id query: %w", err)
	}
	return &row, nil
}

func (r *HotelRepository) ListHotelRoomsByHotelID(ctx context.Context, hotelID string, params hotel.ListHotelRoomsParams) ([]hotel.HotelRoom, error) {
	rows := make([]hotel.HotelRoom, 0, params.Limit)
	queryer := r.queryer(ctx)
	if err := sqlx.SelectContext(ctx, queryer, &rows, listHotelRoomsByHotelIDQuery, hotelID, params.Limit, params.Offset); err != nil {
		return nil, fmt.Errorf("executing list hotel rooms query: %w", err)
	}
	return rows, nil
}

func (r *HotelRepository) CreateHotelBooking(ctx context.Context, input hotel.CreateHotelBookingInput) (*hotel.HotelBooking, error) {
	bookingID := uuid.New().String()
	execer := r.execer(ctx)
	if _, err := execer.ExecContext(ctx, createHotelBookingQuery,
		bookingID, input.UserID, input.HotelID, input.RoomID, input.CheckIn, input.CheckOut, input.Status); err != nil {
		return nil, fmt.Errorf("executing create hotel booking query: %w", err)
	}
	row := hotel.HotelBooking{}
	queryer := r.queryer(ctx)
	if err := sqlx.GetContext(ctx, queryer, &row, getHotelBookingByIDQuery, bookingID); err != nil {
		return nil, fmt.Errorf("selecting created hotel booking: %w", err)
	}
	return &row, nil
}

func (r *HotelRepository) GetHotelBookingByID(ctx context.Context, id string) (*hotel.HotelBooking, error) {
	row := hotel.HotelBooking{}
	queryer := r.queryer(ctx)
	if err := sqlx.GetContext(ctx, queryer, &row, getHotelBookingByIDQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, hotel.ErrHotelBookingNotFound
		}
		return nil, fmt.Errorf("executing get hotel booking by id query: %w", err)
	}
	return &row, nil
}

func (r *HotelRepository) ListHotelBookings(ctx context.Context, params hotel.ListHotelBookingsParams) ([]hotel.HotelBooking, error) {
	rows := make([]hotel.HotelBooking, 0, params.Limit)
	queryer := r.queryer(ctx)
	if err := sqlx.SelectContext(ctx, queryer, &rows, listHotelBookingsQuery, params.Limit, params.Offset); err != nil {
		return nil, fmt.Errorf("executing list hotel bookings query: %w", err)
	}
	return rows, nil
}

func (r *HotelRepository) UpdateHotelBookingStatus(ctx context.Context, id, status string) (*hotel.HotelBooking, error) {
	execer := r.execer(ctx)
	result, err := execer.ExecContext(ctx, updateHotelBookingStatusQuery, status, id)
	if err != nil {
		return nil, fmt.Errorf("executing update hotel booking status query: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("reading rows affected: %w", err)
	}
	if affected == 0 {
		return nil, hotel.ErrHotelBookingNotFound
	}
	row := hotel.HotelBooking{}
	queryer := r.queryer(ctx)
	if err := sqlx.GetContext(ctx, queryer, &row, getHotelBookingByIDQuery, id); err != nil {
		return nil, fmt.Errorf("selecting updated hotel booking: %w", err)
	}
	return &row, nil
}
