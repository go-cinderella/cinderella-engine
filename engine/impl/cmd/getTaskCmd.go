package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
)

var _ engine.Command = (*GetTaskCmd)(nil)

type GetTaskCmd struct {
	UserId        string
	GroupId       string
	Ctx           context.Context
	Transactional bool
}

func (getTaskCmd GetTaskCmd) IsTransactional() bool {
	return getTaskCmd.Transactional
}

func (getTaskCmd GetTaskCmd) Context() context.Context {
	return getTaskCmd.Ctx
}

func (getTaskCmd GetTaskCmd) Execute(ctx engine.Context) (interface{}, error) {
	taskDataManager := datamanager.GetTaskDataManager()
	taskResult, err := taskDataManager.QueryUndoTask(getTaskCmd.UserId, getTaskCmd.GroupId)
	if err != nil {
		return taskResult, err
	}
	return taskResult, err
}
