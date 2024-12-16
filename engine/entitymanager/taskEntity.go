package entitymanager

import (
	bpmn_model "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	"time"
)

var _ delegate.DelegateExecution = (*TaskEntity)(nil)

type TaskEntity struct {
	AbstractEntity
	ExecutionEntity
	VariableScopeImpl
	Assignee       *string
	StartTime      time.Time
	TaskDefineName string
	ClaimTime      *time.Time
	Variables      map[string]interface{}

	Url                       string     `json:"url"`
	Owner                     string     `json:"owner"`
	DelegationState           string     `json:"delegationState"`
	Description               string     `json:"description"`
	DueDate                   *time.Time `json:"dueDate"`
	Priority                  int        `json:"priority"`
	Suspended                 bool       `json:"suspended"`
	TaskDefinitionKey         string     `json:"taskDefinitionKey"`
	ScopeDefinitionId         string     `json:"scopeDefinitionId"`
	ScopeId                   string     `json:"scopeId"`
	SubScopeId                string     `json:"subScopeId"`
	ScopeType                 string     `json:"scopeType"`
	PropagatedStageInstanceId string     `json:"propagatedStageInstanceId"`
	TenantId                  *string    `json:"tenantId"`
	Category                  string     `json:"category"`
	FormKey                   *string    `json:"formKey"`
	ParentTaskId              string     `json:"parentTaskId"`
	ParentTaskUrl             string     `json:"parentTaskUrl"`
	ExecutionId               string     `json:"executionId"`
	ExecutionUrl              string     `json:"executionUrl"`
	ProcessInstanceId         string     `json:"processInstanceId"`
	ProcessInstanceUrl        string     `json:"processInstanceUrl"`
	ProcessDefinitionId       string     `json:"processDefinitionId"`
	ProcessDefinitionUrl      string     `json:"processDefinitionUrl"`
}

func (taskEntiy *TaskEntity) SetAssignee(assignee *string) {
	taskEntiy.Assignee = assignee
}

func (taskEntiy *TaskEntity) GetAssignee() *string {
	return taskEntiy.Assignee
}

func (taskEntiy *TaskEntity) SetStartTime(startTime time.Time) {
	taskEntiy.StartTime = startTime
}

func (taskEntiy *TaskEntity) GetClaimTime() *time.Time {
	return taskEntiy.ClaimTime
}

func (taskEntiy *TaskEntity) SetTaskDefineKey(taskDefineKey string) {
	taskEntiy.TaskDefinitionKey = taskDefineKey
}

func (taskEntiy *TaskEntity) GetTaskDefineKey() string {
	return taskEntiy.TaskDefinitionKey
}

func (taskEntiy *TaskEntity) SetTaskDefineName(taskDefineName string) {
	taskEntiy.TaskDefineName = taskDefineName
}

func (taskEntiy TaskEntity) GetVariable() map[string]interface{} {
	//variableManager := task.GetVariableEntityManager()
	//variables, err := variableManager.SelectByTaskId(taskEntiy.TaskId)
	//if err != nil {
	//	return task.HandleVariable(variables)
	//}
	return nil
}

func (taskEntiy TaskEntity) GetSpecificVariable(variableName string) (variable.Variable, error) {
	variableDataManager := datamanager.GetVariableDataManager()
	return variableDataManager.SelectTaskId(variableName, taskEntiy.GetId())
}

func (taskEntiy *TaskEntity) SetClaimTime(claimTime *time.Time) {
	taskEntiy.ClaimTime = claimTime
}

func (taskEntiy *TaskEntity) GetExecutionId() string {
	return taskEntiy.ExecutionId
}

func (taskEntiy *TaskEntity) SetExecutionId(executionId string) {
	taskEntiy.ExecutionId = executionId
}

func NewTaskEntity(execution delegate.DelegateExecution, userTask bpmn_model.UserTask) TaskEntity {
	return TaskEntity{
		ExecutionId:         execution.GetExecutionId(),
		ProcessInstanceId:   execution.GetProcessInstanceId(),
		ProcessDefinitionId: execution.GetProcessDefinitionId(),
		TaskDefineName:      userTask.Name,
		TaskDefinitionKey:   userTask.Id,
		StartTime:           time.Now().UTC(),
		TenantId:            execution.GetTenantId(),
		FormKey:             userTask.FormKey,
	}
}
