package standalone

import (
	"fmt"
	"reflect"

	"github.com/case-management-suite/api"
	"github.com/case-management-suite/api/controllers"
	"github.com/case-management-suite/casedb"
	"github.com/case-management-suite/caseservice"
	"github.com/case-management-suite/common/config"
	"github.com/case-management-suite/common/factory"
	"github.com/case-management-suite/common/metrics"
	"github.com/case-management-suite/common/server"
	"github.com/case-management-suite/common/service"
	"github.com/case-management-suite/queue"
	"github.com/case-management-suite/rulesengineservice"
	"github.com/case-management-suite/scheduler"
)

type AppFactories struct {
	factory.FactorySet
	api.APIFactories
	CaseServiceFactories  caseservice.CaseServiceFactories
	RulesServiceFactories rulesengineservice.RulesServiceFactories
}

func (f AppFactories) BuildApp(appConfig config.AppConfig) (*api.APIServer, *caseservice.CaseMicroervice, *rulesengineservice.RulesMicoservice, error) {
	if err := factory.ValidateFactorySet(f); err != nil {
		return nil, nil, nil, fmt.Errorf("factory: %s -> %w;", reflect.TypeOf(f).Name(), err)
	}
	// serverUtils := server.NewTestServerUtils()

	// workScheduler := f.WorkSchedulerFactories.BuildWorkScheduler(appConfig)

	apiServer, err := f.BuildAPI(appConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	caseService, err := f.CaseServiceFactories.BuildCaseService(appConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	rulesService, err := f.RulesServiceFactories.BuildRulesService(appConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	return apiServer, caseService, rulesService, nil
}

func App(appConfig config.AppConfig) AppFactories {
	WorkSchedulerFactories := scheduler.WorkSchedulerFactories{
		WorkSchedulerFactory: scheduler.NewWorkScheduler,
		QueueServiceFactory:  queue.QueueServiceFactory(appConfig.RulesServiceConfig.QueueType),
		// TODO: use prod
		ServiceUtilsFactory: service.NewTestServiceUtils,
	}

	return AppFactories{
		APIFactories: api.APIFactories{
			CaseServiceClientFactory:        caseservice.NewCaseServiceClient,
			ControllerFactory:               controllers.GetControllerFactory(),
			APIFactory:                      api.GetAPIFactory(appConfig.API.APIType),
			RulesEngineServiceClientFactory: rulesengineservice.NewRulesServiceClient,
			WorkSchedulerFactories:          WorkSchedulerFactories,
		},
		CaseServiceFactories: caseservice.CaseServiceFactories{
			CaseStorageServiceFactory: casedb.NewSQLCaseStorageService,
			CaseServiceFactory:        caseservice.NewCaseService,
			CaseServiceServerFactory:  caseservice.NewCaseServiceAPIServer,
			WorkSchedulerFactories:    WorkSchedulerFactories,
		},
		RulesServiceFactories: rulesengineservice.RulesServiceFactories{
			MetricsServiceFactory:           metrics.NewCaseMetricsService,
			RuleServiceFactory:              rulesengineservice.NewRulesService,
			RulesEngineServiceServerFactory: rulesengineservice.NewRulesEngineServiceServer,
			WorkSchedulerFactories:          WorkSchedulerFactories,
			// TODO: use prod
			ServerUtilsFactory: server.NewTestServerUtils,
		},
	}

}
