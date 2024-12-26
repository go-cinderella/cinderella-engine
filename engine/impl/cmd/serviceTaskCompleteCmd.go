package cmd

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/impl/cmd/executioncmd"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ executioncmd.IExecutionCmd = (*ServiceTaskCompleteCmd)(nil)

type ServiceTaskCompleteCmd struct {
	executioncmd.NeedsActiveExecutionCmd
}

func (s ServiceTaskCompleteCmd) InternalExecute(CommandContext engine.Context, executionEntity delegate.DelegateExecution) (interface{}, error) {
	contextutil.GetAgenda().PlanTriggerExecutionOperation(executionEntity)
	return nil, nil
}

func NewServiceTaskCompleteCmd(executionId string, options ...executioncmd.Options) ServiceTaskCompleteCmd {
	serviceTaskCompleteCmd := ServiceTaskCompleteCmd{}
	serviceTaskCompleteCmd.NeedsActiveExecutionCmd = executioncmd.NeedsActiveExecutionCmd{
		IExecutionCmd: &serviceTaskCompleteCmd,
		ExecutionId:   executionId,
	}

	for _, option := range options {
		option(&serviceTaskCompleteCmd.NeedsActiveExecutionCmd)
	}

	return serviceTaskCompleteCmd
}
