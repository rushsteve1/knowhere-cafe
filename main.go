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
	var err error

	// Make the context data object that we will fill in as we go
	state := shared.ContextState{}

	// Load the flags
	state.FlagCfg = loadFlags()

	// Connect to the database
	gormCfg := gorm.Config{
		Logger: shared.GormLogger{},
	}

	state.DB, err = gorm.Open(postgres.Open(state.FlagCfg.DbPath), &gormCfg)
	shared.LogOrDie("connect to database", err)

	// Migrate the database
	err = state.DB.AutoMigrate(
		&models.Config{},
		&models.User{},
		&models.Permissions{},
		&models.Invite{},
		&models.Post{},
		&models.Report{},
	)
	shared.LogOrDie("migrate database", err)

	// Load config from the database
	state.Cfg, err = models.QueryConfig(state.DB)
	shared.LogOrDie("load config", err)

	// Connect to mail server(s)
	if state.Cfg.IMAP.Valid {
		state.IMAP, err = mail.IMAPConnect(state.Cfg.IMAP.V)
		shared.LogOrDie("connect to imap", err, "host", state.Cfg.IMAP.V.Host)
	}

	if state.Cfg.SMTP.Valid {
		state.SMTP, err = mail.SMTPConnect(state.Cfg.SMTP.V)
		shared.LogOrDie("connect to smtp", err, "host", state.Cfg.SMTP.V.Host)
	}

	// Start listening on the

	// Prep the HTTP server
	srv := http.Server{
		Addr:              state.Cfg.BindAddr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           web.RootHandler(),
		ErrorLog: slog.NewLogLogger(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}),
			slog.LevelError,
		),
		BaseContext: func(_ net.Listener) context.Context {
			ctx := context.Background()
			context.WithValue(ctx, shared.CTX_STATE_KEY, state)
			return ctx
		},
	}

	slog.Info("starting http server", "addr", state.Cfg.BindAddr)

	// Start the HTTP server and spin
	err = srv.ListenAndServe()
	slog.Error("http server stopped", "error", err)
}

func loadFlags() *models.FlagConfig {
	cfg := models.FlagConfig{}

	// TODO flags

	flag.Parse()

	cfg.DbPath = os.Args[1]

	return &cfg
}
