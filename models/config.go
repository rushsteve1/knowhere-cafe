package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type FlagConfig struct {
	DbPath string
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
	SMTP              sql.Null[ConfigCredentials] `gorm:"embedded;embeddedPrefix:smtp_"`
	IMAP              sql.Null[ConfigCredentials] `gorm:"embedded;embeddedPrefix:imap_"`
}

type ConfigCredentials struct {
	Host     string
	User     string
	Password string
}

func DefaultConfig() Config {
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
	}
}

func QueryConfig(db *gorm.DB) (*Config, error) {
	var cfg Config
	res := db.Last(&cfg)

	if res.RowsAffected == 0 {
		def := DefaultConfig()
		db.Create(&def)
		return QueryConfig(db)
	}

	return &cfg, res.Error
}
