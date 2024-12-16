package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/dto/task"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"math"
)

var _ engine.Command = (*GetCurrentTasksByUserCmd)(nil)

type GetCurrentTasksByUserCmd struct {
	Ctx               context.Context
	ProcessInstanceId string
	UserId            *string
	Groups            []string
	Transactional     bool
}

func (g GetCurrentTasksByUserCmd) IsTransactional() bool {
	return g.Transactional
}

func (g GetCurrentTasksByUserCmd) Execute(commandContext engine.Context) (interface{}, error) {
	if g.UserId == nil && len(g.Groups) == 0 {
		// root/admin user can access all runtime tasks
		req := task.ListRequest{
			ListCommonRequest: request.ListCommonRequest{
				Size: math.MaxInt32,
			},
			ProcessInstanceId: g.ProcessInstanceId,
		}
		taskEntityManager := entitymanager.GetTaskEntityManager()
		tasks, err := taskEntityManager.List(req)
		if err != nil {
			return nil, err
		}
		return tasks, nil
	}

	currentTasks := make([]entitymanager.TaskEntity, 0)

	req := task.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size: math.MaxInt32,
		},
		ProcessInstanceId:   g.ProcessInstanceId,
		CandidateOrAssigned: *g.UserId,
	}
	taskEntityManager := entitymanager.GetTaskEntityManager()
	tasks, err := taskEntityManager.List(req)
	if err != nil {
		return nil, err
	}

	currentTasks = append(currentTasks, tasks...)
	if len(g.Groups) > 0 {
		/**
		SELECT RES.*
		from ACT_RU_TASK RES
		WHERE RES.PROC_INST_ID_ = 'aa48d20a-99f9-11ef-be02-fe7c28f99a6d'
		  and RES.ASSIGNEE_ is null
		  and exists(select LINK.ID_ from ACT_RU_IDENTITYLINK LINK where LINK.TYPE_ = 'candidate' and LINK.TASK_ID_ = RES.ID_ and ((LINK.GROUP_ID_ IN ('role_dev'))))
		order by RES.ID_ asc
		LIMIT 2147483647 OFFSET 0;
		*/
		req = task.ListRequest{
			ListCommonRequest: request.ListCommonRequest{
				Size: math.MaxInt32,
			},
			ProcessInstanceId: g.ProcessInstanceId,
			CandidateGroupIn:  g.Groups,
		}
		tasks, err = taskEntityManager.List(req)
		if err != nil {
			return nil, err
		}
		currentTasks = append(currentTasks, tasks...)
	}
	return currentTasks, nil
}

func (g GetCurrentTasksByUserCmd) Context() context.Context {
	return g.Ctx
}
