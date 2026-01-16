package config

import (
	"encoding/json"
	"os"
)

type ServerConfig struct {
	Port string `json:"port"`
}

type DatabaseConfig struct {
	Driver           string `json:"driver"`
	ConnectionString string `json:"connection_string"`
}

type RankConfig struct {
	BronzeMax int `json:"bronze_max"`
	SilverMax int `json:"silver_max"`
}

type LoggingConfig struct {
	Level    string `json:"level"`     // debug, info, warn, error
	Format   string `json:"format"`    // json, text
	Output   string `json:"output"`    // stdout, file
	FilePath string `json:"file_path"` // path to log file if output is file
}

type Config struct {
	Server    ServerConfig   `json:"server"`
	Database  DatabaseConfig `json:"database"`
	Ranks     RankConfig     `json:"ranks"`
	Logging   LoggingConfig  `json:"logging"`
	JWTSecret string         `json:"jwt_secret"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
