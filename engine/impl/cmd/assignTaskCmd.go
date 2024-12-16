package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/cmd/taskcmd"
	"github.com/samber/lo"
	"time"
)

var _ engine.Command = (*AssignTaskCmd)(nil)

type AssignTaskCmd struct {
	taskcmd.NeedsActiveTaskCmd
	UserId *string
}

func (assignTaskCmd AssignTaskCmd) TaskExecute(commandContext engine.Context, taskEntity entitymanager.TaskEntity) (interface{}, error) {
	taskEntityManager := entitymanager.GetTaskEntityManager()

	if assignTaskCmd.UserId != nil {
		if taskEntity.GetAssignee() != nil {
			if *taskEntity.GetAssignee() == *assignTaskCmd.UserId {
				return nil, nil
			}
		}

		taskEntity.SetClaimTime(lo.ToPtr(time.Now().UTC()))
	} else {
		if taskEntity.GetAssignee() == nil {
			return nil, nil
		}
		taskEntity.SetClaimTime(nil)
	}

	if err := taskEntityManager.ChangeTaskAssignee(taskEntity, assignTaskCmd.UserId); err != nil {
		return nil, err
	}

	return nil, nil
}

func (assignTaskCmd AssignTaskCmd) Context() context.Context {
	return assignTaskCmd.Ctx
}

func NewAssignTaskCmd(ctx context.Context, taskId string, userId *string, options ...taskcmd.Options) AssignTaskCmd {
	assignTaskCmd := AssignTaskCmd{
		UserId: userId,
	}
	assignTaskCmd.NeedsActiveTaskCmd = taskcmd.NeedsActiveTaskCmd{
		ITaskCmd: &assignTaskCmd,
		TaskId:   taskId,
		Ctx:      ctx,
	}

	for _, option := range options {
		option(&assignTaskCmd.NeedsActiveTaskCmd)
	}

	return assignTaskCmd
}
