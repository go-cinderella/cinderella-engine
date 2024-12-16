package interceptor

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/agenda"
	"github.com/go-cinderella/cinderella-engine/engine/query"
)

var _ engine.Context = (*CommandContext)(nil)

type CommandContext struct {
	agenda                     agenda.CinderellaEngineAgenda
	processEngineConfiguration engine.ProcessEngineConfiguration
	query                      *query.Query
}

func (c *CommandContext) SetAgenda(agenda agenda.CinderellaEngineAgenda) {
	c.agenda = agenda
}

func (c *CommandContext) SetProcessEngineConfiguration(processEngineConfiguration engine.ProcessEngineConfiguration) {
	c.processEngineConfiguration = processEngineConfiguration
}

func (c CommandContext) GetQuery() *query.Query {
	return c.query
}

func (c CommandContext) WithQuery(q *query.Query) engine.Context {
	c.query = q
	return &c
}

func (c CommandContext) GetAgenda() agenda.CinderellaEngineAgenda {
	return c.agenda
}

func (c CommandContext) GetProcessEngineConfiguration() engine.ProcessEngineConfiguration {
	return c.processEngineConfiguration
}
