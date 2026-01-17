package store

import "farm/internal/models"

type Repository interface {
	AddCustomer(c *models.Customer) error
	GetCustomer(id string) (*models.Customer, error)
	GetCustomerByEmail(email string) (*models.Customer, error)
	GetAllCustomers() ([]*models.Customer, error)
	UpdateCustomerCredits(id string, credits int) (*models.Customer, error)
	UpdateCustomerRole(id string, role string) (*models.Customer, error)
	UpdateCustomerName(id string, name string) (*models.Customer, error)
	AddProduct(p *models.Product) error
	GetProduct(id string) (*models.Product, error)
	GetAllProducts(visibleOnly bool) ([]*models.Product, error)
	UpdateProduct(p *models.Product) error
	AddActivity(a *models.Activity) error
	GetActivity(id string) (*models.Activity, error)
	GetAllActivities(visibleOnly bool) ([]*models.Activity, error)
	UpdateActivity(a *models.Activity) error
	AddReservation(r *models.Reservation) error
	GetAllReservations() ([]*models.Reservation, error)
	GetReservationsByCustomerID(customerID string) ([]*models.Reservation, error)
	ReserveItem(r *models.Reservation) error

	DeleteProduct(id string) error
	DeleteActivity(id string) error
	DeleteReservation(id string) error
	DeleteCustomer(id string) error
}
