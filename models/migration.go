package models

import (
	"log/slog"
	"reflect"

	"gorm.io/gorm"
)

// MigrateModels calls [[(gorm.DB).AutoMigrate]] and then
// performs additional setup to get a newly created database to a usable state.
func MigrateModels(db *gorm.DB) error {
	var err error

	// Migrations are passed to Gorm in groups
	// groups run in order, and *should* run in order within a group
	var migrationGroups = [][]any{
		{Config{}},
		{ServerStartup{}},
		{WikiPage{}, Note{}},
		{Archive{}},
		{Feed{}, Entry{}},
		{SitePrefs{}, Search{}},
	}

	for i, group := range migrationGroups {
		slog.Debug("auto migration", "group", i)

		// I *believe* that Gorm runs migrations in order.
		// However I don't want to rely on that, which led to the previous
		// migrationsGroups matrix
		err = db.AutoMigrate(group...)
		if err != nil {
			return err
		}
	}

	var migrationFunctions = []func(*gorm.DB) error{
		migrateConfig,
		migrateSearch,
	}
	for i, fn := range migrationFunctions {
		fnname := reflect.ValueOf(fn).Type().Name()
		slog.Debug("migration func", "index", i, "func", fnname)

		err := db.Transaction(fn)
		if err != nil {
			return err
		}
	}

	return err
}

// migrateConfig creates the initial default config in the DB,
// if one does not already exist
func migrateConfig(db *gorm.DB) error {
	var cfg Config
	res := db.Last(&cfg)
	if res.Error == gorm.ErrRecordNotFound && res.RowsAffected == 0 {
		cfg = defaultConfig()
		res = db.Create(&cfg)
	}
	return res.Error
}

// migrateSearch alters the tables to add the search fields
func migrateSearch(db *gorm.DB) error {
	res := db.Exec(
		`ALTER TABLE archives
		ADD COLUMN IF NOT EXISTS search_vector tsvector
		GENERATED ALWAYS AS
			(to_tsvector('english', coalesce(title, '') || coalesce(reader, '')))
		STORED;`,
	)
	if res.Error != nil {
		return res.Error
	}

	res = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_archives_search ON archives USING GIN (search_vector);`,
	)
	if res.Error != nil {
		return res.Error
	}

	res = db.Exec(
		`ALTER TABLE wiki_pages
		ADD COLUMN IF NOT EXISTS search_vector tsvector
		GENERATED ALWAYS AS
			(to_tsvector('english', coalesce(title, '') || coalesce(body, '')))
		STORED;`,
	)
	if res.Error != nil {
		return res.Error
	}

	res = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_wiki_pages_search ON wiki_pages USING GIN (search_vector);`,
	)
	if res.Error != nil {
		return res.Error
	}

	res = db.Exec(
		`ALTER TABLE entries
		ADD COLUMN IF NOT EXISTS search_vector tsvector
		GENERATED ALWAYS AS
			(to_tsvector('english', coalesce(title, '') || coalesce(summary, '') || coalesce(body, '')))
		STORED;`,
	)
	if res.Error != nil {
		return res.Error
	}

	res = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_entries_search ON entries USING GIN (search_vector);`,
	)
	if res.Error != nil {
		return res.Error
	}

	res = db.Exec(
		`ALTER TABLE notes
		ADD COLUMN IF NOT EXISTS search_vector tsvector
		GENERATED ALWAYS AS
			(to_tsvector('english', coalesce(body, '')))
		STORED;`,
	)
	if res.Error != nil {
		return res.Error
	}

	res = db.Exec(
		`CREATE INDEX IF NOT EXISTS idx_notes_search ON notes USING GIN (search_vector);`,
	)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
