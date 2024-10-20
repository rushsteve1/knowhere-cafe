package main

import (
	"context"
	"expvar"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"net/http/cgi"
	"os"
	"time"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
	"knowhere.cafe/src/web"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var err error

	// Start at a high log level
	// TODO at some point of semi-stability upgrade this to a

	// Make the context data object that we will fill in as we go
	state := models.ContextState{}

	// Load the flags
	state.Flags = loadFlags()
	expvar.Publish("flags", expvar.Func(func() any { return state.Flags }))

	// Connect to the database
	gormCfg := gorm.Config{
		// TODO fix the gorm logger
		// Logger: shared.GormLogger{},
	}

	state.DB, err = gorm.Open(postgres.Open(state.Flags.DbUrl), &gormCfg)
	easy.Expect(err, "connect to database", "url", state.Flags.DbUrl)

	// Migrate the database
	err = models.MigrateModels(state.DB)
	easy.Expect(err, "migrate database", "url", state.Flags.DbUrl)

	// The migrate flag causes the program to stop here
	if state.Flags.Migrate {
		return
	}

	// Load config from the database
	cfg, err := state.Config()

	// Set the logger from the flags and config
	slog.SetDefault(setupLogger(state.Flags, cfg))

	mainCtx := context.Background()

	// Load templates
	state.Templ = models.SetupTemplates(TemplateFiles(state.Flags), state.Flags.Dev)

	// Create the context
	mainCtx = context.WithValue(
		mainCtx,
		shared.CTX_STATE_KEY,
		state,
	)

	// Prep the HTTP server
	srv := http.Server{
		Addr:              cfg.BindAddr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           web.Router(StaticFiles(state.Flags)),
		BaseContext:       func(_ net.Listener) context.Context { return mainCtx },
	}

	// Record the server startup
	startup := models.NewServerStartup(cfg)
	res := state.DB.Create(&startup)
	easy.Expect(res.Error, "record startup")
	expvar.Publish("server_startup", expvar.Func(func() any { return startup }))

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

func setupLogger(flags models.FlagConfig, cfg models.Config) *slog.Logger {
	level := cfg.LogLevel
	if flags.Dev {
		level = slog.LevelDebug
	}

	opts := slog.HandlerOptions{
		AddSource: flags.Dev,
		Level:     level,
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, &opts)
	if cfg.LogHandler == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &opts)
	}

	return slog.New(handler)
}
