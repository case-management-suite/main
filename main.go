package main

import (
	"os"

	"github.com/case-management-suite/api/graphql"
	"github.com/case-management-suite/caseservice"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/rulesengineservice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	appConfig := config.NewLocalAppConfig()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	go runCasesService(appConfig)
	go runRulesAPI(appConfig)
	runGraphQLAPI(appConfig)

	// app := runApp(appConfig)
	// app.Run()
}

func runCasesService(appConfig config.AppConfig) {
	caseservice.NewCaseServiceGRPCServer(appConfig).Run()
}

func runGraphQLAPI(appConfig config.AppConfig) {
	graphql.CreateLiteGraphQLAPIServer(appConfig).Run()
}

func runRulesAPI(appConfig config.AppConfig) {
	rulesengineservice.NewRulesServiceCServer(appConfig).Run()
}
