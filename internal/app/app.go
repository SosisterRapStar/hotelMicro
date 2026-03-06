package app

import (
	"github.com/SosisterRapStar/hotels/internal/adapter/controller"
	"github.com/SosisterRapStar/hotels/internal/adapter/controller/middleware"
	v1 "github.com/SosisterRapStar/hotels/internal/adapter/controller/v1"
	"github.com/SosisterRapStar/hotels/internal/config"
)

type App struct {
	Controller *controller.Controller
}

func New(cfg *config.AppConfig) (*App, error) {
	return &App{
		Controller: &controller.Controller{
			Middleware: middleware.NewMiddleware(cfg),
			V1: v1.Controller{
				Dummy: v1.NewDummyController(),
			},
		},
	}, nil
}
