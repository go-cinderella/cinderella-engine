package procdef

import (
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
)

type ListRequest struct {
	request.ListCommonRequest // order allow processInstanceId,processDefinitionId,businessKey,startTime,endTime,duration,tenantId
	request.WithTenant
	ProcessDefinitionKey      string   `json:"processDefinitionKey,omitempty"`
	ProcessDefinitionKeyIn    []string `json:"processDefinitionKeyIn,omitempty"`
	ProcessDefinitionKeyNotIn []string `json:"processDefinitionKeyNotIn,omitempty"`
	ProcessDefinitionName     string   `json:"processDefinitionName,omitempty"`
	ProcessDefinitionVersion  int      `json:"processDefinitionVersion,omitempty"`
	ProcessDefinitionCategory string   `json:"processDefinitionCategory,omitempty"`
	ProcessDefinitionId       string   `json:"processDefinitionId,omitempty"`
	DeploymentId              string   `json:"deploymentId,omitempty"`
	DeploymentIdIn            []string `json:"deploymentIdIn,omitempty"`
}
