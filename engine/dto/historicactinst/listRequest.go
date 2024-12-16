package historicactinst

import "github.com/go-cinderella/cinderella-engine/engine/dto/request"

type ListRequest struct {
	request.ListCommonRequest // order allow activityId,activityName,activityType,duration,endTime,executionId,activityInstanceId,processDefinitionId,processInstanceId,startTime,tenantId
	request.WithTenant
	ActivityId          string `json:"activityId,omitempty"`
	ActivityInstanceId  string `json:"activityInstanceId,omitempty"`
	ActivityName        string `json:"activityName,omitempty"`
	ActivityType        string `json:"activityType,omitempty"`
	ExecutionId         string `json:"executionId,omitempty"`
	Finished            *bool  `json:"finished,omitempty"`
	TaskAssignee        string `json:"taskAssignee,omitempty"`
	ProcessInstanceId   string `json:"processInstanceId,omitempty"`
	ProcessDefinitionId string `json:"processDefinitionId,omitempty"`
}
