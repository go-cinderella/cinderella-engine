package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
)

var _ engine.Command = (*GetProcessStatusCmd)(nil)

type GetProcessStatusCmd struct {
	Ctx               context.Context
	ProcessInstanceId string
	Transactional     bool
}

func (g GetProcessStatusCmd) IsTransactional() bool {
	return g.Transactional
}

func (g GetProcessStatusCmd) Execute(commandContext engine.Context) (interface{}, error) {
	historicProcessInstanceEntityManager := entitymanager.GetHistoricProcessInstanceEntityManager()
	processInstance, err := historicProcessInstanceEntityManager.FindById(g.ProcessInstanceId)
	if err != nil {
		return nil, err
	}

	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	resourceEntity, err := processDefinitionEntityManager.FindResourceEntityByProcessDefinitionById(processInstance.GetProcessDefinitionId())
	if err != nil {
		return nil, err
	}

	bpmnXMLConverter := converter.BpmnXMLConverter{}
	bpmnModel := bpmnXMLConverter.ConvertToBpmnModel(resourceEntity.GetBytes())
	process := bpmnModel.GetProcess()
	flowElement := process.GetFlowElement(processInstance.BusinessStatus)
	return flowElement, nil
}

func (g GetProcessStatusCmd) Context() context.Context {
	return g.Ctx
}
