package api

import (
	"farm/internal/models"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateCredits(c echo.Context) error {
	// Only Admin (Middleware applied in routes)
	id := c.Param("id")
	type Request struct {
		Credits int `json:"credits"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	updated, err := h.store.UpdateCustomerCredits(id, req.Credits)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "customer not found"})
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *Handler) UpdateRole(c echo.Context) error {
	id := c.Param("id")
	type Request struct {
		Role string `json:"role"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Role != models.RoleAdmin && req.Role != models.RoleCustomer {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid role"})
	}

	updated, err := h.store.UpdateCustomerRole(id, req.Role)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "customer not found"})
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *Handler) CreateProduct(c echo.Context) error {
	var p models.Product
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if err := h.store.AddProduct(&p); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, p)
}

func (h *Handler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	var p models.Product
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	p.ID = id
	if err := h.store.UpdateProduct(&p); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, p)
}

func (h *Handler) CreateActivity(c echo.Context) error {
	var a models.Activity
	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if err := h.store.AddActivity(&a); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, a)
}

func (h *Handler) UpdateActivity(c echo.Context) error {
	id := c.Param("id")
	var a models.Activity
	if err := c.Bind(&a); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	a.ID = id
	if err := h.store.UpdateActivity(&a); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, a)
}

func (h *Handler) ListReservations(c echo.Context) error {
	list, err := h.store.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) ListUsers(c echo.Context) error {
	list, err := h.store.GetAllCustomers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// Sanitize passwords
	for _, u := range list {
		u.Password = ""
	}
	return c.JSON(http.StatusOK, list)
}
