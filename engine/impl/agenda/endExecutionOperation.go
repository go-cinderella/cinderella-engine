package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
)

type EndExecutionOperation struct {
	AbstractOperation
}

var deleteReason = "process instance end"

func (end *EndExecutionOperation) Run() (err error) {
	err = deleteDataForExecution(end.Execution)
	if err != nil {
		return err
	}

	executionDataManager := datamanager.GetExecutionDataManager()
	err = executionDataManager.DeleteProcessInstance(end.Execution, nil)
	return err
}

func deleteDataForExecution(execution delegate.DelegateExecution) (err error) {
	processInstanceId := execution.GetProcessInstanceId()

	taskEntityManager := entitymanager.GetTaskEntityManager()
	tasks, err := taskEntityManager.FindByProcessInstanceId(processInstanceId)
	if err != nil {
		return err
	}
	for _, task := range tasks {
		taskEntityManager.DeleteTask(task, &deleteReason)
	}

	linkEntityManager := entitymanager.GetIdentityLinkManager()
	if err = linkEntityManager.DeleteByProcessInstanceId(processInstanceId); err != nil {
		return err
	}

	variableEntityManager := entitymanager.GetVariableEntityManager()
	if err = variableEntityManager.DeleteByProcessInstanceId(processInstanceId); err == nil {
		return err
	}

	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	err = historicActivityInstanceEntityManager.RecordActEndByProcessInstanceId(processInstanceId, &deleteReason)
	return err
}
