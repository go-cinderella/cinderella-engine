package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	. "github.com/go-cinderella/cinderella-engine/engine/variable"
	"time"
)

var _ delegate.DelegateExecution = (*ExecutionEntity)(nil)

type ExecutionEntity struct {
	AbstractEntity
	*VariableScopeImpl
	BusinessKey        string
	CurrentFlowElement delegate.FlowElement
	DeploymentId       string
	ProcessInstanceId  string
	CurrentActivityId  string

	BusinessStatus               string    `json:"businessStatus"`
	Suspended                    bool      `json:"suspended"`
	Ended                        bool      `json:"ended"`
	ProcessDefinitionId          string    `json:"processDefinitionId"`
	ProcessDefinitionName        string    `json:"processDefinitionName"`
	ProcessDefinitionDescription string    `json:"processDefinitionDescription"`
	ActivityId                   string    `json:"activityId"`
	StartUserId                  string    `json:"startUserId"`
	StartTime                    time.Time `json:"StartTime"`
	CallbackId                   string    `json:"callbackId"`
	CallbackType                 string    `json:"callbackType"`
	ReferenceId                  string    `json:"referenceId"`
	ReferenceType                string    `json:"referenceType"`
	TenantId                     *string   `json:"tenantId"`
}

func (execution *ExecutionEntity) GetTenantId() *string {
	return execution.TenantId
}

func (execution *ExecutionEntity) SetCurrentActivityId(currentActivityId string) {
	execution.CurrentActivityId = currentActivityId
}

func (execution *ExecutionEntity) GetExecutionId() string {
	return execution.Id
}

func (execution *ExecutionEntity) SetBusinessKey(businessKey string) {
	execution.BusinessKey = businessKey
}

func (execution ExecutionEntity) GetCurrentFlowElement() delegate.FlowElement {
	return execution.CurrentFlowElement
}

func (execution *ExecutionEntity) SetCurrentFlowElement(flow delegate.FlowElement) {
	execution.CurrentFlowElement = flow
	execution.CurrentActivityId = flow.GetId()
}

func (execution ExecutionEntity) GetDeploymentId() string {
	return execution.DeploymentId
}

func (execution *ExecutionEntity) SetDeploymentId(deploymentId string) {
	execution.DeploymentId = deploymentId
}

func (execution ExecutionEntity) GetProcessInstanceId() string {
	return execution.ProcessInstanceId
}

func (execution *ExecutionEntity) SetProcessInstanceId(processInstanceId string) {
	execution.ProcessInstanceId = processInstanceId
}

func (execution ExecutionEntity) GetProcessDefinitionId() string {
	return execution.ProcessDefinitionId
}

func (execution *ExecutionEntity) SetProcessDefinitionId(processDefineId string) {
	execution.ProcessDefinitionId = processDefineId
}

func (execution ExecutionEntity) GetCurrentActivityId() string {
	return execution.CurrentActivityId
}

func (execution ExecutionEntity) GetProcessVariable() map[string]interface{} {
	variables := execution.GetVariable()
	variableLocal := execution.GetVariableLocal()
	if variableLocal != nil {
		for k, v := range variableLocal {
			variables[k] = v
		}
	}
	return variables
}

func (execution ExecutionEntity) GetVariable() map[string]interface{} {
	variableManager := datamanager.GetVariableDataManager()
	variables, err := variableManager.SelectByProcessInstanceId(execution.GetProcessInstanceId())
	if err == nil {
		return execution.HandleVariable(variables)
	}
	return nil
}

func (execution ExecutionEntity) HandleVariable(variables []Variable) map[string]interface{} {
	variableManager := GetVariableManager()
	variableTypes := variableManager.VariableTypes
	var variableMap = make(map[string]interface{}, 0)
	for _, variable := range variables {
		variableType := variableTypes.GetVariableType(variable.Type_)
		value := variableType.GetValue(&variable)
		variableMap[variable.Name_] = value
	}
	return variableMap
}

func (execution ExecutionEntity) GetSourceActivityExecution() ExecutionEntity {
	return execution
}

func (execution ExecutionEntity) GetSpecificVariable(variableName string) (Variable, error) {
	variableDataManager := datamanager.GetVariableDataManager()
	return variableDataManager.SelectByProcessInstanceIdAndName(variableName, execution.ProcessInstanceId)
}
