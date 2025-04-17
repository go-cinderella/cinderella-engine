package cmd

import (
	"context"
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/dto/execution"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/cmd/executioncmd"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
)

const (
	NUMBER_OF_INSTANCES = "nrOfInstances"
)

var _ executioncmd.IExecutionCmd = (*AddMultiInstanceExecutionCmd)(nil)

// AddMultiInstanceExecutionCmd 用于向多实例活动添加新的执行实例
type AddMultiInstanceExecutionCmd struct {
	executioncmd.NeedsActiveExecutionCmd
	ActivityId         string
	ParentExecutionId  string
	ExecutionVariables map[string]interface{}
}

// InternalExecute 实现 IExecutionCmd 接口
func (cmd *AddMultiInstanceExecutionCmd) InternalExecute(commandContext engine.Context, executionEntity delegate.DelegateExecution) (interface{}, error) {
	executionEntityManager := entitymanager.GetExecutionEntityManager()

	// 查找多实例根执行实例
	miExecution, err := cmd.searchForMultiInstanceActivity(cmd.ActivityId, cmd.ParentExecutionId, executionEntityManager)
	if err != nil {
		return nil, err
	}

	if miExecution == nil {
		return nil, errors.New("未找到活动ID为 " + cmd.ActivityId + " 的多实例执行实例")
	}

	// 创建新的子执行实例
	childExecution := entitymanager.CreateChildExecution(miExecution)
	childExecution.SetCurrentFlowElement(miExecution.GetCurrentFlowElement())

	if err := executionEntityManager.CreateExecution(&childExecution); err != nil {
		return nil, err
	}

	// 获取多实例活动元素
	processUtils := utils.ProcessDefinitionUtil{}
	process, err := processUtils.GetProcess(miExecution.GetProcessDefinitionId())
	if err != nil {
		return nil, err
	}

	miActivityElement := process.GetFlowElement(miExecution.GetCurrentActivityId())
	if miActivityElement == nil {
		return nil, errors.New("未找到活动元素: " + miExecution.GetCurrentActivityId())
	}

	// 检查是否为有效的多实例活动
	activity, ok := miActivityElement.(*model.Activity)
	if !ok {
		return nil, errors.New("流程元素不是有效的活动: " + miExecution.GetCurrentActivityId())
	}

	// 更新多实例计数
	currentNumberOfInstances, err := cmd.getLoopVariable(miExecution, behavior.NUMBER_OF_INSTANCES)
	if err != nil {
		return nil, err
	}

	if err := miExecution.SetVariableLocal(behavior.NUMBER_OF_INSTANCES, currentNumberOfInstances+1); err != nil {
		return nil, err
	}

	// 设置执行变量
	if cmd.ExecutionVariables != nil {
		if err := childExecution.SetVariablesLocal(cmd.ExecutionVariables); err != nil {
			return nil, err
		}
	}

	// 处理并行多实例
	loopCharacteristics := activity.GetLoopCharacteristics()
	if loopCharacteristics != nil && !loopCharacteristics.IsSequential {
		// 设置多实例根执行状态
		if err := miExecution.SetActive(true); err != nil {
			return nil, err
		}

		// 初始化子执行实例
		childExecution.SetCurrentFlowElement(miActivityElement)
		contextutil.GetAgendaFromContext(commandContext).PlanContinueMultiInstanceOperation(&childExecution, miExecution, currentNumberOfInstances)
	}

	return &childExecution, nil
}

// searchForMultiInstanceActivity 递归搜索多实例活动的执行实例
func (cmd *AddMultiInstanceExecutionCmd) searchForMultiInstanceActivity(activityId string, parentExecutionId string, executionEntityManager *entitymanager.ExecutionEntityManager) (delegate.DelegateExecution, error) {
	children, err := executionEntityManager.List(execution.ListRequest{
		ParentId: parentExecutionId,
	})
	if err != nil {
		return nil, err
	}

	var miExecution delegate.DelegateExecution

	for _, child := range children {
		if activityId == child.GetCurrentActivityId() && child.IsMultiInstanceRoot() {
			if miExecution != nil {
				return nil, errors.New("在 " + child.GetExecutionId() + " 中找到多个多实例执行实例")
			}
			miExecution = &child
		}

		// 递归搜索子执行
		childMiExecution, err := cmd.searchForMultiInstanceActivity(activityId, child.GetExecutionId(), executionEntityManager)
		if err != nil {
			return nil, err
		}

		if childMiExecution != nil {
			if miExecution != nil {
				return nil, errors.New("在 " + child.GetExecutionId() + " 中找到多个多实例执行实例")
			}
			miExecution = childMiExecution
		}
	}

	return miExecution, nil
}

// getLoopVariable 获取循环变量
func (cmd *AddMultiInstanceExecutionCmd) getLoopVariable(execution delegate.DelegateExecution, variableName string) (int, error) {
	value, ok, err := execution.GetVariableLocal(variableName)
	if err != nil {
		return 0, err
	}

	if ok {
		return value.(int), nil
	}

	return 0, nil
}

// NewAddMultiInstanceExecutionCmd 创建添加多实例执行命令
func NewAddMultiInstanceExecutionCmd(ctx context.Context, activityId string, parentExecutionId string, executionVariables map[string]interface{}, options ...executioncmd.Options) *AddMultiInstanceExecutionCmd {
	cmd := &AddMultiInstanceExecutionCmd{
		ActivityId:         activityId,
		ParentExecutionId:  parentExecutionId,
		ExecutionVariables: executionVariables,
	}

	cmd.NeedsActiveExecutionCmd = executioncmd.NeedsActiveExecutionCmd{
		IExecutionCmd: cmd,
		ExecutionId:   parentExecutionId,
		Ctx:           ctx,
		Transactional: true,
	}

	for _, option := range options {
		option(&cmd.NeedsActiveExecutionCmd)
	}

	return cmd
}
