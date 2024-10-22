package models

import (
	"log/slog"

	"gorm.io/gorm"
)

type FlagConfig struct {
	DbUrl   string
	Migrate bool
	Cgi     bool
	Dev     bool
}

type Config struct {
	gorm.Model
	BindAddr   string
	Public     bool
	Hostname   string
	LogLevel   slog.Level
	LogHandler string // either "text" or "json"
}

func defaultConfig() Config {
	return Config{
		BindAddr:   ":9999",
		Public:     false,
		Hostname:   "knowhere",
		LogLevel:   slog.LevelWarn,
		LogHandler: "text",
	}
}
