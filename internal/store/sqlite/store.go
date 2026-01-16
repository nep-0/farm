package sqlite

import (
	"database/sql"
	"farm/internal/config"
	"farm/internal/models"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

type SQLiteStore struct {
	db     *sql.DB
	Config *config.Config
}

func NewSQLiteStore(cfg *config.Config) (*SQLiteStore, error) {
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.ConnectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &SQLiteStore{db: db, Config: cfg}
	if err := store.InitDB(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *SQLiteStore) InitDB() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS customers (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE,
			password TEXT,
			salt TEXT,
			name TEXT,
			credits INTEGER,
			rank INTEGER,
			role TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS products (
			id TEXT PRIMARY KEY,
			name TEXT,
			quantity INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS activities (
			id TEXT PRIMARY KEY,
			name TEXT,
			capacity INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS reservations (
			id TEXT PRIMARY KEY,
			customer_id TEXT,
			item_id TEXT,
			type TEXT,
			priority_rank INTEGER,
			timestamp DATETIME,
			status TEXT
		);`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

// CalculateRank Helper
func (s *SQLiteStore) calculateRank(credits int) models.Rank {
	if credits <= s.Config.Ranks.BronzeMax {
		return models.RankBronze
	} else if credits <= s.Config.Ranks.SilverMax {
		return models.RankSilver
	} else {
		return models.RankGold
	}
}
