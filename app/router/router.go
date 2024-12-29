package router

import (
	"gemini-poc/utils/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"go.uber.org/zap"
)

func InitRoute(
	app *fiber.App,
	conf *config.Config,
	log *zap.Logger,
) {
	app.Use(proxy.Balancer(proxy.Config{
		Servers: []string{
			conf.App.TargetHost,
		},
		Timeout: conf.App.ProxyTimeout,
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())
			return nil
		},
		ModifyResponse: func(c *fiber.Ctx) error {
			return nil
		},
	}))
}
