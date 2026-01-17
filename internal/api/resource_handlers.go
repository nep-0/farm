package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) ListProducts(c echo.Context) error {
	products, err := h.store.GetAllProducts(true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) ListActivities(c echo.Context) error {
	activities, err := h.store.GetAllActivities(true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, activities)
}

func (h *Handler) ListAllProducts(c echo.Context) error {
	// Admin handler - returns all
	products, err := h.store.GetAllProducts(false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (h *Handler) ListAllActivities(c echo.Context) error {
	// Admin handler - returns all
	activities, err := h.store.GetAllActivities(false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, activities)
}
