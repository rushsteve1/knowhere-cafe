package models

import (
	"io/fs"
	"log/slog"
	"reflect"

	"fmt"

	"gorm.io/gorm"

	// MigrateModels calls [[(gorm.DB).AutoMigrate]] and then
	// performs additional setup to get a newly created database to a usable state.
	"time"
)

func MigrateModels(db *gorm.DB, sqlFiles fs.ReadFileFS) error {
	var err error

	// Migrations are passed to Gorm in groups
	// groups run in order, and *should* run in order within a group
	var migrationGroups = [][]any{
		{Config{}},
		{ServerStartup{}},
		{User{}, Permissions{}},
		{Login{}, Invite{}, Group{}},
		{Post{}},
		{Report{}, Vote{}},
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

	var cfg *Config

	var migrationFunctions = []func(*gorm.DB) error{
		func(db *gorm.DB) error { cfg, err = migrateConfig(db); return err },
		func(db *gorm.DB) error { return migrateUsers(db, cfg) },
		migratePresentationsView,
	}
	for i, fn := range migrationFunctions {
		fnname := reflect.ValueOf(fn).Type().Name()
		slog.Debug("migration func", "index", i, "func", fnname)

		err := fn(db)
		if err != nil {
			return err
		}
	}

	var manualMigrations = []string{}
	for i, mmpath := range manualMigrations {
		slog.Debug("migration file", "index", i, "file", mmpath)

		qbytes, err := sqlFiles.ReadFile(mmpath)
		if err != nil {
			return err
		}
		db = db.Exec(string(qbytes))
	}

	return nil
}

// migrateConfig creates the initial default config in the DB,
// if one does not already exist
func migrateConfig(db *gorm.DB) (*Config, error) {
	var cfg Config
	res := db.Last(&cfg)
	if res.Error == gorm.ErrRecordNotFound && res.RowsAffected == 0 {
		cfg = defaultConfig()
		res = db.Create(&cfg)
	}
	return &cfg, res.Error
}

// migrateUsers creates the special root user and default permissions entries
func migrateUsers(db *gorm.DB, cfg *Config) error {
	defaultRoot := User{
		Email:      fmt.Sprintf("%s@%s", cfg.RootName, cfg.DomainName),
		LastSeenAt: time.Now(),
	}

	var root User
	res := db.Where(User{Email: defaultRoot.Email}).
		Attrs(defaultRoot).
		FirstOrCreate(&root)
	if res.Error != nil {
		return res.Error
	}

	defaultAnonPerm := defaultAnonPermissions()
	res = db.Where(Permissions{AppliesTo: defaultAnonPerm.AppliesTo}).
		Attrs(defaultAnonPerm).
		FirstOrCreate(&defaultAnonPerm)
	if res.Error != nil {
		return res.Error
	}

	defaultRegPerm := defaultRegisteredPermissions()
	res = db.Where(Permissions{AppliesTo: defaultRegPerm.AppliesTo}).
		Attrs(defaultRegPerm).
		FirstOrCreate(&defaultRegPerm)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func migratePresentationsView(db *gorm.DB) error {
	// return db.Migrator().CreateView("presentations", gorm.ViewOption{
	// 	Replace: true,
	// 	Query:   db, // TODO
	// })
	return nil
}
