package sqlite

import (
	"farm/internal/models"
)

// Product Implementation

func (s *SQLiteStore) AddProduct(p *models.Product) error {
	_, err := s.db.Exec("INSERT INTO products (id, name, quantity) VALUES (?, ?, ?)", p.ID, p.Name, p.Quantity)
	return err
}

func (s *SQLiteStore) GetProduct(id string) (*models.Product, error) {
	var p models.Product
	err := s.db.QueryRow("SELECT id, name, quantity FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Quantity)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Activity Implementation

func (s *SQLiteStore) AddActivity(a *models.Activity) error {
	_, err := s.db.Exec("INSERT INTO activities (id, name, capacity) VALUES (?, ?, ?)", a.ID, a.Name, a.Capacity)
	return err
}

func (s *SQLiteStore) GetActivity(id string) (*models.Activity, error) {
	var a models.Activity
	err := s.db.QueryRow("SELECT id, name, capacity FROM activities WHERE id = ?", id).Scan(&a.ID, &a.Name, &a.Capacity)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
