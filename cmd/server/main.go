package main

import (
	"context"
	"farm/internal/api"
	"farm/internal/auth"
	"farm/internal/config"
	"farm/internal/logger"
	"farm/internal/store"
	"farm/internal/store/postgres"
	"farm/internal/store/sqlite"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// 2. Setup Logger
	if err := logger.Setup(&cfg.Logging); err != nil {
		slog.Error("Failed to setup logger", "error", err)
		os.Exit(1)
	}

	// 3. Init Store
	slog.Info("Connecting to database", "driver", cfg.Database.Driver, "connection_string", cfg.Database.ConnectionString)

	var s store.Repository
	var errStore error

	switch cfg.Database.Driver {
	case "sqlite":
		s, errStore = sqlite.NewSQLiteStore(cfg)
	case "postgres":
		s, errStore = postgres.NewPostgresStore(cfg)
	default:
		slog.Error("Unsupported database driver", "driver", cfg.Database.Driver)
		os.Exit(1)
	}

	if errStore != nil {
		slog.Error("Failed to connect to database", "error", errStore)
		os.Exit(1)
	}

	// 4. Init Handlers
	handler := api.NewHandler(s, cfg)

	// 5. Init Echo
	e := echo.New()
	e.HideBanner = true // Use our own logs

	// Middleware: Recovery
	e.Use(middleware.Recover())

	// Middleware: Request Logger (slog integration)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		HandleError: true, // Forward error to the global error handler
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("uri", v.URI),
				slog.String("method", v.Method),
				slog.Int("status", v.Status),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("err", v.Error.Error()))
				slog.LogAttrs(context.Background(), slog.LevelError, "http_request", attrs...)
			} else {
				slog.LogAttrs(context.Background(), slog.LevelInfo, "http_request", attrs...)
			}
			return nil
		},
	}))

	// Public Routes
	e.POST("/signup", handler.Signup)
	e.POST("/login", handler.Login)

	// Protected Routes
	// Configuration for JWT Middleware
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTClaims)
		},
		SigningKey: []byte(cfg.JWTSecret),
	}
	r := e.Group("/api")
	r.Use(echojwt.WithConfig(jwtConfig))

	r.GET("/me", handler.GetMe)
	r.PUT("/me", handler.UpdateMe)
	r.GET("/products", handler.ListProducts)
	r.GET("/activities", handler.ListActivities)
	r.POST("/reserve", handler.CreateReservation)

	// Admin Routes
	admin := r.Group("/admin")
	admin.Use(handler.AdminOnly)

	admin.POST("/products", handler.CreateProduct)
	admin.PUT("/products/:id", handler.UpdateProduct)
	admin.GET("/products", handler.ListAllProducts)
	admin.POST("/activities", handler.CreateActivity)
	admin.PUT("/activities/:id", handler.UpdateActivity)
	admin.GET("/activities", handler.ListAllActivities)
	admin.GET("/reservations", handler.ListReservations)
	admin.GET("/users", handler.ListUsers)
	admin.POST("/users/:id/credits", handler.UpdateCredits)
	admin.POST("/users/:id/role", handler.UpdateRole)

	// 5. Start Server
	slog.Info("Starting server", "port", cfg.Server.Port)
	if err := e.Start(cfg.Server.Port); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
