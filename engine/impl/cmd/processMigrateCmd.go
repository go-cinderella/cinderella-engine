package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
	"github.com/samber/lo"
)

var _ engine.Command = (*ProcessMigrateCmd)(nil)

type DiffActivity struct {
	ActivityId    string
	ActivityName  string
	NewActivityId string
}

type ProcessMigrateCmd struct {
	Transactional bool
	Ctx           context.Context

	OldDeploymentId string
	NewDeploymentId string

	DiffActivities []DiffActivity
}

func (g ProcessMigrateCmd) Execute(commandContext engine.Context) (interface{}, error) {
	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	oldProcessDefinitionEntity, err := processDefinitionEntityManager.FindByDeploymentId(g.OldDeploymentId)
	if err != nil {
		return nil, err
	}

	newProcessDefinitionEntity, err := processDefinitionEntityManager.FindByDeploymentId(g.NewDeploymentId)
	if err != nil {
		return nil, err
	}

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	err = executionEntityManager.MigrateProcessInstanceProcDefIdAndStartActId(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	err = executionEntityManager.MigrateExecutionProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	taskEntityManager := entitymanager.GetTaskEntityManager()
	err = taskEntityManager.MigrateProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	err = historicActivityInstanceEntityManager.MigrateProcDefID(oldProcessDefinitionEntity, newProcessDefinitionEntity)
	if err != nil {
		return nil, err
	}

	bpmnXMLConverter := converter.BpmnXMLConverter{}
	bpmnModel := bpmnXMLConverter.ConvertToBpmnModel(newProcessDefinitionEntity.ResourceContent)
	process := bpmnModel.GetProcess()

	lo.ForEachWhile(g.DiffActivities, func(item DiffActivity, index int) (goon bool) {
		err = executionEntityManager.MigrateProcessInstanceBusinessStatus(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId)
		if err != nil {
			return false
		}

		err = executionEntityManager.MigrateExecutionActID(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId)
		if err != nil {
			return false
		}

		newFlowElement := process.GetFlowElement(item.NewActivityId)
		err = taskEntityManager.MigrateNameAndTaskDefKey(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId, newFlowElement.GetName())
		if err != nil {
			return false
		}

		err = historicActivityInstanceEntityManager.MigrateAct(newProcessDefinitionEntity, item.ActivityId, item.NewActivityId, newFlowElement.GetName(), newFlowElement.GetHandlerType())
		if err != nil {
			return false
		}

		return true
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (g ProcessMigrateCmd) Context() context.Context {
	return g.Ctx
}

func (g ProcessMigrateCmd) IsTransactional() bool {
	return g.Transactional
}
