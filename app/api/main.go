package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	conf "github.com/ardanlabs/conf/v3"
	"github.com/thebigyovadiaz/golang-base-structure/app/api/handlers"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/database/postgres"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/logger"
	"github.com/thebigyovadiaz/golang-base-structure/foundation/web"
)

var build = "dev"
var env = os.Getenv("ENV")

func main() {
	var logLevel logger.Level
	level := os.Getenv("APP_LOG_LEVEL")

	switch level {
	case "INFO":
		logLevel = logger.LevelInfo
	case "DEBUG":
		logLevel = logger.LevelDebug
	case "WARN":
		logLevel = logger.LevelWarn
	case "ERROR":
		logLevel = logger.LevelError
	default:
		level = "INFO"
		logLevel = logger.LevelInfo
	}

	traceFunc := func(ctx context.Context) []any {
		v := web.GetValues(ctx)

		fields := make([]any, 2, 4)
		fields[0], fields[1] = "traceID", v.TraceID

		if v.RUT != "" {
			fields = append(fields, "RUT", v.RUT)
		}

		return fields
	}

	log := logger.New(os.Stdout, logLevel, "go-base-structure", traceFunc)

	ctx := context.Background()

	log.Info(ctx, "startup - service version...", "version", build, "cores", runtime.GOMAXPROCS(0), "logLevelAt", level)

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "service error, shutting down", "errorDetails", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	defer log.Info(ctx, "shutdown - complete")

	/*==========================================================================
		App Configuration
	==========================================================================*/

	var cfg Config

	const prefix = "APP"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup - set config to...", "config", out)

	/*==========================================================================
		Hidden App Configuration
	==========================================================================*/

	hiddenAppConfig, err := readSecretConfig()
	if err != nil {
		return fmt.Errorf("read secret config: %w", err)
	}

	/*==========================================================================
		Postgres Database Support
	==========================================================================*/
	db, err := postgres.Open(postgres.Config{
		User:            hiddenAppConfig.Postgres.User,
		Password:        hiddenAppConfig.Postgres.Password,
		Host:            hiddenAppConfig.Postgres.Host,
		Port:            hiddenAppConfig.Postgres.Port,
		Name:            hiddenAppConfig.Postgres.Name,
		MaxIdleConns:    hiddenAppConfig.Postgres.MaxIdleConns,
		MaxOpenConns:    hiddenAppConfig.Postgres.MaxOpenConns,
		IdleConnTimeout: hiddenAppConfig.Postgres.ConnMaxIdleTime,
		EnableTLS:       hiddenAppConfig.Postgres.EnableTLS,
		ApplicationName: "go-base-structure",
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	defer func() {
		log.Info(ctx, "shutdown - stopping  database support", "host", hiddenAppConfig.Postgres.Host, "port", hiddenAppConfig.Postgres.Port)
		if err := db.Close(); err != nil {
			log.Info(ctx, "shutdown - cannot stop database support gracefully", "host", hiddenAppConfig.Postgres.Host, "port", hiddenAppConfig.Postgres.Port, "error", err.Error())
			return
		}

		log.Info(ctx, "shutdown - database support stopped", "host", hiddenAppConfig.Postgres.Host, "port", hiddenAppConfig.Postgres.Port)
	}()

	/*==========================================================================
		API Service
	==========================================================================*/

	// log.Info(ctx, "startup - initializing API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(
		handlers.APIMuxConfig{
			Build:    build,
			Shutdown: shutdown,
			Log:      log,
			DB:       db,
		},
	)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup - API router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	/*==========================================================================
		Wait for shutdown signal or server error
	==========================================================================*/

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown - started", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("cannot stop server gracefully: %w", err)
		}
	}

	return nil
}
