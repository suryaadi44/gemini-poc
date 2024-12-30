package main

import (
	"context"
	"net"
	"strconv"
	"time"

	"gemini-poc/utils/config"
	"gemini-poc/utils/logger"
	"gemini-poc/utils/shutdown"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"gemini-poc/app/controller"
	"gemini-poc/app/router"
)

var SERVICE_ONE = "http://localhost:8000"
var SERVICE_TWO = "http://localhost:8001"

func main() {
	conf := config.Load("config")
	log := logger.InitLogger(conf, "REST")
	ctx, cancel := context.WithCancel(context.Background())

	local, err := time.LoadLocation(conf.App.Timezone)
	if err != nil {
		log.Fatal("[Config] Error loading timezone", zap.Error(err))
	}
	time.Local = local

	fiberzapConfig := fiberzap.Config{
		Logger: log.Named("Fiber"),
	}
	fiberLogger := fiberzap.New(fiberzapConfig)
	app := fiber.New(
		fiber.Config{
			Prefork:      conf.Server.Rest.Prefork,
			AppName:      conf.App.Service,
			ErrorHandler: controller.HandlerError,
		},
	)
	app.Use(fiberLogger)

	router.InitRoute(app, conf, log)

	listen := net.JoinHostPort(conf.Server.Rest.Host, strconv.Itoa(conf.Server.Rest.Port))
	go func() {
		if err := app.Listen(listen); err != nil {
			log.Fatal("[Server] Error running server", zap.Error(err))
		}
	}()
	log.Info("[Server] Running", zap.String("listen", listen))

	// operations to be executed on shutdown
	ops := map[string]shutdown.Operation{

		"rest-server": func(ctx context.Context) error {
			return app.Shutdown()
		},

		"background": func(ctx context.Context) error {
			cancel()
			return nil
		},
	}
	// listen for interrupt signal
	<-shutdown.GracefulShutdown(ctx, conf.App.ShutdownTimeout, ops)
}
