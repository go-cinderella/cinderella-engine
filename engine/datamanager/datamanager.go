package datamanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/datamanager/abstract"
	"github.com/go-cinderella/cinderella-engine/engine/model"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
)

var (
	deploymentDataManager = &DeploymentDataManager{DataManager: abstract.DataManager{&model.ActReDeployment{}}}

	processDefinitionDataManager = &ProcessDefinitionDataManager{DataManager: abstract.DataManager{&model.ActReProcdef{}}}

	executionDataManager = &ExecutionDataManager{DataManager: abstract.DataManager{&model.ActRuExecution{}}}

	historicActinstDataManager = &HistoricActinstDataManager{DataManager: abstract.DataManager{&model.ActHiActinst{}}}

	historicIdentityLinkDataManager = &HistoricIdentityLinkDataManager{DataManager: abstract.DataManager{&model.ActHiIdentitylink{}}}

	historicProcessDataManager = &HistoricProcessDataManager{DataManager: abstract.DataManager{&model.ActHiProcinst{}}}

	historicTaskDataManager = &HistoricTaskDataManager{DataManager: abstract.DataManager{&model.ActHiTaskinst{}}}

	historicVariableDataManager = &HistoricVariableDataManager{DataManager: abstract.DataManager{&variable.Variable{}}}

	identityLinkDataManager = &IdentityLinkDataManager{DataManager: abstract.DataManager{&model.ActRuIdentitylink{}}}

	resourceDataManager = &ResourceDataManager{DataManager: abstract.DataManager{AbstractModel: &model.ActGeBytearray{}}}

	taskDataManager = &TaskDataManager{DataManager: abstract.DataManager{&model.ActRuTask{}}}

	variableDataManager = &VariableDataManager{DataManager: abstract.DataManager{AbstractModel: &variable.Variable{}}}
)

func GetDeploymentDataManager() *DeploymentDataManager {
	return deploymentDataManager
}

func GetProcessDefinitionDataManager() *ProcessDefinitionDataManager {
	return processDefinitionDataManager
}

func GetExecutionDataManager() *ExecutionDataManager {
	return executionDataManager
}

func GetHistoricActinstDataManager() *HistoricActinstDataManager {
	return historicActinstDataManager
}

func GetHistoricIdentityLinkDataManager() *HistoricIdentityLinkDataManager {
	return historicIdentityLinkDataManager
}

func GetHistoricProcessDataManager() *HistoricProcessDataManager {
	return historicProcessDataManager
}

func GetHistoricTaskDataManager() *HistoricTaskDataManager {
	return historicTaskDataManager
}

func GetHistoricVariableDataManager() *HistoricVariableDataManager {
	return historicVariableDataManager
}

func GetIdentityLinkDataManager() *IdentityLinkDataManager {
	return identityLinkDataManager
}

func GetResourceDataManager() *ResourceDataManager {
	return resourceDataManager
}

func GetTaskDataManager() *TaskDataManager {
	return taskDataManager
}

func GetVariableDataManager() *VariableDataManager {
	return variableDataManager
}
