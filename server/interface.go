package server

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/functionx/go-sdk/db"
)

type Server interface {
	// Init service and dynamic check config
	Init(ctx context.Context, db db.DB) error
	// Start service and keep the goroutine of the blocked
	Start(group *errgroup.Group, ctx context.Context) error
	// Close service and release resources(e.g. http connect)
	Close() error
}

type Config interface {
	Name() string
}
