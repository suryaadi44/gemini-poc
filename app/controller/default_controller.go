package controller

import (
	"errors"
	"gemini-poc/utils/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func HandlerError(ctx *fiber.Ctx, err error) error {
	code, message := errorTranslator(err)

	// Return status code with error JSON
	return ctx.Status(code).JSON(
		dto.NewErrorResponse(
			code,
			message,
			[]string{message},
		),
	)
}

func errorTranslator(err error) (code int, message string) {
	message = err.Error()
	code = fiber.StatusInternalServerError

	// Status code from errors if they implement *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	switch err {
	case fasthttp.ErrTimeout:
		code = fiber.StatusGatewayTimeout
		message = http.StatusText(fiber.StatusGatewayTimeout)
	case fasthttp.ErrBodyTooLarge:
		code = fiber.StatusRequestEntityTooLarge
		message = http.StatusText(fiber.StatusRequestEntityTooLarge)
	}

	return
}
