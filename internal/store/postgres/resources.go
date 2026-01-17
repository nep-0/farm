package postgres

import (
	"farm/internal/models"
)

// Product Implementation

func (s *PostgresStore) AddProduct(p *models.Product) error {
	_, err := s.db.Exec("INSERT INTO products (id, name, description, image_url, quantity, visible) VALUES ($1, $2, $3, $4, $5, $6)",
		p.ID, p.Name, p.Description, p.ImageURL, p.Quantity, p.Visible)
	return err
}

func (s *PostgresStore) GetProduct(id string) (*models.Product, error) {
	var p models.Product
	err := s.db.QueryRow("SELECT id, name, description, image_url, quantity, visible FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.Quantity, &p.Visible)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PostgresStore) GetAllProducts(visibleOnly bool) ([]*models.Product, error) {
	query := "SELECT id, name, description, image_url, quantity, visible FROM products"
	if visibleOnly {
		query += " WHERE visible = true"
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

func (s *PostgresStore) UpdateProduct(p *models.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = $1, description = $2, image_url = $3, quantity = $4, visible = $5 WHERE id = $6",
		p.Name, p.Description, p.ImageURL, p.Quantity, p.Visible, p.ID)
	return err
}

func (s *PostgresStore) DeleteProduct(id string) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}

// Activity Implementation

func (s *PostgresStore) AddActivity(a *models.Activity) error {
	_, err := s.db.Exec("INSERT INTO activities (id, name, description, image_url, capacity, visible) VALUES ($1, $2, $3, $4, $5, $6)",
		a.ID, a.Name, a.Description, a.ImageURL, a.Capacity, a.Visible)
	return err
}

func (s *PostgresStore) GetActivity(id string) (*models.Activity, error) {
	var a models.Activity
	err := s.db.QueryRow("SELECT id, name, description, image_url, capacity, visible FROM activities WHERE id = $1", id).
		Scan(&a.ID, &a.Name, &a.Description, &a.ImageURL, &a.Capacity, &a.Visible)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *PostgresStore) GetAllActivities(visibleOnly bool) ([]*models.Activity, error) {
	query := "SELECT id, name, description, image_url, capacity, visible FROM activities"
	if visibleOnly {
		query += " WHERE visible = true"
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

func (s *PostgresStore) UpdateActivity(a *models.Activity) error {
	_, err := s.db.Exec("UPDATE activities SET name = $1, description = $2, image_url = $3, capacity = $4, visible = $5 WHERE id = $6",
		a.Name, a.Description, a.ImageURL, a.Capacity, a.Visible, a.ID)
	return err
}

func (s *PostgresStore) DeleteActivity(id string) error {
	_, err := s.db.Exec("DELETE FROM activities WHERE id = $1", id)
	return err
}
