package migration

import (
	"context"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/functionx/go-sdk/db"
	"github.com/functionx/go-sdk/log"
	"github.com/functionx/go-sdk/server"
)

var _ server.Server = (*Server)(nil)

type Server struct {
	logger log.Logger
	config Config
}

func NewMigration(logger log.Logger, config Config) *Server {
	return &Server{
		logger: logger.With("server", "migration"),
		config: config,
	}
}

func (s *Server) Init(ctx context.Context, db db.DB) error {
	if !s.config.Enabled {
		return nil
	}

	if err := s.config.Check(); err != nil {
		return err
	}
	s.logger.Infof("enable migration server")

	driver := db.GetDriver()
	migrationDriver, err := driver.GetMigrationsDriver()
	if err != nil {
		return errors.Wrap(err, "migration source driver error")
	}
	databaseName, dbDriver, err := driver.ToMigrateDriver(os.ExpandEnv(db.GetSource()))
	if err != nil {
		return errors.WithMessage(err, "to migrate driver error")
	}
	migrateInstance, err := migrate.NewWithInstance("httpfs", migrationDriver, databaseName, dbDriver)
	if err != nil {
		return errors.Wrap(err, "migrations new error")
	}
	defer func() {
		_, _ = migrateInstance.Close()
	}()
	if err = migrateInstance.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "migrations up error")
	}
	return nil
}

func (s *Server) Start(_ *errgroup.Group, _ context.Context) error {
	return nil
}

func (s *Server) Close() error {
	return nil
}
