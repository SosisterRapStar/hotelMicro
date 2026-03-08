package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SosisterRapStar/hotels/internal/domain/hotel"
	"github.com/go-chi/chi/v5"
)

// RoomController — CRUD для номеров отеля (без update).
type RoomController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type roomController struct {
	mod hotel.Module
}

func NewRoomController(mod hotel.Module) RoomController {
	return &roomController{mod: mod}
}

// Create создаёт номер в отеле (hotelId в URL).
//
// @Summary  Create hotel room
// @Description Создаёт номер в указанном отеле
// @Tags     rooms
// @Accept   json
// @Produce  json
// @Param    hotelId  path      string               true   "ID отеля"
// @Param    body     body      createHotelRoomRequest  true  "Данные номера"
// @Success  201      {object}  hotel.HotelRoom
// @Failure  400      {object}  errorResponse
// @Failure  500      {object}  errorResponse
// @Router   /hotels/{hotelId}/rooms [post]
func (c *roomController) Create(w http.ResponseWriter, r *http.Request) {
	hotelID := chi.URLParam(r, "hotelId")
	if hotelID == "" {
		writeError(w, http.StatusBadRequest, "hotelId is required")
		return
	}
	var req createHotelRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input := hotel.CreateHotelRoomInput{
		HotelID:        hotelID,
		RoomType:       req.RoomType,
		RoomsTotal:     req.RoomsTotal,
		RoomsAvailable: req.RoomsAvailable,
		Price:          req.Price,
	}
	created, err := c.mod.CreateHotelRoom(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// Get возвращает номер по ID.
//
// @Summary  Get room by ID
// @Description Возвращает номер по ID
// @Tags     rooms
// @Produce  json
// @Param    id   path      string  true  "ID номера (UUID)"
// @Success  200  {object}  hotel.HotelRoom
// @Failure  404  {object}  errorResponse
// @Failure  500  {object}  errorResponse
// @Router   /rooms/{id} [get]
func (c *roomController) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	item, err := c.mod.GetHotelRoomByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, hotel.ErrHotelRoomNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get room")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// List возвращает номера отеля с пагинацией.
//
// @Summary  List hotel rooms
// @Description Список номеров отеля
// @Tags     rooms
// @Produce  json
// @Param    hotelId  path   string  true   "ID отеля"
// @Param    limit    query  int     false  "Лимит (default 20)"
// @Param    offset   query  int     false  "Смещение (default 0)"
// @Success  200  {array}  hotel.HotelRoom
// @Failure  500  {object}  errorResponse
// @Router   /hotels/{hotelId}/rooms [get]
func (c *roomController) List(w http.ResponseWriter, r *http.Request) {
	hotelID := chi.URLParam(r, "hotelId")
	if hotelID == "" {
		writeError(w, http.StatusBadRequest, "hotelId is required")
		return
	}
	limit := parseIntQuery(r, "limit", 20)
	offset := parseIntQuery(r, "offset", 0)
	params := hotel.ListHotelRoomsParams{Limit: limit, Offset: offset}
	items, err := c.mod.ListHotelRoomsByHotelID(r.Context(), hotelID, params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list rooms")
		return
	}
	if items == nil {
		items = []hotel.HotelRoom{}
	}
	writeJSON(w, http.StatusOK, items)
}

// Delete удаляет номер по ID.
//
// @Summary  Delete room
// @Description Удаляет номер. Сначала удалите бронирования этого номера.
// @Tags     rooms
// @Param    id  path  string  true  "ID номера (UUID)"
// @Success  204
// @Failure  404  {object}  errorResponse
// @Failure  500  {object}  errorResponse
// @Router   /rooms/{id} [delete]
func (c *roomController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := c.mod.DeleteHotelRoom(r.Context(), id); err != nil {
		if errors.Is(err, hotel.ErrHotelRoomNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
