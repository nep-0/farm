package api

import (
	"database/sql"
	"farm/internal/auth"
	"farm/internal/models"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Signup(c echo.Context) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Generate a unique salt for the user
	salt, err := auth.GenerateSalt()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error generating salt"})
	}

	// Hash password with Salt (and still can use Pepper if desired, but user asked for "different salt for each user")
	// Combination: password + salt
	hash, err := auth.HashPassword(req.Password + salt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error processing password"})
	}

	customer := &models.Customer{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: hash,
		Salt:     salt,
		Name:     req.Name,
		Credits:  0,
		Role:     models.RoleCustomer, // Default role
	}

	if err := h.store.AddCustomer(customer); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "user likely already exists"})
	}

	// Don't return sensitive info
	customer.Password = ""
	customer.Salt = ""
	return c.JSON(http.StatusCreated, customer)
}

func (h *Handler) Login(c echo.Context) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	customer, err := h.store.GetCustomerByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server error"})
	}

	// Verify password + salt
	if !auth.CheckPasswordHash(req.Password+customer.Salt, customer.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	token, err := auth.GenerateToken(customer.ID, customer.Role, h.config.JWTSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) GetMe(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JWTClaims)

	customer, err := h.store.GetCustomer(claims.UserID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	// Sanitize
	customer.Password = ""
	customer.Salt = ""
	return c.JSON(http.StatusOK, customer)
}

func (h *Handler) UpdateMe(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JWTClaims)

	type Request struct {
		Name string `json:"name"`
	}
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name cannot be empty"})
	}

	updated, err := h.store.UpdateCustomerName(claims.UserID, req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not update user"})
	}

	updated.Password = ""
	updated.Salt = ""
	return c.JSON(http.StatusOK, updated)
}
