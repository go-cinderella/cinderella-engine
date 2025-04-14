package delegate

// DelegateExecution 执行实例接口
type DelegateExecution interface {
	VariableScope
	// Deprecated:
	GetId() string
	// Deprecated:
	SetId(id string)

	SetBusinessKey(businessKey string)

	GetCurrentFlowElement() FlowElement

	SetCurrentFlowElement(flow FlowElement)

	GetDeploymentId() string

	SetDeploymentId(deploymentId string)

	GetProcessInstanceId() string

	SetProcessInstanceId(processInstanceId string)

	GetProcessDefinitionId() string

	SetProcessDefinitionId(processDefineId string)

	GetCurrentActivityId() string

	SetCurrentActivityId(currentActivityId string)

	GetExecutionId() string

	GetParent() (DelegateExecution, error)

	GetParentId() string

	GetTenantId() *string

	// IsMultiInstanceRoot returns whether this execution is the root of a multi instance execution.
	IsMultiInstanceRoot() bool

	// SetMultiInstanceRoot changes whether this execution is a multi instance root or not.
	SetMultiInstanceRoot(isMultiInstanceRoot bool)
	
	SetActive(active bool) error
	IsActive() bool
}
