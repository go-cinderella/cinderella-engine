package cmd

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/impl/cmd/executioncmd"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ executioncmd.IExecutionCmd = (*EvaluateConditionalEventsCmd)(nil)

type EvaluateConditionalEventsCmd struct {
	executioncmd.NeedsActiveExecutionCmd
}

func (evaluateConditionalEventsCmd EvaluateConditionalEventsCmd) InternalExecute(CommandContext engine.Context, executionEntity delegate.DelegateExecution) (interface{}, error) {
	contextutil.GetAgenda().PlanEvaluateConditionalEventsOperation(executionEntity)
	return nil, nil
}

func NewEvaluateConditionalEventsCmd(executionId string, options ...executioncmd.Options) EvaluateConditionalEventsCmd {
	conditionalEventsCmd := EvaluateConditionalEventsCmd{}
	conditionalEventsCmd.NeedsActiveExecutionCmd = executioncmd.NeedsActiveExecutionCmd{
		IExecutionCmd: &conditionalEventsCmd,
		ExecutionId:   executionId,
	}

	for _, option := range options {
		option(&conditionalEventsCmd.NeedsActiveExecutionCmd)
	}

	return conditionalEventsCmd
}
