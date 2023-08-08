package main

import (
	"context"
	"log"
	"math/rand"
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
	tracer = otel.Tracer("github.com/efumagal/geo-3d-otel")
}

func randFloat(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
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

	// Exclude instrumentation for /health 
	app.Use(otelfiber.Middleware(otelfiber.WithNext(func(c *fiber.Ctx) bool {
		return c.Path() == "/health"
	})))

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(compress.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	// Hello World HTML page
	app.Get("/hello", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString("<h1>Hello, World !</h1>")
	})

	// Simplified not to pass params and generates random 3D points
	app.Get("/distance", func(c *fiber.Ctx) error {
		// Create a child span to show time taken just for the computation
		_, childSpan := tracer.Start(c.UserContext(), "distance_computation")
		start := geo.NewCoord3d(randFloat(-90, 90), randFloat(-180, 180), randFloat(0, 10000))
		end := geo.NewCoord3d(randFloat(-90, 90), randFloat(-180, 180), randFloat(0, 10000))

		// Distance in metres between two 3D coordinates
		distance := geo.Distance3D(start, end)
		childSpan.End()

		return c.JSON(map[string]any{"distance": distance})
	})

	// Just to simulate a more CPU intensive load
	app.Post("/action/cpu-load", func(c *fiber.Ctx) error {
		// Create a child span
		_, childSpan := tracer.Start(c.UserContext(), "cpu_computation")

		for i := 0; i < 1_000; i++ {
			start := geo.NewCoord3d(randFloat(-90, 90), randFloat(-180, 180), randFloat(0, 10000))
			end := geo.NewCoord3d(randFloat(-90, 90), randFloat(-180, 180), randFloat(0, 10000))
			geo.Distance3D(start, end)
		}

		childSpan.End()

		return c.SendStatus(http.StatusOK)
	})

	err = app.Listen("0.0.0.0:8080")
	if err != nil {
		log.Panic(err)
	}
}
