package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SosisterRapStar/hotels/internal/domain/hotel"
	"github.com/go-chi/chi/v5"
)

// HotelController — CRUD для отелей (без update).
type HotelController interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type hotelController struct {
	mod hotel.Module
}

func NewHotelController(mod hotel.Module) HotelController {
	return &hotelController{mod: mod}
}

// Create создаёт отель.
//
// @Summary  Create hotel
// @Description Создаёт новый отель
// @Tags     hotels
// @Accept   json
// @Produce  json
// @Param    body  body  createHotelRequest  true  "Данные отеля"
// @Success  201   {object}  hotel.Hotel
// @Failure  400   {object}  errorResponse
// @Failure  500   {object}  errorResponse
// @Router   /hotels [post]
func (c *hotelController) Create(w http.ResponseWriter, r *http.Request) {
	var req createHotelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input := hotel.CreateHotelInput{
		Name:    req.Name,
		City:    req.City,
		Address: req.Address,
	}
	created, err := c.mod.CreateHotel(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// Get возвращает отель по ID.
//
// @Summary  Get hotel by ID
// @Description Возвращает отель по ID
// @Tags     hotels
// @Produce  json
// @Param    id   path      string  true  "ID отеля (UUID)"
// @Success  200  {object}  hotel.Hotel
// @Failure  404  {object}  errorResponse
// @Failure  500  {object}  errorResponse
// @Router   /hotels/{id} [get]
func (c *hotelController) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	item, err := c.mod.GetHotelByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, hotel.ErrHotelNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get hotel")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

// List возвращает список отелей с пагинацией.
//
// @Summary  List hotels
// @Description Список отелей с limit/offset
// @Tags     hotels
// @Produce  json
// @Param    limit   query  int  false  "Лимит (default 20)"
// @Param    offset  query  int  false  "Смещение (default 0)"
// @Success  200  {array}  hotel.Hotel
// @Failure  500  {object}  errorResponse
// @Router   /hotels [get]
func (c *hotelController) List(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 20)
	offset := parseIntQuery(r, "offset", 0)
	params := hotel.ListHotelsParams{Limit: limit, Offset: offset}
	items, err := c.mod.ListHotels(r.Context(), params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list hotels")
		return
	}
	if items == nil {
		items = []hotel.Hotel{}
	}
	writeJSON(w, http.StatusOK, items)
}

// Delete удаляет отель по ID.
//
// @Summary  Delete hotel
// @Description Удаляет отель. Сначала нужно удалить номера и бронирования.
// @Tags     hotels
// @Param    id  path  string  true  "ID отеля (UUID)"
// @Success  204
// @Failure  404  {object}  errorResponse
// @Failure  500  {object}  errorResponse
// @Router   /hotels/{id} [delete]
func (c *hotelController) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}
	if err := c.mod.DeleteHotel(r.Context(), id); err != nil {
		if errors.Is(err, hotel.ErrHotelNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
