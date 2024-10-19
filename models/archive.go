package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Archive struct {
	gorm.Model
	URL    string `gorm:"unique"`
	Title  string `gorm:"index"`
	Reader sql.NullString
	HTMl   sql.NullString
	Notes  []Note `gorm:"many2many:archive_notes"`
}
