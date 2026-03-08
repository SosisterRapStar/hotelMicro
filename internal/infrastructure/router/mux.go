package router

import (
	chi "github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/SosisterRapStar/hotels/internal/adapter/controller"
	"github.com/SosisterRapStar/hotels/internal/config"
	_ "github.com/SosisterRapStar/hotels/internal/docs"
)

func NewMux(cfg *config.AppConfig, c *controller.Controller) chi.Router {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(c.Middleware.Logger)
		r.Use(c.Middleware.Monitoring)

		r.Get("/dummy", c.V1.Dummy.GetDummy)

		r.Route("/hotels", func(r chi.Router) {
			r.Route("/{hotelId}/rooms", func(r chi.Router) {
				r.Get("/", c.V1.Room.List)
				r.Post("/", c.V1.Room.Create)
			})
			r.Get("/", c.V1.Hotel.List)
			r.Post("/", c.V1.Hotel.Create)
			r.Get("/{id}", c.V1.Hotel.Get)
			r.Delete("/{id}", c.V1.Hotel.Delete)
		})
		r.Route("/rooms", func(r chi.Router) {
			r.Get("/{id}", c.V1.Room.Get)
			r.Delete("/{id}", c.V1.Room.Delete)
		})
		r.Route("/bookings", func(r chi.Router) {
			r.Get("/", c.V1.Booking.List)
			r.Post("/", c.V1.Booking.Create)
			r.Get("/{id}", c.V1.Booking.Get)
			r.Delete("/{id}", c.V1.Booking.Delete)
		})
	})
	return r
}
