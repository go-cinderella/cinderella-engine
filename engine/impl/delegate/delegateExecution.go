package delegate

import (
	"github.com/go-cinderella/cinderella-engine/engine/variable"
)

type DelegateExecution interface {
	GetId() string
	SetId(id string)

	SetBusinessKey(businessKey string)

	GetCurrentFlowElement() FlowElement

	SetCurrentFlowElement(flow FlowElement)

	GetDeploymentId() string

	SetDeploymentId(deploymentId string)

	GetProcessInstanceId() string

	SetProcessInstanceId(processInstanceId string)

	GetProcessDefinitionId() string

	SetProcessDefinitionId(processDefineId string)

	GetCurrentActivityId() string

	SetCurrentActivityId(currentActivityId string)

	//SetVariable(execution ExecutionEntity,variables map[string]interface{}) error

	GetSpecificVariable(variableName string) (variable.Variable, error)

	GetVariable() map[string]interface{}

	GetProcessVariable() map[string]interface{}

	GetExecutionId() string

	GetTenantId() *string
}
