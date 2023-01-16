package main

import (
	"context"
	"os"
	"time"

	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/main/standalone"
	"github.com/rs/zerolog/log"
)

// func main() {
// 	appConfig := config.NewLocalAppConfig()
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// 	go runCasesService(appConfig)
// 	go runRulesAPI(appConfig)
// 	runGraphQLAPI(appConfig)

// 	// app := runApp(appConfig)
// 	// app.Run()
// }

// func runCasesService(appConfig config.AppConfig) {
// 	caseservice.NewCaseServiceGRPCServer(appConfig).Run()
// }

// func runGraphQLAPI(appConfig config.AppConfig) {
// 	graphql.CreateLiteGraphQLAPIServer(appConfig).Run()
// }

// func runRulesAPI(appConfig config.AppConfig) {
// 	rulesengineservice.NewRulesServiceCServer(appConfig).Run()
// }

func main() {
	appConfig := config.NewLocalAppConfigWithParams(config.GraphQL)
	// app, stop := factory.NewFxApp(appConfig)
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// defer stop(ctx)

	// app.Run()
	api, caseServiceServer, rulesServer, err := standalone.App(appConfig).BuildApp(appConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build the app")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	defer func(ctx context.Context) {
		defer api.Stop(ctx)
		defer caseServiceServer.Stop(ctx)
		defer rulesServer.Stop(ctx)
	}(ctx)

	err = rulesServer.Start(ctx)
	if err != nil {
		panic("rule server")
	}
	err = caseServiceServer.Start(ctx)
	if err != nil {
		panic("case server")
	}
	err = api.Start(ctx)
	if err != nil {
		panic("api")
	}

	interrupt := make(chan os.Signal, 1)
	<-interrupt
}
