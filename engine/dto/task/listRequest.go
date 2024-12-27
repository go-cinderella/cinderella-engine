package task

import (
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/dto/variable"
	"time"
)

type ListRequest struct {
	request.ListCommonRequest                                 // order allow id,name,description,dueDate,createTime,priority,executionId,processInstanceId,tenantId
	Name                           string                     `json:"name,omitempty"`
	NameLike                       string                     `json:"nameLike,omitempty"`
	Description                    string                     `json:"description,omitempty"`
	DescriptionLike                string                     `json:"descriptionLike,omitempty"`
	Priority                       int                        `json:"priority,omitempty"`
	MinimumPriority                int                        `json:"minimumPriority,omitempty"`
	MaximumPriority                int                        `json:"maximumPriority,omitempty"`
	Assignee                       string                     `json:"assignee,omitempty"`
	AssigneeLike                   string                     `json:"assigneeLike,omitempty"`
	Owner                          string                     `json:"owner,omitempty"`
	OwnerLike                      string                     `json:"ownerLike,omitempty"`
	Unassigned                     *bool                      `json:"unassigned,omitempty"`
	DelegationState                string                     `json:"delegationState,omitempty"`
	CandidateUser                  string                     `json:"candidateUser,omitempty"`
	CandidateGroup                 string                     `json:"candidateGroup,omitempty"`
	CandidateGroupIn               []string                   `json:"candidateGroupIn,omitempty"`
	InvolvedUser                   string                     `json:"involvedUser,omitempty"`
	ProcessInstanceId              string                     `json:"processInstanceId,omitempty"`
	TaskDefinitionKey              string                     `json:"taskDefinitionKey,omitempty"`
	TaskDefinitionKeys             []string                   `json:"taskDefinitionKeys,omitempty"`
	TaskDefinitionKeyLike          string                     `json:"taskDefinitionKeyLike,omitempty"`
	ProcessInstanceBusinessKey     string                     `json:"processInstanceBusinessKey,omitempty"`
	ProcessInstanceBusinessKeyLike string                     `json:"processInstanceBusinessKeyLike,omitempty"`
	ProcessDefinitionId            string                     `json:"processDefinitionId,omitempty"`
	ProcessDefinitionKey           string                     `json:"processDefinitionKey,omitempty"`
	ProcessDefinitionKeyLike       string                     `json:"processDefinitionKeyLike,omitempty"`
	ProcessDefinitionName          string                     `json:"processDefinitionName,omitempty"`
	ProcessDefinitionNameLike      string                     `json:"processDefinitionNameLike,omitempty"`
	ExecutionId                    string                     `json:"executionId,omitempty"`
	CreatedOn                      *time.Time                 `json:"createdOn,omitempty"`
	CreatedBefore                  *time.Time                 `json:"createdBefore,omitempty"`
	CreatedAfter                   *time.Time                 `json:"createdAfter,omitempty"`
	DueOn                          *time.Time                 `json:"dueOn,omitempty"`
	DueBefore                      *time.Time                 `json:"dueBefore,omitempty"`
	DueAfter                       *time.Time                 `json:"dueAfter,omitempty"`
	WithoutDueDate                 *time.Time                 `json:"withoutDueDate,omitempty"`
	ExcludeSubTasks                *bool                      `json:"excludeSubTasks,omitempty"`
	Active                         *bool                      `json:"active,omitempty"`
	IncludeTaskLocalVariables      *bool                      `json:"includeTaskLocalVariables,omitempty"`
	IncludeProcessVariables        *bool                      `json:"includeProcessVariables,omitempty"`
	CandidateOrAssigned            string                     `json:"candidateOrAssigned,omitempty"`
	Category                       string                     `json:"category,omitempty"`
	TaskVariables                  []variable.VariableRequest `json:"taskVariables,omitempty"`
	ProcessInstanceVariables       []variable.VariableRequest `json:"processInstanceVariables,omitempty"`
	request.WithTenant
}
