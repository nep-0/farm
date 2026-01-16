package sqlite

import (
	"errors"
	"farm/internal/models"
)

// Reservation Implementation

func (s *SQLiteStore) AddReservation(r *models.Reservation) error {
	_, err := s.db.Exec("INSERT INTO reservations (id, customer_id, item_id, type, priority_rank, timestamp, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		r.ID, r.CustomerID, r.ItemID, r.Type, r.PriorityRank, r.Timestamp, r.Status)
	return err
}

func (s *SQLiteStore) GetAllReservations() ([]*models.Reservation, error) {
	rows, err := s.db.Query("SELECT id, customer_id, item_id, type, priority_rank, timestamp, status FROM reservations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*models.Reservation
	for rows.Next() {
		var r models.Reservation
		if err := rows.Scan(&r.ID, &r.CustomerID, &r.ItemID, &r.Type, &r.PriorityRank, &r.Timestamp, &r.Status); err != nil {
			return nil, err
		}
		reservations = append(reservations, &r)
	}
	return reservations, nil
}

func (s *SQLiteStore) ReserveItem(r *models.Reservation) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Verify Customer
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM customers WHERE id = ?", r.CustomerID).Scan(&count)
	if err != nil || count == 0 {
		return errors.New("customer not found")
	}

	// 2. Check and Decrement Stock
	if r.Type == models.ReservationProduct {
		var qty int
		err = tx.QueryRow("SELECT quantity FROM products WHERE id = ?", r.ItemID).Scan(&qty)
		if err != nil {
			return errors.New("product not found")
		}
		if qty <= 0 {
			return errors.New("product out of stock")
		}
		if _, err := tx.Exec("UPDATE products SET quantity = quantity - 1 WHERE id = ?", r.ItemID); err != nil {
			return err
		}
	} else if r.Type == models.ReservationActivity {
		var cap int
		err = tx.QueryRow("SELECT capacity FROM activities WHERE id = ?", r.ItemID).Scan(&cap)
		if err != nil {
			return errors.New("activity not found")
		}
		if cap <= 0 {
			return errors.New("activity fully booked")
		}
		if _, err := tx.Exec("UPDATE activities SET capacity = capacity - 1 WHERE id = ?", r.ItemID); err != nil {
			return err
		}
	} else {
		return errors.New("invalid reservation type")
	}

	// 3. Create Reservation
	r.Status = "confirmed"
	_, err = tx.Exec("INSERT INTO reservations (id, customer_id, item_id, type, priority_rank, timestamp, status) VALUES (?, ?, ?, ?, ?, ?, ?)",
		r.ID, r.CustomerID, r.ItemID, r.Type, r.PriorityRank, r.Timestamp, r.Status)
	if err != nil {
		return err
	}

	return tx.Commit()
}
