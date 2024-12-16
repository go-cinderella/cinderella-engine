package historicprocess

import (
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/dto/variable"
	"time"
)

type ListRequest struct {
	request.ListCommonRequest // order allow processInstanceId,processDefinitionId,businessKey,startTime,endTime,duration,tenantId
	request.WithTenant
	ProcessInstanceId                 string                            `json:"processInstanceId,omitempty"`
	ProcessInstanceIds                []string                          `json:"processInstanceIds,omitempty"`
	ProcessInstanceName               string                            `json:"processInstanceName,omitempty"`
	ProcessInstanceNameLike           string                            `json:"processInstanceNameLike,omitempty"`
	ProcessInstanceNameLikeIgnoreCase string                            `json:"processInstanceNameLikeIgnoreCase,omitempty"`
	ProcessBusinessKey                string                            `json:"processBusinessKey,omitempty"`
	ProcessBusinessKeyLike            string                            `json:"processBusinessKeyLike,omitempty"`
	ProcessDefinitionKey              string                            `json:"processDefinitionKey,omitempty"`
	ProcessDefinitionKeyIn            []string                          `json:"processDefinitionKeyIn,omitempty"`
	ProcessDefinitionKeyNotIn         []string                          `json:"processDefinitionKeyNotIn,omitempty"`
	ProcessDefinitionName             string                            `json:"processDefinitionName,omitempty"`
	ProcessDefinitionVersion          int                               `json:"processDefinitionVersion,omitempty"`
	ProcessDefinitionCategory         string                            `json:"processDefinitionCategory,omitempty"`
	ProcessDefinitionId               string                            `json:"processDefinitionId,omitempty"`
	DeploymentId                      string                            `json:"deploymentId,omitempty"`
	DeploymentIdIn                    []string                          `json:"deploymentIdIn,omitempty"`
	InvolvedUser                      string                            `json:"involvedUser,omitempty"`
	Finished                          *bool                             `json:"finished,omitempty"`
	SuperProcessInstanceId            string                            `json:"superProcessInstanceId,omitempty"`
	ExcludeSubprocesses               *bool                             `json:"excludeSubprocesses,omitempty"`
	FinishedAfter                     *time.Time                        `json:"finishedAfter,omitempty"`
	FinishedBefore                    *time.Time                        `json:"finishedBefore,omitempty"`
	StartedAfter                      *time.Time                        `json:"startedAfter,omitempty"`
	StartedBefore                     *time.Time                        `json:"startedBefore,omitempty"`
	StartedBy                         string                            `json:"startedBy,omitempty"`
	IncludeProcessVariables           *bool                             `json:"includeProcessVariables,omitempty"`
	CallbackId                        string                            `json:"callbackId,omitempty"`
	CallbackType                      string                            `json:"callbackType,omitempty"`
	Variables                         *[]variable.VariableSearchRequest `json:"variable,omitempty"`
}
