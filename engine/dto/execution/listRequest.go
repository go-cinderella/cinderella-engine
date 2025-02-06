package execution

import "github.com/go-cinderella/cinderella-engine/engine/dto/request"

type ListRequest struct {
	request.ListCommonRequest          // order allow id,processDefinitionId,tenantId,processDefinitionKey
	ProcessInstanceId         string   `json:"processInstanceId,omitempty"`
	ProcessInstanceIds        []string `json:"processInstanceIds,omitempty"`
	ProcessDefinitionKey      string   `json:"processDefinitionKey,omitempty"`
	ProcessDefinitionId       string   `json:"processDefinitionId,omitempty"`
	BusinessKey               string   `json:"businessKey,omitempty"`
	InvolvedUser              string   `json:"involvedUser,omitempty"`
	Suspended                 *bool    `json:"suspended,omitempty"`
	SuperProcessInstanceId    string   `json:"superProcessInstanceId,omitempty"`
	SubProcessInstanceId      string   `json:"subProcessInstanceId,omitempty"`
	ExcludeSubprocesses       *bool    `json:"excludeSubprocesses,omitempty"`
	IncludeProcessVariables   *bool    `json:"includeProcessVariables,omitempty"`
	ParentId                  string   `json:"parentId,omitempty"`
	ChildOnly                 *bool    `json:"childOnly,omitempty"`
	request.WithTenant
}
