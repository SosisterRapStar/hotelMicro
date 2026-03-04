package controller

import (
	"github.com/SosisterRapStar/hotels/internal/adapter/controller/middleware"
	v1 "github.com/SosisterRapStar/hotels/internal/adapter/controller/v1"
)

type Controller struct {
	Middleware middleware.Middleware
	V1         v1.Controller
}
