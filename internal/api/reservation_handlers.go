package api

import (
	"farm/internal/auth"
	"farm/internal/models"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateReservation(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JWTClaims)

	type Request struct {
		ItemID string                 `json:"item_id"`
		Type   models.ReservationType `json:"type"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Fetch customer to get rank - using ID from token
	customer, err := h.store.GetCustomer(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "customer not found"})
	}

	reservation := &models.Reservation{
		ID:           uuid.New().String(),
		CustomerID:   customer.ID,
		ItemID:       req.ItemID,
		Type:         req.Type,
		PriorityRank: customer.Rank,
		Timestamp:    time.Now(),
		Status:       "pending",
	}

	if err := h.store.ReserveItem(reservation); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, reservation)
}
