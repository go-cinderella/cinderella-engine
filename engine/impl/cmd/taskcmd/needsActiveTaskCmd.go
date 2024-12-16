package taskcmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

var _ engine.Command = (*NeedsActiveTaskCmd)(nil)

type NeedsActiveTaskCmd struct {
	ITaskCmd
	Ctx           context.Context
	TaskId        string
	Transactional bool
}

func (needsActiveTaskCmd NeedsActiveTaskCmd) IsTransactional() bool {
	return needsActiveTaskCmd.Transactional
}

func (needsActiveTaskCmd NeedsActiveTaskCmd) Context() context.Context {
	return needsActiveTaskCmd.Ctx
}

func (needsActiveTaskCmd NeedsActiveTaskCmd) Execute(commandContext engine.Context) (interface{}, error) {
	taskEntityManager := entitymanager.GetTaskEntityManager()
	taskEntity, err := taskEntityManager.FindById(needsActiveTaskCmd.TaskId)
	if err != nil {
		return nil, err
	}
	execute, err := needsActiveTaskCmd.TaskExecute(commandContext, taskEntity)
	return execute, err
}

type Options func(*NeedsActiveTaskCmd)

func WithTransactional(transactional bool) Options {
	return func(cmd *NeedsActiveTaskCmd) {
		cmd.Transactional = transactional
	}
}

func WithContext(ctx context.Context) Options {
	return func(cmd *NeedsActiveTaskCmd) {
		cmd.Ctx = ctx
	}
}
