package sqlite

import (
	"farm/internal/models"
)

// Product Implementation

func (s *SQLiteStore) AddProduct(p *models.Product) error {
	_, err := s.db.Exec("INSERT INTO products (id, name, description, image_url, quantity, visible) VALUES (?, ?, ?, ?, ?, ?)",
		p.ID, p.Name, p.Description, p.ImageURL, p.Quantity, p.Visible)
	return err
}

func (s *SQLiteStore) GetProduct(id string) (*models.Product, error) {
	var p models.Product
	err := s.db.QueryRow("SELECT id, name, description, image_url, quantity, visible FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.Quantity, &p.Visible)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *SQLiteStore) GetAllProducts(visibleOnly bool) ([]*models.Product, error) {
	query := "SELECT id, name, description, image_url, quantity, visible FROM products"
	if visibleOnly {
		query += " WHERE visible = 1" // SQLite stores booleans as 1/0
	}
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.Quantity, &p.Visible); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}

func (s *SQLiteStore) UpdateProduct(p *models.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, image_url = ?, quantity = ?, visible = ? WHERE id = ?",
		p.Name, p.Description, p.ImageURL, p.Quantity, p.Visible, p.ID)
	return err
}

// Activity Implementation

func (s *SQLiteStore) AddActivity(a *models.Activity) error {
	_, err := s.db.Exec("INSERT INTO activities (id, name, description, image_url, capacity, visible) VALUES (?, ?, ?, ?, ?, ?)",
		a.ID, a.Name, a.Description, a.ImageURL, a.Capacity, a.Visible)
	return err
}

func (s *SQLiteStore) GetActivity(id string) (*models.Activity, error) {
	var a models.Activity
	err := s.db.QueryRow("SELECT id, name, description, image_url, capacity, visible FROM activities WHERE id = ?", id).
		Scan(&a.ID, &a.Name, &a.Description, &a.ImageURL, &a.Capacity, &a.Visible)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *SQLiteStore) GetAllActivities(visibleOnly bool) ([]*models.Activity, error) {
	query := "SELECT id, name, description, image_url, capacity, visible FROM activities"
	if visibleOnly {
		query += " WHERE visible = 1"
	}
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*models.Activity
	for rows.Next() {
		var a models.Activity
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.ImageURL, &a.Capacity, &a.Visible); err != nil {
			return nil, err
		}
		activities = append(activities, &a)
	}
	return activities, nil
}

func (s *SQLiteStore) UpdateActivity(a *models.Activity) error {
	_, err := s.db.Exec("UPDATE activities SET name = ?, description = ?, image_url = ?, capacity = ?, visible = ? WHERE id = ?",
		a.Name, a.Description, a.ImageURL, a.Capacity, a.Visible, a.ID)
	return err
}
