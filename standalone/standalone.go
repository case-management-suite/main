package standalone

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/case-management-suite/api"
	"github.com/case-management-suite/caseservice"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/rulesengineservice"
)

type StandaloneApp struct {
	APIServer         *api.APIServer
	CaseMicroservice  *caseservice.CaseMicroervice
	RulesMicroservice *rulesengineservice.RulesMicoservice
}

var START_TIMEOUT time.Duration = 10 * time.Second
var STOP_TIMEOUT time.Duration = 20 * time.Second

func (s StandaloneApp) Start(ctx context.Context) error {
	errchans := server.RunAllAsync(
		ctx,
		START_TIMEOUT,
		s.APIServer.Start,
		s.CaseMicroservice.Start,
		s.RulesMicroservice.Start,
	)

	for _, v := range errchans {
		if err := <-v; err != nil {
			return err
		}
	}
	return nil
}

func (s StandaloneApp) Stop(ctx context.Context) error {
	errchans := server.RunAllAsync(
		ctx,
		START_TIMEOUT,
		s.APIServer.Stop,
		s.CaseMicroservice.Stop,
		s.RulesMicroservice.Stop,
	)

	for _, v := range errchans {
		if err := <-v; err != nil {
			return err
		}
	}
	return nil
}

func (s StandaloneApp) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Start(ctx); err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer func() {
		if err := s.Stop(ctx); err != nil {
			log.Fatal(err)
		}
	}()
}

func NewStandaloneApp(appConfig config.AppConfig) (*StandaloneApp, error) {
	apiServer, caseService, rulesService, err := App(appConfig).BuildApp(appConfig)
	if err != nil {
		return nil, err
	}
	return &StandaloneApp{
		APIServer:         apiServer,
		CaseMicroservice:  caseService,
		RulesMicroservice: rulesService,
	}, nil
}
