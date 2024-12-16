package interceptor

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/impl/agenda"
)

type CommandContextFactory struct {
	ProcessEngineConfiguration engine.ProcessEngineConfiguration
}

func (factory CommandContextFactory) CreateCommandContext() engine.Context {
	agendaInstance := &agenda.DefaultCinderellaEngineAgenda{}
	context := &CommandContext{
		agenda:                     agendaInstance,
		processEngineConfiguration: factory.ProcessEngineConfiguration,
	}
	agendaInstance.SetContext(context)
	return context
}
