package responses

import "github.com/gofiber/fiber/v2"

type BaseResponse struct {
	Success  bool       `json:"success"`
	Data     any        `json:"data"`
	Error    *ErrorBody `json:"error,omitempty"`
	Metadata any        `json:"metadata,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"detail"`
}

func Success(c *fiber.Ctx, status int, data any, meta any) error {
	return c.Status(status).JSON(BaseResponse{
		Success:  true,
		Data:     data,
		Metadata: meta,
	})
}

func Error(c *fiber.Ctx, status int, code, msg string, details any) error {
	return c.Status(status).JSON(BaseResponse{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: msg,
			Details: details,
		},
	})
}
