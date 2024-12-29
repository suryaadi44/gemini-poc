package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type Handler struct {
	name            string
	targetHost      string
	destinationHost string
	methods         []string
	endpoint        []string
}

func NewHandler(
	name string,
	targetHost string,
	destinationHost string,
	methods []string,
	endpoint []string,
) *Handler {
	return &Handler{
		name:            name,
		targetHost:      targetHost,
		destinationHost: destinationHost,
		methods:         methods,
		endpoint:        endpoint,
	}
}

func (h *Handler) Redirect(c *fiber.Ctx) error {
	if err := proxy.Do(c, h.destinationHost); err != nil {
		return err
	}

	return nil
}
