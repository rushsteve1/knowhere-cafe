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
	BindAddr string
	DisplayName string
	DomainName  string
	ShowPionts bool
	AllowRegistration bool
	AllowInvites bool
	DevEndpoints bool
	JSONRepr bool
	XMLRepr bool
	AdminUser   uint              `gorm:"references:users(id)"`
	SMTP        sql.Null[ConfigCredentials] `gorm:"embedded;embeddedPrefix:smtp_"`
	IMAP        sql.Null[ConfigCredentials] `gorm:"embedded;embeddedPrefix:imap_"`
}

type ConfigCredentials struct {
	Host     string
	User     string
	Password string
}

func DefaultConfig() Config {
	return Config{
		DisplayName: "Knowhere Cafe",
	}
}

func QueryConfig(db *gorm.DB) (cfg *Config, error) {
	res := db.Last(cfg)

	if res.RowsAffected == 0 {
		db.Create(DefaultConfig())
		return QueryConfig(db)
	}

	return cfg, res.Error
}
