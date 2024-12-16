package entitymanager

import "time"

type HistoricActivityInstanceEntity struct {
	AbstractEntity
	ActivityId              string     `json:"activityId"`
	ActivityName            string     `json:"activityName"`
	ActivityType            string     `json:"activityType"`
	ProcessDefinitionId     string     `json:"processDefinitionId"`
	ProcessInstanceId       string     `json:"processInstanceId"`
	ExecutionId             string     `json:"executionId"`
	TaskId                  string     `json:"taskId"`
	CalledProcessInstanceId string     `json:"calledProcessInstanceId"`
	Assignee                string     `json:"assignee"`
	StartTime               time.Time  `json:"startTime"`
	EndTime                 *time.Time `json:"endTime"`
	DurationInMillis        int        `json:"durationInMillis"`
	TenantId                *string    `json:"tenantId"`
	DeleteReason            *string    `json:"deleteReason"`
}
