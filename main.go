package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"net/http/cgi"
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

	// Start at a high log level
	// TODO at some point of semi-stability upgrade this to a
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Make the context data object that we will fill in as we go
	state := models.ContextState{}

	// Load the flags
	state.Flags = loadFlags()

	// Connect to the database
	gormCfg := gorm.Config{
		// TODO fix the gorm logger
		// Logger:                           shared.GormLogger{},
	}

	state.DB, err = gorm.Open(postgres.Open(state.Flags.DbUrl), &gormCfg)
	shared.LogOrPanic("connect to database", err, "url", state.Flags.DbUrl)

	// Migrate the database
	err = models.MigrateModels(state.DB)
	shared.LogOrPanic("migrate database", err, "url", state.Flags.DbUrl)

	// The migrate flag causes the program to stop here
	if state.Flags.Migrate {
		return
	}

	// Load config from the database
	cfg, err := state.Config()

	// Set the log level from the config
	if !state.Flags.Dev {
		// TODO don't warn if the level doesn't change
		slog.Warn("setting log level", "level", cfg.LogLevel.String())
	}

	// Connect to mail server(s)
	if cfg.IMAP.Valid {
		state.IMAP, err = mail.IMAPConnect(cfg.IMAP.V)
		shared.LogOrPanic("connect to imap", err, "host", cfg.IMAP.V.Host)
	}

	if cfg.SMTP.Valid {
		state.SMTP, err = mail.SMTPConnect(cfg.SMTP.V)
		shared.LogOrPanic("connect to smtp", err, "host", cfg.SMTP.V.Host)
	}

	// Create the context
	mainCtx := context.WithValue(
		context.Background(),
		shared.CTX_STATE_KEY,
		state,
	)
	
	// Load templates
	state.Templ = models.SetupTemplates(web.TemplateFiles(mainCtx))

	// Prep the HTTP server
	srv := http.Server{
		Addr:              cfg.BindAddr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           web.RootHandler(),
		ErrorLog: slog.NewLogLogger(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}),
			slog.LevelError,
		),
		BaseContext: func(_ net.Listener) context.Context { return mainCtx },
	}

	// Record the server startup
	res := state.DB.Create(models.NewServerStartup(cfg))
	shared.LogOrPanic("record startup", res.Error)

	// Start the HTTP server
	slog.Info("starting http server", "addr", cfg.BindAddr)

	if state.Flags.Cgi {
		// Start as CGI
		cgi.Serve(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r = r.WithContext(srv.BaseContext(nil))
				srv.Handler.ServeHTTP(w, r)
			}),
		)
	} else {
		// Start the HTTP server and spin
		err = srv.ListenAndServe()
		slog.Error("http server stopped", "error", err)
	}
}

// loadFlags loads the CLI flags using the standard [[flag]] package
func loadFlags() models.FlagConfig {
	cfg := models.FlagConfig{}

	flag.BoolVar(&cfg.Migrate, "migrate", false, "migrate database then exit")
	flag.BoolVar(&cfg.Cgi, "cgi", false, "run in Common Gateway Interface mode")
	flag.BoolVar(&cfg.Dev, "dev", false, "enable developer mode")
	// TODO more flags

	flag.Parse()

	cfg.DbUrl = flag.Arg(0)

	return cfg
}
