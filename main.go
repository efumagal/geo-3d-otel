package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/contrib/otelfiber"

	"github.com/efumagal/geo-3d-otel/otel_instrumentation"
)

func main() {
	ctx := context.Background()

	tp, exp, err := otel_instrumentation.InitializeGlobalTracerProvider(ctx)

	// Handle shutdown to ensure all sub processes are closed correctly and telemetry is exported
	defer func() {
		_ = exp.Shutdown(ctx)
		_ = tp.Shutdown(ctx)
	}()

	if err != nil {
		log.Fatalf("failed to initialize OpenTelemetry: %e", err)
	}

	app := fiber.New()

	app.Use(otelfiber.Middleware(otelfiber.WithNext(func(c *fiber.Ctx) bool {
		return c.Path() == "/health"
	})))

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(compress.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	err = app.Listen("0.0.0.0:8099")
	if err != nil {
		log.Panic(err)
	}
}
