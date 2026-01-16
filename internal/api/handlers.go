package api

import (
	"farm/internal/auth"
	"farm/internal/config"
	"farm/internal/models"
	"farm/internal/store"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store  store.Repository
	config *config.Config
}

func NewHandler(store store.Repository, cfg *config.Config) *Handler {
	return &Handler{store: store, config: cfg}
}

// --- Middleware Helpers ---

func (h *Handler) AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JWTClaims)
		if claims.Role != models.RoleAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "admin access required"})
		}
		return next(c)
	}
}
