package controller

import (
	"auth_service/internal/model/dto/requests"
	"auth_service/internal/model/dto/responses"
	"auth_service/internal/service"
	"auth_service/internal/validation"
	"github.com/gofiber/fiber/v2"
)

type IController interface {
	Register(ctx *fiber.Ctx) error
}

type controller struct {
	svc service.IService
}

func NewController(svc service.IService) IController {
	return &controller{
		svc: svc,
	}
}

func (c *controller) Register(ctx *fiber.Ctx) error {
	var reqData requests.RegisterRequest
	if err := ctx.BodyParser(&reqData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	errs := validation.Validate(reqData)
	if len(errs) > 0 {
		return responses.Error(ctx, 400, "0400", "Error validation", errs)
	}
	err := c.svc.Register(ctx.Context(), reqData)
	if err != nil {
		return responses.Error(ctx, 400, "0001", err.Error(), nil)
	}
	return responses.Success(ctx, 200, nil, nil)
}
