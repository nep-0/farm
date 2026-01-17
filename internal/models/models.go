package models

import (
	"time"
)

type Rank int

const (
	RankBronze Rank = iota
	RankSilver
	RankGold
)

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

func (r Rank) String() string {
	switch r {
	case RankBronze:
		return "Bronze"
	case RankSilver:
		return "Silver"
	case RankGold:
		return "Gold"
	default:
		return "Unknown"
	}
}

type Customer struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Hash
	Salt     string `json:"-"`
	Name     string `json:"name"`
	Credits  int    `json:"credits"`
	Rank     Rank   `json:"rank"`
	Role     string `json:"role"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Quantity    int    `json:"quantity"`
	Visible     bool   `json:"visible"`
}

type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Capacity    int    `json:"capacity"`
	Visible     bool   `json:"visible"`
}

type ReservationType string

const (
	ReservationProduct  ReservationType = "product"
	ReservationActivity ReservationType = "activity"
)

type Reservation struct {
	ID           string          `json:"id"`
	CustomerID   string          `json:"customer_id"`
	ItemID       string          `json:"item_id"` // ProductID or ActivityID
	Type         ReservationType `json:"type"`
	PriorityRank Rank            `json:"priority_rank"`
	Timestamp    time.Time       `json:"timestamp"`
	Status       string          `json:"status"` // "confirmed", "waitlist"
}
