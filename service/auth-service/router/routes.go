package router

import (
	"auth_service/internal/controller"
	"github.com/gofiber/fiber/v2"
)

type IRoutes interface {
	Init()
}

type RouteParams struct {
	AuthController controller.IController
}

type routes struct {
	app  *fiber.App
	ctrl controller.IController
}

func NewRoutes(params RouteParams, app *fiber.App) IRoutes {
	return &routes{ctrl: params.AuthController, app: app}
}

func (routes *routes) Init() {
	routes.app.Post("/register", routes.ctrl.Register)
}
