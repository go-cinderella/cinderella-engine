package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historictask"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/samber/lo"
	"math"
	"slices"
	"sort"
)

var _ engine.Command = (*GetRecallableActivitiesCmd)(nil)

type GetRecallableActivitiesCmd struct {
	ProcessInstanceId string
	UserId            *string
	Groups            []string
	Ctx               context.Context
	Transactional     bool
}

func (g GetRecallableActivitiesCmd) IsTransactional() bool {
	return g.Transactional
}

type SelectOption struct {
	Label string
	Value string
}

type SelectOptions struct {
	Options []*SelectOption
}

func (g GetRecallableActivitiesCmd) Execute(commandContext engine.Context) (interface{}, error) {
	result := &SelectOptions{
		Options: make([]*SelectOption, 0),
	}

	historicProcessInstanceEntityManager := entitymanager.GetHistoricProcessInstanceEntityManager()
	historicProcessInstanceEntity, err := historicProcessInstanceEntityManager.FindById(g.ProcessInstanceId)
	if err != nil {
		return nil, err
	}

	if historicProcessInstanceEntity.EndTime != nil {
		return result, nil
	}

	hasPermission := true
	if g.UserId != nil {
		getIdentityLinksForProcessInstanceResult, err := GetIdentityLinksForProcessInstanceCmd{
			ProcessInstanceId: g.ProcessInstanceId,
			Ctx:               g.Ctx,
		}.Execute(commandContext)
		if err != nil {
			return result, err
		}

		involvedUsers := getIdentityLinksForProcessInstanceResult.([]*entitymanager.HistoricIdentityLinkEntity)
		hasPermission = lo.CountBy(involvedUsers, func(item *entitymanager.HistoricIdentityLinkEntity) bool {
			if item.UserID == nil {
				return false
			}
			return *item.UserID == *g.UserId
		}) > 0
	}
	if !hasPermission {
		return result, nil
	}

	historicTaskInstanceEntityManager := entitymanager.GetHistoricTaskInstanceEntityManager()
	histTaskInsts, err := historicTaskInstanceEntityManager.List(historictask.ListRequest{
		ListCommonRequest: request.ListCommonRequest{
			Size:  math.MaxInt32,
			Sort:  "start",
			Order: "desc",
		},
		ProcessInstanceId: g.ProcessInstanceId,
		Finished:          lo.ToPtr(true),
	})
	if err != nil {
		return nil, err
	}

	if len(histTaskInsts) == 0 {
		return result, nil
	}

	getCurrentTasksByUserResult, err := GetCurrentTasksByUserCmd{
		Ctx:               g.Ctx,
		UserId:            g.UserId,
		Groups:            g.Groups,
		ProcessInstanceId: g.ProcessInstanceId,
	}.Execute(commandContext)
	if err != nil {
		return result, err
	}
	currentTasks := getCurrentTasksByUserResult.([]entitymanager.TaskEntity)
	// 如果查到了当前登录用户有权限操作的任务，则过滤掉
	lo.ForEach(currentTasks, func(currentTask entitymanager.TaskEntity, index int) {
		histTaskInsts = lo.Filter(histTaskInsts, func(item entitymanager.HistoricTaskInstanceEntity, index int) bool {
			return item.TaskDefinitionKey != currentTask.TaskDefinitionKey
		})
	})
	// uniqTaskActivityIds only for sort
	uniqTaskActivityIds := lo.Uniq(lo.Map[entitymanager.HistoricTaskInstanceEntity, string](histTaskInsts, func(item entitymanager.HistoricTaskInstanceEntity, index int) string {
		return item.TaskDefinitionKey
	}))
	taskInstMap := lo.SliceToMap[entitymanager.HistoricTaskInstanceEntity, string, *SelectOption](histTaskInsts, func(item entitymanager.HistoricTaskInstanceEntity) (string, *SelectOption) {
		return item.TaskDefinitionKey, &SelectOption{
			Value: item.TaskDefinitionKey,
			Label: item.Name,
		}
	})
	result.Options = lo.Values(taskInstMap)
	sort.SliceStable(result.Options, func(i, j int) bool {
		optionI := result.Options[i]
		optionJ := result.Options[j]
		return slices.Index(uniqTaskActivityIds, optionI.Value) < slices.Index(uniqTaskActivityIds, optionJ.Value)
	})
	return result, nil
}

func (g GetRecallableActivitiesCmd) Context() context.Context {
	return g.Ctx
}
