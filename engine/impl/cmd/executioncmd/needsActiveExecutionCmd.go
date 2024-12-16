package executioncmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

var _ engine.Command = (*NeedsActiveExecutionCmd)(nil)

type NeedsActiveExecutionCmd struct {
	IExecutionCmd
	Ctx           context.Context
	ExecutionId   string
	Transactional bool
}

func (n NeedsActiveExecutionCmd) IsTransactional() bool {
	return n.Transactional
}

func (n NeedsActiveExecutionCmd) Execute(commandContext engine.Context) (interface{}, error) {
	executionEntityManager := entitymanager.GetExecutionEntityManager()
	executionEntity, err := executionEntityManager.FindById(n.ExecutionId)
	if err != nil {
		return nil, err
	}

	return n.InternalExecute(commandContext, &executionEntity)
}

func (n NeedsActiveExecutionCmd) Context() context.Context {
	return n.Ctx
}

type Options func(*NeedsActiveExecutionCmd)

func WithTransactional(transactional bool) Options {
	return func(cmd *NeedsActiveExecutionCmd) {
		cmd.Transactional = transactional
	}
}

func WithContext(ctx context.Context) Options {
	return func(cmd *NeedsActiveExecutionCmd) {
		cmd.Ctx = ctx
	}
}
