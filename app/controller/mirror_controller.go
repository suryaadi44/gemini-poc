package controller

import (
	"encoding/json"
	"gemini-poc/app/service"
	"gemini-poc/utils/dto"

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
	var baseResponse dto.BaseResponse
	err := json.Unmarshal(c.Response().Body(), &baseResponse)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
			dto.NewErrorResponse(
				fiber.StatusBadGateway,
				"Error while unmarshalling target response body, but the destination request is still sent",
				[]string{err.Error()},
			),
		)
	}

	m.service.MirrorRequest(c.Path(), c.Method(), c.Queries(), c.GetReqHeaders(), c.Body(), baseResponse.Status, c.GetRespHeaders())

	return nil
}
