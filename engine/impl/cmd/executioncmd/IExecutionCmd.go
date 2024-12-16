package executioncmd

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type IExecutionCmd interface {
	InternalExecute(CommandContext engine.Context, executionEntity delegate.DelegateExecution) (interface{}, error)
}
