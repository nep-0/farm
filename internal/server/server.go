package server

import (
	"context"
	"farm/internal/api"
	"farm/internal/auth"
	"farm/internal/config"
	"farm/internal/logger"
	"farm/internal/store"
	"farm/internal/store/postgres"
	"farm/internal/store/sqlite"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e     *echo.Echo
	cfg   *config.Config
	store store.Repository
}

func New(configPath string) (*Server, error) {
	// 1. Load Config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Setup Logger
	if err := logger.Setup(&cfg.Logging); err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
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
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	if errStore != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", errStore)
	}

	// 4. Init Handlers
	handler := api.NewHandler(s, cfg)

	// 5. Init Echo
	e := echo.New()
	e.HideBanner = true

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
	r.GET("/reservations", handler.ListMyReservations)
	r.GET("/products", handler.ListProducts)
	r.GET("/activities", handler.ListActivities)
	r.POST("/reserve", handler.CreateReservation)

	// Admin Routes
	admin := r.Group("/admin")
	admin.Use(handler.AdminOnly)

	admin.POST("/products", handler.CreateProduct)
	admin.PUT("/products/:id", handler.UpdateProduct)
	admin.DELETE("/products/:id", handler.DeleteProduct)
	admin.GET("/products", handler.ListAllProducts)
	admin.POST("/activities", handler.CreateActivity)
	admin.PUT("/activities/:id", handler.UpdateActivity)
	admin.DELETE("/activities/:id", handler.DeleteActivity)
	admin.GET("/activities", handler.ListAllActivities)
	admin.GET("/reservations", handler.ListReservations)
	admin.DELETE("/reservations/:id", handler.DeleteReservation)
	admin.GET("/users", handler.ListUsers)
	admin.DELETE("/users/:id", handler.DeleteUser)
	admin.POST("/users/:id/credits", handler.UpdateCredits)
	admin.POST("/users/:id/role", handler.UpdateRole)

	return &Server{
		e:     e,
		cfg:   cfg,
		store: s,
	}, nil
}

func (s *Server) Start() error {
	slog.Info("Starting server", "port", s.cfg.Server.Port)
	return s.e.Start(s.cfg.Server.Port)
}
