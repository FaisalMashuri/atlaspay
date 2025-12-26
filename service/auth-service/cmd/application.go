package cmd

import (
	"auth_service/logger"
	"auth_service/middleware"
	"context"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.opentelemetry.io/otel"
	"log"
	"time"
)

func Run(ctx context.Context) error {
	// 1️⃣ Init Logger
	logInstance := logger.InitLogger()
	logMiddleware := middleware.NewLogMiddleware(logInstance)

	// 2️⃣ Init OpenTelemetry
	shutdownOtel, err := logger.InitOtel(ctx)
	if err != nil {
		return err
	}
	defer shutdownOtel(ctx)

	// 3️⃣ Setup Fiber
	app := fiber.New(fiber.Config{})

	app.Use(otelfiber.Middleware(
		otelfiber.WithTracerProvider(otel.GetTracerProvider()),
	))

	app.Use(logMiddleware.Start())
	app.Use(recover.New())

	// 4️⃣ Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("AtlasPay OK")
	})

	// 5️⃣ Start Server (non-blocking)
	errCh := make(chan error, 1)

	go func() {
		log.Println("fiber listening on :3029")
		if err := app.Listen(":3029"); err != nil {
			errCh <- err
		}
	}()

	// 6️⃣ Graceful Shutdown Handling
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

	case err := <-errCh:

		return err

	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return app.ShutdownWithContext(shutdownCtx)
}
