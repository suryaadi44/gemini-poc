package controller

import (
	"gemini-poc/app/service"

	"github.com/gofiber/fiber/v2"
)

type MirrorController struct {
	service *service.MirrorService
}

func NewMirrorController(
	service *service.MirrorService,
) *MirrorController {
	return &MirrorController{
		service: service,
	}
}

func (m *MirrorController) MirrorRequest(c *fiber.Ctx) error {
	m.service.MirrorRequest(
		c.Path(),
		c.Method(),
		c.Queries(),
		c.GetReqHeaders(),
		c.Body(),
		c.Response().StatusCode(),
		c.GetRespHeaders(),
	)

	return nil
}
