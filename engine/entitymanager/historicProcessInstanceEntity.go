package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	"time"
)

type HistoricProcessInstanceEntity struct {
	AbstractEntity
	BusinessKey            string              `json:"businessKey"`
	BusinessStatus         string              `json:"businessStatus"`
	ProcessDefinitionId    string              `json:"processDefinitionId"`
	ProcessDefinitionUrl   string              `json:"processDefinitionUrl"`
	StartTime              *time.Time          `json:"StartTime"`
	EndTime                *time.Time          `json:"endTime"`
	DurationInMillis       int                 `json:"durationInMillis"`
	StartUserId            string              `json:"startUserId"`
	StartActivityId        string              `json:"startActivityId"`
	EndActivityId          string              `json:"endActivityId"`
	DeleteReason           string              `json:"deleteReason"`
	SuperProcessInstanceId string              `json:"superProcessInstanceId"`
	Url                    string              `json:"url"`
	Variables              []variable.Variable `json:"variable"`
	TenantId               string              `json:"tenantId"`
	ProcessDefinitionName  string              `json:"processDefinitionName"`
}

func (historicProcessInstanceEntity HistoricProcessInstanceEntity) GetProcessDefinitionId() string {
	return historicProcessInstanceEntity.ProcessDefinitionId
}

func (historicProcessInstanceEntity *HistoricProcessInstanceEntity) SetProcessDefinitionId(processDefinitionId string) {
	historicProcessInstanceEntity.ProcessDefinitionId = processDefinitionId
}
