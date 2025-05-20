package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/errs"
	"github.com/go-cinderella/cinderella-engine/engine/impl/cmd/taskcmd"
	"github.com/samber/lo"
	"time"
)

var _ engine.Command = (*ClaimTaskCmd)(nil)

type ClaimTaskCmd struct {
	taskcmd.NeedsActiveTaskCmd
	UserId *string
}

func (claimTaskCmd ClaimTaskCmd) TaskExecute(commandContext engine.Context, taskEntity entitymanager.TaskEntity) (interface{}, error) {
	taskEntityManager := entitymanager.GetTaskEntityManager()

	if claimTaskCmd.UserId != nil {
		if taskEntity.GetAssignee() != nil {
			if *taskEntity.GetAssignee() != *claimTaskCmd.UserId {
				return nil, errs.CinderellaError{Msg: "Task has already been claimed by another user"}
			}
			return nil, nil
		}

		taskEntity.SetClaimTime(lo.ToPtr(time.Now().UTC()))
		if err := taskEntityManager.ChangeTaskAssignee(taskEntity, claimTaskCmd.UserId); err != nil {
			return nil, err
		}
	} else {
		if taskEntity.GetAssignee() != nil {
			taskEntity.SetClaimTime(nil)

			if err := taskEntityManager.ChangeTaskAssignee(taskEntity, nil); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (claimTaskCmd ClaimTaskCmd) Context() context.Context {
	return claimTaskCmd.Ctx
}
