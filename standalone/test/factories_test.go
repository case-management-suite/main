package test

import (
	"os"
	"testing"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/main/standalone"
	"github.com/case-management-suite/testutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestValidateFactories(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false}).Level(zerolog.DebugLevel)
	appConfig := config.NewLocalAppConfigWithParams(config.GraphQL)
	apiServer, caseServiceServer, rulesServer, err := standalone.App(appConfig).BuildApp(appConfig)
	testutil.AssertNilError(err, t)
	testutil.AssertNonNil(apiServer, t)
	testutil.AssertNonNil(caseServiceServer, t)
	testutil.AssertNonNil(rulesServer, t)
}
