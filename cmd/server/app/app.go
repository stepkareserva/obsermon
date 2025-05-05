package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/stepkareserva/obsermon/internal/server/config"
	"github.com/stepkareserva/obsermon/internal/server/db"
	"github.com/stepkareserva/obsermon/internal/server/metrics/dbstorage"
	"github.com/stepkareserva/obsermon/internal/server/metrics/handlers"
	"github.com/stepkareserva/obsermon/internal/server/metrics/memstorage"
	"github.com/stepkareserva/obsermon/internal/server/metrics/persistence"
	"github.com/stepkareserva/obsermon/internal/server/metrics/service"
	"github.com/stepkareserva/obsermon/internal/server/routing"
	"github.com/stepkareserva/obsermon/internal/server/server"
	"go.uber.org/zap"
)

type App struct {
	database *db.Database
	storage  service.Storage
	service  handlers.Service
	handler  http.Handler
	server   *server.Server
	log      *zap.Logger
}

func New(ctx context.Context, cfg config.Config, log *zap.Logger) (*App, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if log == nil {
		log = zap.NewNop()
	}

	app := App{log: log}

	if err := app.initDatabase(cfg); err != nil {
		if closeErr := app.Close(); closeErr != nil {
			log.Error("app close", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("init storage: %w", err)
	}

	if err := app.initStorage(cfg); err != nil {
		if closeErr := app.Close(); closeErr != nil {
			log.Error("app close", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("init storage: %w", err)
	}

	if err := app.initService(cfg); err != nil {
		if closeErr := app.Close(); closeErr != nil {
			log.Error("app close", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("init service: %w", err)
	}

	if err := app.initHandler(ctx, cfg); err != nil {
		if closeErr := app.Close(); closeErr != nil {
			log.Error("app close", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("init handler: %w", err)
	}

	if err := app.initServer(cfg); err != nil {
		if closeErr := app.Close(); closeErr != nil {
			log.Error("app close", zap.Error(closeErr))
		}
		return nil, fmt.Errorf("init server: %w", err)
	}

	return &app, nil
}

func (a *App) Close() error {
	if a == nil {
		return nil
	}

	var closingErrs error
	// cancel server if exists
	if a.server != nil {
		context, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.server.Shutdown(context); err != nil && !errors.Is(err, http.ErrServerClosed) {
			closingErrs = errors.Join(err, fmt.Errorf("server shutdown: %w", err))
		} else {
			a.log.Info("server stopped")
		}
		a.server = nil
	}

	// close storage, if it can be closed
	if a.storage != nil {
		if c, ok := a.storage.(io.Closer); ok {
			if err := c.Close(); err != nil {
				closingErrs = errors.Join(err, fmt.Errorf("storage closing: %w", err))
			} else {
				a.log.Info("storage closed")
			}
		} else {
			a.log.Info("storage does not implement io.Closer")
		}
		a.storage = nil
	}

	// close database if it can be closed
	if a.database != nil {
		if err := a.database.Close(); err != nil {
			closingErrs = errors.Join(err, fmt.Errorf("database closing: %w", err))
		} else {
			a.log.Info("database closed")
		}
		a.database = nil
	}

	return closingErrs
}

func (a *App) Run(ctx context.Context) error {
	if a == nil {
		return fmt.Errorf("app not exists")
	}
	if a.server == nil {
		return fmt.Errorf("server not exists")
	}

	// start server
	serverErrCh := a.server.Start()

	// wait for cancel or server error (it's critical)
	select {
	case <-ctx.Done():
		a.log.Info("app received cancel signal")
		return nil
	case srvErr := <-serverErrCh:
		return fmt.Errorf("server running: %w", srvErr)
	}
}

func (a *App) initDatabase(cfg config.Config) error {
	if cfg.DBConnection == "" {
		a.log.Info("Database connection not performed, don't use database")
		return nil
	}

	db, err := db.New(cfg.DBConnection)
	if err != nil {
		return fmt.Errorf("db connect: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}

	a.database = db

	return nil
}

func (a *App) initStorage(cfg config.Config) error {
	// omg I don't know how to refactor it.
	// we need to firstly load storage state if exists,
	// then create one of storage impl, then apply to it
	// loaded state.

	// try to load current state if passed
	var currentState *persistence.State
	if cfg.Restore {
		storage := persistence.NewJSONStateStorage(cfg.FileStoragePath)
		state, err := storage.LoadState()
		if err != nil {
			a.log.Warn("service config restoring", zap.Error(err))
		} else {
			currentState = state
			a.log.Info("restored state loaded", zap.String("location", cfg.FileStoragePath))
		}
	}

	// create database, presistent or memory storage
	if cfg.DBConnection != "" {
		// use db storage
		storage, err := dbstorage.New(a.database, a.log)
		if err != nil {
			return fmt.Errorf("init db storage: %w", err)
		}
		a.storage = storage
	} else {
		// storage
		a.storage = memstorage.New()

		// wrap onto persistent, if corresponding param passed
		if cfg.FileStoragePath != "" {
			// wrap onto persistent storage
			persistenStorage := persistence.NewJSONStateStorage(cfg.FileStoragePath)
			persistenceCfg := persistence.Config{
				StateStorage:  &persistenStorage,
				StoreInterval: cfg.StoreInterval(),
			}
			var err error
			a.storage, err = persistence.New(persistenceCfg, a.storage, a.log)
			if err != nil {
				return fmt.Errorf("persistent storage: %w", err)
			}
		}
	}

	// restore state, if state exists
	if currentState != nil {
		if err := currentState.Export(a.storage); err != nil {
			return fmt.Errorf("restoring state: %w", err)
		}
	}

	return nil
}

func (a *App) initService(cfg config.Config) error {
	// service
	service, err := service.New(a.storage)
	if err != nil {
		return fmt.Errorf("service creation: %w", err)
	}

	a.service = service

	return nil
}

func (a *App) initHandler(ctx context.Context, cfg config.Config) error {
	// create router
	router := routing.New(a.log)

	// add metrics handlers
	if err := router.AddMetricsHandlers(ctx, a.service); err != nil {
		return fmt.Errorf("metrics handler: %w", err)
	}

	// add database handler
	if err := router.AddDatabaseHandlers(ctx, a.database); err != nil {
		return fmt.Errorf("database handler: %w", err)
	}

	// set handlers
	handler, err := router.Handler()
	if err != nil {
		return fmt.Errorf("router handler: %w", err)
	}
	a.handler = handler

	return nil
}

func (a *App) initServer(cfg config.Config) error {
	a.server = server.New(cfg.Endpoint, a.handler)
	return nil
}
