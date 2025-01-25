package contextutil

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/agenda"
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/go-cinderella/cinderella-engine/engine/idgenerator"
	"github.com/go-cinderella/cinderella-engine/engine/query"
	"github.com/go-resty/resty/v2"
)

func GetCommandExecutor(context engine.Context) engine.Executor {
	return context.GetProcessEngineConfiguration().GetCommandExecutor()
}

func GetAgenda() agenda.CinderellaEngineAgenda {
	return MustGetCommandContext().GetAgenda()
}

func GetAgendaFromContext(context engine.Context) agenda.CinderellaEngineAgenda {
	return context.GetAgenda()
}

func GetIDGenerator() idgenerator.IDGenerator {
	return MustGetCommandContext().GetProcessEngineConfiguration().GetIDGenerator()
}

func GetDeploymentSettings() map[string]interface{} {
	return MustGetCommandContext().GetProcessEngineConfiguration().GetDeploymentSettings()
}

func GetQuery() *query.Query {
	return query.Use(db.DB())
}

func GetBpmnDeployer() engine.Deployer {
	return MustGetCommandContext().GetProcessEngineConfiguration().GetBpmnDeployer()
}

func GetHttpClient() *resty.Client {
	return MustGetCommandContext().GetProcessEngineConfiguration().GetHttpClient()
}
