package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/dto/task"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/internal/datamanager"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/stringutils"
	"math"
	"slices"
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

	taskDataManager := datamanager.GetTaskDataManager()
	req := task.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size: math.MaxInt32,
		},
		ProcessInstanceId: g.ProcessInstanceId,
	}
	uniqTaskDefKeys, err := taskDataManager.GetUniqTaskDefKeys(req)
	if err != nil {
		return nil, err
	}
	if len(uniqTaskDefKeys) == 0 {
		return currentTasks, nil
	}

	processUtils := utils.ProcessDefinitionUtil{}
	var process model.Process
	process, err = processUtils.GetProcess(*uniqTaskDefKeys[0].ProcDefID_)
	if err != nil {
		return nil, err
	}

	var taskDefKeys []string
	var assigneeOrCandidateUsersElKeys []string
	var candidateGroupsElKeys []string

	for _, item := range uniqTaskDefKeys {
		flowElement := process.GetFlowElement(*item.TaskDefKey_)
		userTask := flowElement.(*model.UserTask)

		assigneeEl := userTask.Assignee
		if assigneeEl != nil && !utils.IsExpr(*assigneeEl) {
			// Assignee is literal
			assignee := *assigneeEl
			if g.UserId != nil && *g.UserId == assignee {
				taskDefKeys = append(taskDefKeys, *item.TaskDefKey_)
				continue
			}
		} else if assigneeEl != nil && utils.IsExpr(*assigneeEl) {
			assigneeOrCandidateUsersElKeys = append(assigneeOrCandidateUsersElKeys, *item.TaskDefKey_)
		}

		candidateUsersEl := userTask.CandidateUsers
		if candidateUsersEl != nil && !utils.IsExpr(*candidateUsersEl) {
			candidateUsersStr := *candidateUsersEl
			candidateUsers := stringutils.Split(candidateUsersStr, ",")
			if g.UserId != nil && slices.Contains(candidateUsers, *g.UserId) {
				taskDefKeys = append(taskDefKeys, *item.TaskDefKey_)
				continue
			}

			if g.UserId != nil {
				_, hasExpr := lo.Find(candidateUsers, func(item string) bool {
					return utils.IsExpr(item)
				})
				if hasExpr {
					assigneeOrCandidateUsersElKeys = append(assigneeOrCandidateUsersElKeys, *item.TaskDefKey_)
				}
			}
		} else if candidateUsersEl != nil && utils.IsExpr(*candidateUsersEl) {
			assigneeOrCandidateUsersElKeys = append(assigneeOrCandidateUsersElKeys, *item.TaskDefKey_)
		}

		candidateGroupsEl := userTask.CandidateGroups
		if candidateGroupsEl != nil && !utils.IsExpr(*candidateGroupsEl) {
			candidateGroupsStr := *candidateGroupsEl
			candidateGroups := stringutils.Split(candidateGroupsStr, ",")
			if len(g.Groups) > 0 && len(lo.Intersect(candidateGroups, g.Groups)) > 0 {
				taskDefKeys = append(taskDefKeys, *item.TaskDefKey_)
				continue
			}

			if len(g.Groups) > 0 {
				_, hasExpr := lo.Find(candidateGroups, func(item string) bool {
					return utils.IsExpr(item)
				})
				if hasExpr {
					candidateGroupsElKeys = append(candidateGroupsElKeys, *item.TaskDefKey_)
				}
			}
		} else if candidateGroupsEl != nil && utils.IsExpr(*candidateGroupsEl) {
			candidateGroupsElKeys = append(candidateGroupsElKeys, *item.TaskDefKey_)
		}
	}

	taskEntityManager := entitymanager.GetTaskEntityManager()

	if len(taskDefKeys) > 0 {
		req = task.ListRequest{
			ListCommonRequest: request.ListCommonRequest{
				Size: math.MaxInt32,
			},
			ProcessInstanceId:  g.ProcessInstanceId,
			TaskDefinitionKeys: taskDefKeys,
		}
		tasks, err := taskEntityManager.List(req)
		if err != nil {
			return nil, err
		}
		currentTasks = append(currentTasks, tasks...)
	}

	if len(assigneeOrCandidateUsersElKeys) > 0 {
		req = task.ListRequest{
			ListCommonRequest: request.ListCommonRequest{
				Size: math.MaxInt32,
			},
			ProcessInstanceId:   g.ProcessInstanceId,
			CandidateOrAssigned: *g.UserId,
			TaskDefinitionKeys:  assigneeOrCandidateUsersElKeys,
		}
		tasks, err := taskEntityManager.List(req)
		if err != nil {
			return nil, err
		}

		currentTasks = append(currentTasks, tasks...)
	}

	if len(g.Groups) > 0 && len(candidateGroupsElKeys) > 0 {
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
			ProcessInstanceId:  g.ProcessInstanceId,
			CandidateGroupIn:   g.Groups,
			TaskDefinitionKeys: candidateGroupsElKeys,
		}
		tasks, err := taskEntityManager.List(req)
		if err != nil {
			return nil, err
		}
		currentTasks = append(currentTasks, tasks...)
	}

	currentTasks = lo.UniqBy[entitymanager.TaskEntity, string, []entitymanager.TaskEntity](currentTasks, func(item entitymanager.TaskEntity) string {
		return item.Id
	})

	return currentTasks, nil
}

func (g GetCurrentTasksByUserCmd) Context() context.Context {
	return g.Ctx
}
