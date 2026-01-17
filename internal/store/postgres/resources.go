package postgres

import (
	"farm/internal/models"
)

// Product Implementation

func (s *PostgresStore) AddProduct(p *models.Product) error {
	_, err := s.db.Exec("INSERT INTO products (id, name, quantity) VALUES ($1, $2, $3)", p.ID, p.Name, p.Quantity)
	return err
}

func (s *PostgresStore) GetProduct(id string) (*models.Product, error) {
	var p models.Product
	err := s.db.QueryRow("SELECT id, name, quantity FROM products WHERE id = $1", id).Scan(&p.ID, &p.Name, &p.Quantity)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Activity Implementation

func (s *PostgresStore) AddActivity(a *models.Activity) error {
	_, err := s.db.Exec("INSERT INTO activities (id, name, capacity) VALUES ($1, $2, $3)", a.ID, a.Name, a.Capacity)
	return err
}

func (s *PostgresStore) GetActivity(id string) (*models.Activity, error) {
	var a models.Activity
	err := s.db.QueryRow("SELECT id, name, capacity FROM activities WHERE id = $1", id).Scan(&a.ID, &a.Name, &a.Capacity)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
