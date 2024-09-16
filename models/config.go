package models

import (
	"database/sql"
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
	BindAddr          string
	DisplayName       string
	DomainName        string
	ShowPoints        bool
	AllowRegistration bool
	AllowInvites      bool
	DevEndpoints      bool
	FeedEndpoints     bool
	JSONRepr          bool
	XMLRepr           bool
	Lang              string
	RootName          string
	LogLevel          slog.Level
	SMTP              sql.Null[ConfigCredentials] `gorm:"embedded;embeddedPrefix:smtp_"`
	IMAP              sql.Null[ConfigCredentials] `gorm:"embedded;embeddedPrefix:imap_"`
}

type ConfigCredentials struct {
	Host     string
	User     string
	Password string
}

func defaultConfig() Config {
	return Config{
		BindAddr:          ":9999",
		DisplayName:       "Knowhere Cafe",
		DomainName:        "knowhere.cafe",
		ShowPoints:        true,
		AllowRegistration: false,
		AllowInvites:      true,
		DevEndpoints:      true,
		FeedEndpoints:     true,
		JSONRepr:          true,
		XMLRepr:           true,
		Lang:              "en-us",
		RootName:          "root",
		LogLevel:          slog.LevelWarn,
	}
}
