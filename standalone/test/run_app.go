package test

import (
	"context"

	"os"
	"testing"
	"time"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/main/standalone"
	"github.com/case-management-suite/testutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func RunApp(t *testing.T, fn func(t *testing.T)) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})
	appConfig := config.NewLocalAppConfigWithParams(config.GraphQL)

	app, err := standalone.NewStandaloneApp(appConfig)
	testutil.AssertNilError(err, t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = app.Start(ctx)
	testutil.AssertNilError(err, t)

	fn(t)

	defer func() {
		err := app.Stop(ctx)
		testutil.AssertNilError(err, t)
	}()
}
