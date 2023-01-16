package standalone

import (
	"context"

	"github.com/case-management-suite/common/config"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func NewFxApp(appConfig config.AppConfig) (*fx.App, func(context.Context)) {
	apiServer, caseServiceServer, rulesServer, err := App(appConfig).BuildApp(appConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build the application")
	}
	close := func(ctx context.Context) {
		defer apiServer.Stop(ctx)
		defer caseServiceServer.Stop(ctx)
		defer rulesServer.Stop(ctx)
	}
	return fx.New(apiServer.GetFxOption(), caseServiceServer.GetFxOption(), rulesServer.GetFxOption()), close
}
