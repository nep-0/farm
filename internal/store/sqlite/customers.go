package sqlite

import (
	"farm/internal/models"
)

// Customer Implementation

func (s *SQLiteStore) AddCustomer(c *models.Customer) error {
	c.Rank = s.calculateRank(c.Credits) // Ensure rank is set correctly on creation
	_, err := s.db.Exec("INSERT INTO customers (id, email, password, salt, name, credits, rank, role) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		c.ID, c.Email, c.Password, c.Salt, c.Name, c.Credits, c.Rank, c.Role)
	return err
}

func (s *SQLiteStore) GetCustomer(id string) (*models.Customer, error) {
	var c models.Customer
	err := s.db.QueryRow("SELECT id, email, password, salt, name, credits, rank, role FROM customers WHERE id = ?", id).
		Scan(&c.ID, &c.Email, &c.Password, &c.Salt, &c.Name, &c.Credits, &c.Rank, &c.Role)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *SQLiteStore) GetCustomerByEmail(email string) (*models.Customer, error) {
	var c models.Customer
	err := s.db.QueryRow("SELECT id, email, password, salt, name, credits, rank, role FROM customers WHERE email = ?", email).
		Scan(&c.ID, &c.Email, &c.Password, &c.Salt, &c.Name, &c.Credits, &c.Rank, &c.Role)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *SQLiteStore) GetAllCustomers() ([]*models.Customer, error) {
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

func (s *SQLiteStore) UpdateCustomerCredits(id string, credits int) (*models.Customer, error) {
	rank := s.calculateRank(credits)
	_, err := s.db.Exec("UPDATE customers SET credits = ?, rank = ? WHERE id = ?", credits, rank, id)
	if err != nil {
		return nil, err
	}
	return s.GetCustomer(id)
}

func (s *SQLiteStore) UpdateCustomerRole(id string, role string) (*models.Customer, error) {
	_, err := s.db.Exec("UPDATE customers SET role = ? WHERE id = ?", role, id)
	if err != nil {
		return nil, err
	}
	return s.GetCustomer(id)
}

func (s *SQLiteStore) UpdateCustomerName(id string, name string) (*models.Customer, error) {
	_, err := s.db.Exec("UPDATE customers SET name = ? WHERE id = ?", name, id)
	if err != nil {
		return nil, err
	}
	return s.GetCustomer(id)
}
