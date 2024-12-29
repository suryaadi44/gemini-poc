package app

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) MirrorRequest(c *fiber.Ctx) error {
	return nil
}
