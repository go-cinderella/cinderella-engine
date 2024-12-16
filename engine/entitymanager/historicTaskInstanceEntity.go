package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	"time"
)

type HistoricTaskInstanceEntity struct {
	AbstractEntity
	ProcessDefinitionId  string              `json:"processDefinitionId"`
	ProcessDefinitionUrl string              `json:"processDefinitionUrl"`
	ProcessInstanceId    string              `json:"processInstanceId"`
	ProcessInstanceUrl   string              `json:"processInstanceUrl"`
	ExecutionId          string              `json:"executionId"`
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	DeleteReason         string              `json:"deleteReason"`
	Owner                string              `json:"owner"`
	Assignee             string              `json:"Assignee"`
	StartTime            time.Time           `json:"StartTime"`
	EndTime              *time.Time          `json:"endTime"`
	DurationInMillis     int                 `json:"durationInMillis"`
	WorkTimeInMillis     int                 `json:"workTimeInMillis"`
	ClaimTime            *time.Time          `json:"claimTime"`
	TaskDefinitionKey    string              `json:"taskDefinitionKey"`
	FormKey              string              `json:"formKey"`
	Priority             int                 `json:"priority"`
	DueDate              *time.Time          `json:"dueDate"`
	ParentTaskId         string              `json:"parentTaskId"`
	Url                  string              `json:"url"`
	Variables            []variable.Variable `json:"variables"`
	TenantId             string              `json:"tenantId"`
	Category             string              `json:"category"`
}

func (historicTaskInstanceEntity HistoricTaskInstanceEntity) SetAssignee(assignee string) {
	historicTaskInstanceEntity.Assignee = assignee
}
