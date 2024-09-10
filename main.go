package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"knowhere.cafe/src/mail"
	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/web"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load the flags
	fcfg := loadFlags()

	// Connect to the database
	gormCfg := gorm.Config{
		Logger: shared.GormLogger{},
	}

	db, err := gorm.Open(postgres.Open(fcfg.DbPath), &gormCfg)
	shared.LogOrDie("connect to database", err)

	// Migrate the database
	err = db.AutoMigrate(
		&models.Config{},
		&models.User{},
		&models.Permissions{},
		&models.Invite{},
		&models.Post{},
		&models.Report{},
	)
	shared.LogOrDie("migrate database", err)

	// Load config from the database
	cfg, err := models.QueryConfig(db)
	shared.LogOrDie("load config", err)

	// Connect to mail server(s)
	imapClient, err := mail.IMAPConnect(cfg.IMAP)
	shared.LogOrDie("connect to imap", err, "host", cfg.IMAP.Host)

	smtpClient, err := mail.SMTPConnect(cfg.SMTP)
	shared.LogOrDie("connect to smtp", err, "host", cfg.SMTP.Host)

	// Make the context connections object
	ctxConn := shared.ContextData{
		FlagCfg: &fcfg,
		Cfg:     cfg,
		DB:      db,
		IMAP:    imapClient,
		SMTP:    smtpClient,
	}

	// Start listening on the

	// Prep the HTTP server
	srv := http.Server{
		Addr:              cfg.BindAddr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           web.RootMux(),
		ErrorLog: slog.NewLogLogger(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}),
			slog.LevelError,
		),
		BaseContext: func(_ net.Listener) context.Context {
			ctx := context.Background()
			context.WithValue(ctx, shared.CTX_DATA_KEY, ctxConn)
			return ctx
		},
	}

	// Start the HTTP server and spin
	err = srv.ListenAndServe()
	slog.Error("http server stopped", "error", err)
}

func loadFlags() models.FlagConfig {
	cfg := models.FlagConfig{}

	// TODO flags

	flag.Parse()

	cfg.DbPath = os.Args[1]

	return cfg
}
