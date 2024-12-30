package router

import (
	"gemini-poc/app/adapter"
	"gemini-poc/app/controller"
	"gemini-poc/app/service"
	"gemini-poc/utils/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func InitRoute(
	app *fiber.App,
	conf *config.Config,
	log *zap.Logger,
) {
	destination_adapter := adapter.NewDestinationAdapter(conf.App.DestinationHost)

	pool_service := service.NewWorkerPool(conf.App.MaxMirrorWorker, conf.App.MaxMirrorWorkerQueue, log.Named("WorkerPool"))
	pool_service.Run()

	mirror_service := service.NewMirrorService(destination_adapter, pool_service, conf.App.Mirrors, log.Named("MirrorService"))
	mirror_controller := controller.NewMirrorController(mirror_service)

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(proxy.Balancer(proxy.Config{
		Servers: []string{
			conf.App.TargetHost,
		},
		Timeout: conf.App.ProxyTimeout,
		ModifyRequest: func(c *fiber.Ctx) error {
			c.Request().Header.Add("X-Real-IP", c.IP())
			return nil
		},
		ModifyResponse: mirror_controller.MirrorRequest,
	}))
}
