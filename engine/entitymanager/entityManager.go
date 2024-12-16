package entitymanager

var (
	processDefinitionEntityManager        *ProcessDefinitionEntityManager
	executionEntityManager                *ExecutionEntityManager
	taskEntityManager                     *TaskEntityManager
	identityLinkManager                   *IdentityLinkEntityManager
	variableEntityManager                 *VariableEntityManager
	historicTaskInstanceEntityManager     *HistoricTaskInstanceEntityManager
	historicActivityInstanceEntityManager *HistoricActivityInstanceEntityManager
	deploymentEntityManager               *DeploymentEntityManager
	resourceEntityManager                 *ResourceEntityManager
	historicIdentityLinkEntityManager     *HistoricIdentityLinkEntityManager
	historicProcessInstanceEntityManager  *HistoricProcessInstanceEntityManager
)

func init() {
	processDefinitionEntityManager = &ProcessDefinitionEntityManager{}
	executionEntityManager = &ExecutionEntityManager{}
	taskEntityManager = &TaskEntityManager{}
	historicProcessInstanceEntityManager = &HistoricProcessInstanceEntityManager{}
	variableEntityManager = &VariableEntityManager{}
	identityLinkManager = &IdentityLinkEntityManager{}
	historicTaskInstanceEntityManager = &HistoricTaskInstanceEntityManager{}
	historicActivityInstanceEntityManager = &HistoricActivityInstanceEntityManager{}
	deploymentEntityManager = &DeploymentEntityManager{}
	resourceEntityManager = &ResourceEntityManager{}
	historicIdentityLinkEntityManager = &HistoricIdentityLinkEntityManager{}
}

type EntityManager interface {
	Insert(interface{}) error

	GetById(id string, data interface{}) interface{}

	Delete(entity Entity)
}

func GetTaskEntityManager() *TaskEntityManager {
	return taskEntityManager
}
func GetProcessDefinitionEntityManager() *ProcessDefinitionEntityManager {
	return processDefinitionEntityManager
}

func GetDeploymentEntityManager() *DeploymentEntityManager {
	return deploymentEntityManager
}

func GetResourceEntityManager() *ResourceEntityManager {
	return resourceEntityManager
}
func GetExecutionEntityManager() *ExecutionEntityManager {
	return executionEntityManager
}

func GetIdentityLinkManager() *IdentityLinkEntityManager {
	return identityLinkManager
}

func GetVariableEntityManager() *VariableEntityManager {
	return variableEntityManager
}

func GetHistoricTaskInstanceEntityManager() *HistoricTaskInstanceEntityManager {
	return historicTaskInstanceEntityManager
}

func GetHistoricActivityInstanceEntityManager() *HistoricActivityInstanceEntityManager {
	return historicActivityInstanceEntityManager
}

func GetHistoricIdentityLinkEntityManager() *HistoricIdentityLinkEntityManager {
	return historicIdentityLinkEntityManager
}

func GetHistoricProcessInstanceEntityManager() *HistoricProcessInstanceEntityManager {
	return historicProcessInstanceEntityManager
}
