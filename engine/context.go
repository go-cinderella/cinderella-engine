package engine

import (
	"github.com/go-cinderella/cinderella-engine/engine/agenda"
)

type Context interface {
	SetAgenda(agenda agenda.CinderellaEngineAgenda)
	SetProcessEngineConfiguration(processEngineConfiguration ProcessEngineConfiguration)
	GetAgenda() agenda.CinderellaEngineAgenda
	GetProcessEngineConfiguration() ProcessEngineConfiguration
}
