package pprof_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/functionx/go-sdk/log"
	"github.com/functionx/go-sdk/pprof"
)

func TestNewServer(t *testing.T) {
	config := pprof.NewDefConfig()
	config.ListenAddr = "localhost:6061"
	server := pprof.NewServer(log.NewNopLogger(), config)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Millisecond)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	assert.NoError(t, server.Init(ctx, nil))
	err := server.Start(group, ctx)
	assert.NoError(t, err)

	assert.Panics(t, func() {
		http.HandleFunc("/debug/pprof/", func(w http.ResponseWriter, r *http.Request) {
		})
	})

	<-ctx.Done()
	assert.NoError(t, server.Close())
	assert.NoError(t, group.Wait())
}
