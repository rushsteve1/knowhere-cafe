package models

import (
	"gorm.io/gorm"
	"knowhere.cafe/src/shared/log"
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
		{User{}, Permissions{}},
		{Login{}, Invite{}, Group{}},
		{Post{}},
		{Report{}, Vote{}},
	}

	for i, group := range migrationGroups {
		log.Debug("migrating", "group", i)

		// I *believe* that Gorm runs migrations in order.
		// However I don't want to rely on that, which led to the previous
		// migrationsGroups matrix
		err = db.AutoMigrate(group...)
		if err != nil {
			return err
		}
	}

	cfg, err := migrateConfig(db)
	if err != nil {
		return err
	}

	err = migrateUsers(db, cfg)
	if err != nil {
		return err
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
	defaultRoot := defaultRootUser(cfg)
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
