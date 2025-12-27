package cmd

import (
	"auth_service/config"
	"auth_service/infrastructure/database"
	"auth_service/internal/controller"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"auth_service/logger"
	"auth_service/middleware"
	"auth_service/router"
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

	configData, errLoadConfig := config.LoadConfig()
	if errLoadConfig != nil {
		log.Fatal(errLoadConfig)
	}

	db, errConnectDatabase := database.ConnectDatabase(configData.Database)
	if errConnectDatabase != nil {
		log.Fatal(errConnectDatabase)
	}
	txManager := database.NewTxManager(db)
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, *txManager)
	ctrl := controller.NewController(svc)

	routes := router.NewRoutes(router.RouteParams{
		AuthController: ctrl,
	}, app)
	routes.Init()

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
