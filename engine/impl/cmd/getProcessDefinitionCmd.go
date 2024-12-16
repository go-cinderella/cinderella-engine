package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

var _ engine.Command = (*GetProcessDefinitionCmd)(nil)

type GetProcessDefinitionCmd struct {
	ProcessDefinitionId string
	Ctx                 context.Context
	Transactional       bool
}

func (getProcessDefinitionCmd GetProcessDefinitionCmd) IsTransactional() bool {
	return getProcessDefinitionCmd.Transactional
}

func (getProcessDefinitionCmd GetProcessDefinitionCmd) Context() context.Context {
	return getProcessDefinitionCmd.Ctx
}

func (getProcessDefinitionCmd GetProcessDefinitionCmd) Execute(ctx engine.Context) (interface{}, error) {
	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	processDefinitionEntity, err := processDefinitionEntityManager.FindProcessDefinitionById(getProcessDefinitionCmd.ProcessDefinitionId)
	return processDefinitionEntity, err
}
