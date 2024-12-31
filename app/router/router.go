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
	destinationAdapter := adapter.NewDestinationAdapter(conf.App.DestinationHost, log.Named("DestinationAdapter"))

	poolService := service.NewWorkerPool(conf.App.MaxMirrorWorker, conf.App.MaxMirrorWorkerQueue, log.Named("WorkerPool"))
	poolService.Run()

	authService := service.NewAuthService(destinationAdapter, &conf.App.Auth, log.Named("AuthService"))

	mirrorService := service.NewMirrorService(authService, destinationAdapter, poolService, conf.App.Mirrors, conf.App.MaxMirrorRetry, log.Named("MirrorService"))
	mirrorController := controller.NewMirrorController(mirrorService)

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
		ModifyResponse: mirrorController.MirrorRequest,
	}))
}
