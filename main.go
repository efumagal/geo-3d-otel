package main

import (
	"context"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/contrib/otelfiber"

	"github.com/efumagal/geo-3d-otel/otel_instrumentation"
	geo "github.com/emanuelef/go-geo-3d"
)

var tracer trace.Tracer

func init() {
	// Name the tracer after the package, or the service if you are in main
	tracer = otel.Tracer("github.com/efumagal/geo-3d-otel")
}

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

	app.Get("/distance", func(c *fiber.Ctx) error {
		// Create a child span
		_, childSpan := tracer.Start(c.UserContext(), "distance_computation")
		start := geo.NewCoord3d(51.39674, -0.36148, 1104.9)
		end := geo.NewCoord3d(51.38463, -0.36819, 1219.2)

		// Distance in metres between two 3D coordinates
		distance := geo.Distance3D(start, end)
		childSpan.End()
		
		return c.JSON(map[string]any{"distance": distance})
	})

	err = app.Listen("0.0.0.0:8099")
	if err != nil {
		log.Panic(err)
	}
}
