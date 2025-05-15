package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/dto/historictask"
	"github.com/go-cinderella/cinderella-engine/engine/dto/request"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
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

/**
 * 执行获取可回撤活动命令，返回流程实例中可回撤的历史任务节点
 *
 * 此方法执行以下操作：
 * 1. 检查流程实例是否已结束
 * 2. 验证用户是否有权限操作该流程实例
 * 3. 获取流程实例的已完成历史任务
 * 4. 排除当前用户正在处理的任务
 * 5. 排除在当前流程定义中已不存在的节点（处理流程迁移情况）
 * 6. 对结果进行去重和排序
 *
 * @param commandContext 命令上下文
 * @return 可回撤活动的下拉选项列表和可能的错误
 */
func (g GetRecallableActivitiesCmd) Execute(commandContext engine.Context) (interface{}, error) {
	result := &SelectOptions{
		Options: make([]*SelectOption, 0),
	}

	// 获取历史流程实例
	historicProcessInstanceEntityManager := entitymanager.GetHistoricProcessInstanceEntityManager()
	historicProcessInstanceEntity, err := historicProcessInstanceEntityManager.FindById(g.ProcessInstanceId)
	if err != nil {
		return nil, err
	}

	// 如果流程已结束，返回空结果
	if historicProcessInstanceEntity.EndTime != nil {
		return result, nil
	}

	// 权限检查
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

	// 获取历史任务实例
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

	// 获取当前任务
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
	
	// 过滤掉当前用户正在操作的任务
	lo.ForEach(currentTasks, func(currentTask entitymanager.TaskEntity, index int) {
		histTaskInsts = lo.Filter(histTaskInsts, func(item entitymanager.HistoricTaskInstanceEntity, index int) bool {
			return item.TaskDefinitionKey != currentTask.TaskDefinitionKey
		})
	})
	
	// 获取当前流程定义，以处理流程迁移情况
	processDefinitionEntityManager := entitymanager.GetProcessDefinitionEntityManager()
	processDefinitionEntity, err := processDefinitionEntityManager.FindProcessDefinitionById(historicProcessInstanceEntity.ProcessDefinitionId)
	if err != nil {
		return nil, err
	}
	
	// 解析流程模型，获取当前流程定义中的所有节点
	bpmnXMLConverter := converter.BpmnXMLConverter{}
	bpmnModel := bpmnXMLConverter.ConvertToBpmnModel(processDefinitionEntity.ResourceContent)
	process := bpmnModel.GetProcess()
	
	// 过滤掉在当前流程定义中已不存在的节点（处理流程迁移情况）
	histTaskInsts = lo.Filter(histTaskInsts, func(item entitymanager.HistoricTaskInstanceEntity, index int) bool {
		// 检查节点是否在当前流程定义中存在
		flowElement := process.GetFlowElement(item.TaskDefinitionKey)
		return flowElement != nil
	})
	
	if len(histTaskInsts) == 0 {
		return result, nil
	}
	
	// 对结果进行去重和排序
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
