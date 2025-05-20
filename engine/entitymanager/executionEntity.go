package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
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
	StartUserId                  string    `json:"startUserId"`
	StartTime                    time.Time `json:"StartTime"`
	CallbackId                   string    `json:"callbackId"`
	CallbackType                 string    `json:"callbackType"`
	ReferenceId                  string    `json:"referenceId"`
	ReferenceType                string    `json:"referenceType"`
	TenantId                     *string   `json:"tenantId"`
	isMultiInstanceRoot          bool
	ParentId                     string
	Parent                       delegate.DelegateExecution
}

func (execution *ExecutionEntity) RemoveVariablesLocal(variableNames []string) error {
	variableDataManager := datamanager.GetVariableDataManager()
	return variableDataManager.DeleteByExecutionIdAndNames(execution.GetExecutionId(), variableNames)
}

func (execution *ExecutionEntity) GetVariableLocal(variableName string) (value interface{}, ok bool, err error) {
	variableDataManager := datamanager.GetVariableDataManager()

	vari, ok, err := variableDataManager.SelectByExecutionIdAndName(variableName, execution.GetExecutionId())
	if err != nil {
		return nil, false, err
	}

	if ok {
		return execution.getValue(vari), ok, nil
	}

	return nil, false, nil
}

func (execution *ExecutionEntity) GetParent() (delegate.DelegateExecution, error) {
	if execution.Parent != nil {
		return execution.Parent, nil
	}

	if stringutils.IsNotEmpty(execution.ParentId) {
		parentExection, err := executionEntityManager.FindById(execution.ParentId)
		if err != nil {
			return nil, err
		}
		return &parentExection, nil
	}

	return nil, nil
}

func (execution *ExecutionEntity) SetParent(parent delegate.DelegateExecution) {
	execution.Parent = parent
}

func (execution *ExecutionEntity) GetParentId() string {
	return execution.ParentId
}

func (execution *ExecutionEntity) SetParentId(parentId string) {
	execution.ParentId = parentId
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

func (execution ExecutionEntity) GetVariables() (map[string]interface{}, error) {
	var variables []variable.Variable

	variableManager := datamanager.GetVariableDataManager()
	localVariables, err := variableManager.SelectByExecutionId(execution.GetExecutionId())
	if err != nil {
		return nil, err
	}

	variables = append(variables, localVariables...)

	parent, err := execution.GetParent()
	if err != nil {
		return nil, err
	}

	for parent != nil {
		localVariables, err = variableManager.SelectByExecutionId(parent.GetExecutionId())
		if err != nil {
			return nil, err
		}

		variables = append(variables, localVariables...)

		parent, err = parent.GetParent()
		if err != nil {
			return nil, err
		}
	}

	return execution.handleVariables(variables), nil
}

func (execution ExecutionEntity) handleVariables(variables []variable.Variable) map[string]interface{} {
	var variableMap = make(map[string]interface{}, 0)
	for _, item := range variables {
		if lo.HasKey(variableMap, item.Name_) {
			continue
		}
		variableMap[item.Name_] = execution.getValue(item)
	}
	return variableMap
}

func (execution ExecutionEntity) getValue(vari variable.Variable) interface{} {
	variableManager := variable.GetVariableManager()
	variableTypes := variableManager.VariableTypes
	variableType := variableTypes.GetVariableType(vari.Type_)
	return variableType.GetValue(&vari)
}

func IsIntegral(val float64) bool {
	return val == float64(int(val))
}

// SetProcessVariables 保存流程变量
func (execution ExecutionEntity) SetProcessVariables(variables map[string]interface{}) error {
	return execution.doSetVariablesLocal(variables, execution.GetProcessInstanceId(), "")
}

// SetVariablesLocal 保存执行实例变量
func (execution ExecutionEntity) SetVariablesLocal(variables map[string]interface{}) error {
	return execution.doSetVariablesLocal(variables, execution.GetExecutionId(), "")
}

func (execution ExecutionEntity) SetVariableLocal(variableName string, value interface{}) error {
	variables := make(map[string]interface{})
	variables[variableName] = value
	return execution.SetVariablesLocal(variables)
}

func (execution ExecutionEntity) doSetVariablesLocal(variables map[string]interface{}, executionId string, taskId string) error {
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
		if stringutils.IsNotEmpty(taskId) {
			vari.TaskID_ = &taskId
		}

		if err := variableEntityManager.UpsertVariable(vari); err != nil {
			return err
		}
	}

	return nil
}

func (execution *ExecutionEntity) SetActive(active bool) error {
	executionDataManager := datamanager.GetExecutionDataManager()
	return executionDataManager.UpdateActive(execution.GetExecutionId(), active)
}

func (execution *ExecutionEntity) IsActive() bool {
	executionDataManager := datamanager.GetExecutionDataManager()
	isActive, err := executionDataManager.IsActive(execution.GetExecutionId())
	if err != nil {
		return false
	}
	return isActive
}

func CreateChildExecution(parentExecution delegate.DelegateExecution) ExecutionEntity {
	newExecution := ExecutionEntity{
		ProcessInstanceId:   parentExecution.GetProcessInstanceId(),
		ProcessDefinitionId: parentExecution.GetProcessDefinitionId(),
		ParentId:            parentExecution.GetExecutionId(),
		Parent:              parentExecution,
		StartTime:           time.Now().UTC(),
	}
	return newExecution
}

func CreateExecution(execution delegate.DelegateExecution) ExecutionEntity {
	newExecution := ExecutionEntity{
		ProcessInstanceId:   execution.GetProcessInstanceId(),
		ProcessDefinitionId: execution.GetProcessDefinitionId(),
		ParentId:            execution.GetParentId(),
		StartTime:           time.Now().UTC(),
	}
	return newExecution
}
