package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"reflect"
	"time"
)

var _ delegate.DelegateExecution = (*ExecutionEntity)(nil)

type ExecutionEntity struct {
	AbstractEntity
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
	isMultiInstanceRoot          bool
	parentId                     string
	parent                       delegate.DelegateExecution
}

func (execution *ExecutionEntity) GetParent() (delegate.DelegateExecution, error) {
	if execution.parent != nil {
		return execution.parent, nil
	}

	if stringutils.IsNotEmpty(execution.parentId) {
		parentExection, err := executionEntityManager.FindById(execution.parentId)
		if err != nil {
			return nil, err
		}
		return &parentExection, nil
	}

	return nil, nil
}

func (execution *ExecutionEntity) SetParent(parent delegate.DelegateExecution) {
	execution.parent = parent
}

func (execution *ExecutionEntity) GetParentId() string {
	return execution.parentId
}

func (execution *ExecutionEntity) SetParentId(parentId string) {
	execution.parentId = parentId
}

func (execution *ExecutionEntity) GetVariablesLocal() (map[string]interface{}, error) {
	variableManager := datamanager.GetVariableDataManager()
	variables, err := variableManager.SelectByExecutionId(execution.GetExecutionId())
	if err != nil {
		return nil, err
	}
	return execution.handleVariables(variables), nil
}

func (execution *ExecutionEntity) IsMultiInstanceRoot() bool {
	return execution.isMultiInstanceRoot
}

func (execution *ExecutionEntity) SetMultiInstanceRoot(isMultiInstanceRoot bool) {
	execution.isMultiInstanceRoot = isMultiInstanceRoot
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

func (execution ExecutionEntity) GetProcessVariables() (map[string]interface{}, error) {
	variableManager := datamanager.GetVariableDataManager()
	variables, err := variableManager.SelectByExecutionId(execution.GetProcessInstanceId())
	if err != nil {
		return nil, err
	}
	return execution.handleVariables(variables), nil
}

func (execution ExecutionEntity) handleVariables(variables []variable.Variable) map[string]interface{} {
	variableManager := variable.GetVariableManager()
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

func (execution ExecutionEntity) GetSpecificVariable(variableName string) (variable.Variable, bool) {
	variableDataManager := datamanager.GetVariableDataManager()

	vari, ok, _ := variableDataManager.SelectByExecutionIdAndName(variableName, execution.GetExecutionId())
	if ok {
		return vari, true
	}

	vari, ok, _ = variableDataManager.SelectByExecutionIdAndName(variableName, execution.GetProcessInstanceId())
	return vari, ok
}

func IsIntegral(val float64) bool {
	return val == float64(int(val))
}

// SetProcessVariables 保存流程变量
func (execution ExecutionEntity) SetProcessVariables(variables map[string]interface{}) error {
	return execution.doSetVariablesLocal(variables, execution.GetProcessInstanceId())
}

// SetVariablesLocal 保存执行实例变量
func (execution ExecutionEntity) SetVariablesLocal(variables map[string]interface{}) error {
	return execution.doSetVariablesLocal(variables, execution.GetExecutionId())
}

func (execution ExecutionEntity) SetVariableLocal(variableName string, value interface{}) error {
	variables := make(map[string]interface{})
	variables[variableName] = value
	return execution.SetVariablesLocal(variables)
}

func (execution ExecutionEntity) doSetVariablesLocal(variables map[string]interface{}, executionId string) error {
	variableManager := variable.GetVariableManager()
	variableTypes := variableManager.VariableTypes
	idGenerator := contextutil.GetIDGenerator()

	for k, v := range variables {
		if v == nil {
			continue
		}

		kind := reflect.TypeOf(v).Kind()
		if kind == reflect.Float64 {
			if IsIntegral(v.(float64)) {
				v = cast.ToInt(v)
				kind = reflect.Int
			}
		}

		variableType := variableTypes.GetVariableType(kind.String())
		if variableType == nil {
			continue
		}

		vari := variable.Variable{}
		vari.ID_, _ = idGenerator.NextID()
		vari.Name_ = k
		vari.Type_ = variableType.GetTypeName()
		vari.SetValue(v, variableType)
		vari.ProcInstID_ = lo.ToPtr(execution.GetProcessInstanceId())
		vari.ExecutionID_ = &executionId

		if err := variableEntityManager.UpsertVariable(vari); err != nil {
			return err
		}
	}

	return nil
}
