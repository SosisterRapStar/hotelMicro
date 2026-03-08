package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/SosisterRapStar/hotels/internal/domain/hotel"
	"github.com/go-chi/chi/v5"
)

const dateLayout = "2006-01-02"

// BookingController — CRUD для бронирований отеля (без update).
type BookingController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type bookingController struct {
	mod hotel.Module
}

func NewBookingController(mod hotel.Module) BookingController {
	return &bookingController{mod: mod}
}

// Create создаёт бронирование.
//
// @Summary  Create hotel booking
// @Description Создаёт бронирование номера (check_in, check_out в формате YYYY-MM-DD)
// @Tags     bookings
// @Accept   json
// @Produce  json
// @Param    body   body  createHotelBookingRequest  true  "Данные бронирования"
// @Success  201    {object}  hotel.HotelBooking
// @Failure  400    {object}  errorResponse
// @Failure  500    {object}  errorResponse
// @Router   /bookings [post]
func (c *bookingController) Create(w http.ResponseWriter, r *http.Request) {
	var req createHotelBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	checkIn, err := time.Parse(dateLayout, req.CheckIn)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid check_in format, use YYYY-MM-DD")
		return
	}
	checkOut, err := time.Parse(dateLayout, req.CheckOut)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid check_out format, use YYYY-MM-DD")
		return
	}
	input := hotel.CreateHotelBookingInput{
		UserID:   req.UserID,
		HotelID:  req.HotelID,
		RoomID:   req.RoomID,
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Status:   hotel.BookingStatusReserved,
	}
	if req.Status != "" {
		input.Status = req.Status
	}
	created, err := c.mod.CreateHotelBooking(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// Get возвращает бронирование по ID.
//
// @Summary  Get booking by ID
// @Description Возвращает бронирование по ID
// @Tags     bookings
// @Produce  json
// @Param    id   path      string  true  "ID бронирования (UUID)"
// @Success  200  {object}  hotel.HotelBooking
// @Failure  404  {object}  errorResponse
// @Failure  500  {object}  errorResponse
// @Router   /bookings/{id} [get]
func (c *bookingController) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	item, err := c.mod.GetHotelBookingByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, hotel.ErrHotelBookingNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get booking")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// List возвращает список бронирований с пагинацией.
//
// @Summary  List hotel bookings
// @Description Список бронирований
// @Tags     bookings
// @Produce  json
// @Param    limit   query  int  false  "Лимит (default 20)"
// @Param    offset  query  int  false  "Смещение (default 0)"
// @Success  200  {array}  hotel.HotelBooking
// @Failure  500  {object}  errorResponse
// @Router   /bookings [get]
func (c *bookingController) List(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 20)
	offset := parseIntQuery(r, "offset", 0)
	params := hotel.ListHotelBookingsParams{Limit: limit, Offset: offset}
	items, err := c.mod.ListHotelBookings(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list bookings")
		return
	}
	if items == nil {
		items = []hotel.HotelBooking{}
	}
	writeJSON(w, http.StatusOK, items)
}

// Delete удаляет бронирование по ID.
//
// @Summary  Delete booking
// @Description Удаляет бронирование
// @Tags     bookings
// @Param    id  path  string  true  "ID бронирования (UUID)"
// @Success  204
// @Failure  404  {object}  errorResponse
// @Failure  500  {object}  errorResponse
// @Router   /bookings/{id} [delete]
func (c *bookingController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := c.mod.DeleteHotelBooking(r.Context(), id); err != nil {
		if errors.Is(err, hotel.ErrHotelBookingNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
