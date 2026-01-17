package postgres

import (
	"farm/internal/models"
)

// Customer Implementation

func (s *PostgresStore) AddCustomer(c *models.Customer) error {
	c.Rank = s.calculateRank(c.Credits) // Ensure rank is set correctly on creation
	_, err := s.db.Exec("INSERT INTO customers (id, email, password, salt, name, credits, rank, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		c.ID, c.Email, c.Password, c.Salt, c.Name, c.Credits, c.Rank, c.Role)
	return err
}

func (s *PostgresStore) GetCustomer(id string) (*models.Customer, error) {
	var c models.Customer
	err := s.db.QueryRow("SELECT id, email, password, salt, name, credits, rank, role FROM customers WHERE id = $1", id).
		Scan(&c.ID, &c.Email, &c.Password, &c.Salt, &c.Name, &c.Credits, &c.Rank, &c.Role)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *PostgresStore) GetCustomerByEmail(email string) (*models.Customer, error) {
	var c models.Customer
	err := s.db.QueryRow("SELECT id, email, password, salt, name, credits, rank, role FROM customers WHERE email = $1", email).
		Scan(&c.ID, &c.Email, &c.Password, &c.Salt, &c.Name, &c.Credits, &c.Rank, &c.Role)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *PostgresStore) GetAllCustomers() ([]*models.Customer, error) {
	rows, err := s.db.Query("SELECT id, email, password, salt, name, credits, rank, role FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Email, &c.Password, &c.Salt, &c.Name, &c.Credits, &c.Rank, &c.Role); err != nil {
			return nil, err
		}
		customers = append(customers, &c)
	}
	return customers, nil
}

func (s *PostgresStore) UpdateCustomerCredits(id string, credits int) (*models.Customer, error) {
	rank := s.calculateRank(credits)
	_, err := s.db.Exec("UPDATE customers SET credits = $1, rank = $2 WHERE id = $3", credits, rank, id)
	if err != nil {
		return nil, err
	}
	return s.GetCustomer(id)
}

func (s *PostgresStore) UpdateCustomerRole(id string, role string) (*models.Customer, error) {
	_, err := s.db.Exec("UPDATE customers SET role = $1 WHERE id = $2", role, id)
	if err != nil {
		return nil, err
	}
	return s.GetCustomer(id)
}

func (s *PostgresStore) UpdateCustomerName(id string, name string) (*models.Customer, error) {
	_, err := s.db.Exec("UPDATE customers SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return nil, err
	}
	return s.GetCustomer(id)
}
