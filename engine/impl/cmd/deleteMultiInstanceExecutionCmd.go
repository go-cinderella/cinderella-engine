package cmd

import (
	"context"
	"errors"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/cmd/executioncmd"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
)

const (
	DELETE_REASON_MI_DELETE = "Delete MI execution"
)

var _ executioncmd.IExecutionCmd = (*DeleteMultiInstanceExecutionCmd)(nil)

// DeleteMultiInstanceExecutionCmd 用于删除多实例活动的执行实例
// 这是基于Java中的DeleteMultiInstanceExecutionCmd实现的
type DeleteMultiInstanceExecutionCmd struct {
	executioncmd.NeedsActiveExecutionCmd
	ExecutionId          string
	ExecutionIsCompleted bool
}

// InternalExecute 实现 IExecutionCmd 接口
func (cmd *DeleteMultiInstanceExecutionCmd) InternalExecute(commandContext engine.Context, execution delegate.DelegateExecution) (interface{}, error) {
	executionEntityManager := entitymanager.GetExecutionEntityManager()

	// 获取多实例活动元素
	processUtils := utils.ProcessDefinitionUtil{}
	process, err := processUtils.GetProcess(execution.GetProcessDefinitionId())
	if err != nil {
		return nil, err
	}

	miActivityElement := process.GetFlowElement(execution.GetCurrentActivityId())
	if miActivityElement == nil {
		return nil, errors.New("未找到活动元素: " + execution.GetCurrentActivityId())
	}

	// 检查是否为有效的多实例活动
	activity, ok := miActivityElement.(*model.Activity)
	if !ok {
		return nil, errors.New("流程元素不是有效的活动: " + execution.GetCurrentActivityId())
	}

	loopCharacteristics := activity.GetLoopCharacteristics()
	if loopCharacteristics == nil {
		return nil, errors.New("未找到多实例循环特性: " + execution.GetCurrentActivityId())
	}

	// 获取多实例根执行实例
	miExecution, err := cmd.getMultiInstanceRootExecution(execution)
	if err != nil {
		return nil, err
	}

	// 删除当前执行实例及其子实例
	// 直接使用GetExecutionEntityManager().CollectChildren方法获取子执行
	childExecutions, err := executionEntityManager.CollectChildren(execution.GetExecutionId())
	if err != nil {
		return nil, err
	}

	if err = executionEntityManager.DeleteChildExecutions(childExecutions, lo.ToPtr(DELETE_REASON_MI_DELETE)); err != nil {
		return nil, err
	}

	if err = executionEntityManager.DeleteRelatedDataForExecution(execution.GetExecutionId(), lo.ToPtr(DELETE_REASON_MI_DELETE)); err != nil {
		return nil, err
	}

	if err = executionEntityManager.DeleteExecution(execution.GetExecutionId()); err != nil {
		return nil, err
	}

	// 获取循环计数器
	loopCounter := 0
	if loopCharacteristics.IsSequential {
		// 顺序多实例获取循环计数器
		behaviorInstance, ok := activity.GetBehavior().(*behavior.SequentialMultiInstanceBehavior)
		if ok {
			value, err := behaviorInstance.GetLoopVariable(execution, behavior.CollectionElementIndexVariable)
			if err == nil {
				loopCounter = value
			}
		}
	}

	// 处理已完成的实例
	if cmd.ExecutionIsCompleted {
		completedValue, ok, err := miExecution.GetVariableLocal(behavior.NUMBER_OF_COMPLETED_INSTANCES)
		if err != nil {
			return nil, err
		}

		if ok {
			// 增加已完成实例的计数
			numberOfCompletedInstances := 0
			if completedValue != nil {
				numberOfCompletedInstances = completedValue.(int)
			}

			if err := miExecution.SetVariableLocal(behavior.NUMBER_OF_COMPLETED_INSTANCES, numberOfCompletedInstances+1); err != nil {
				return nil, err
			}
		}

		loopCounter++
	} else {
		// 减少实例总数
		instancesValue, ok, err := miExecution.GetVariableLocal(behavior.NUMBER_OF_INSTANCES)
		if err != nil {
			return nil, err
		}

		if ok && instancesValue != nil {
			currentNumberOfInstances := instancesValue.(int)
			if err := miExecution.SetVariableLocal(behavior.NUMBER_OF_INSTANCES, currentNumberOfInstances-1); err != nil {
				return nil, err
			}
		}
	}

	// 处理顺序多实例
	if loopCharacteristics.IsSequential {
		// 为顺序多实例创建新的子执行实例并继续执行
		childExecution := entitymanager.CreateChildExecution(miExecution)
		childExecution.SetCurrentFlowElement(miExecution.GetCurrentFlowElement())

		if err := executionEntityManager.CreateExecution(&childExecution); err != nil {
			return nil, err
		}

		// 继续执行顺序多实例
		behaviorInstance, ok := activity.GetBehavior().(*behavior.SequentialMultiInstanceBehavior)
		if ok {
			if err := behaviorInstance.ContinueSequentialMultiInstance(&childExecution, loopCounter, miExecution); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

// getMultiInstanceRootExecution 获取多实例根执行实例
func (cmd *DeleteMultiInstanceExecutionCmd) getMultiInstanceRootExecution(executionEntity delegate.DelegateExecution) (delegate.DelegateExecution, error) {
	var multiInstanceRootExecution delegate.DelegateExecution
	currentExecution := executionEntity

	for currentExecution != nil && multiInstanceRootExecution == nil && currentExecution.GetParentId() != "" {
		if currentExecution.IsMultiInstanceRoot() {
			multiInstanceRootExecution = currentExecution
		} else {
			var err error
			currentExecution, err = currentExecution.GetParent()
			if err != nil {
				return nil, err
			}
		}
	}

	return multiInstanceRootExecution, nil
}

// NewDeleteMultiInstanceExecutionCmd 创建删除多实例执行命令
func NewDeleteMultiInstanceExecutionCmd(ctx context.Context, executionId string, executionIsCompleted bool, options ...executioncmd.Options) *DeleteMultiInstanceExecutionCmd {
	cmd := &DeleteMultiInstanceExecutionCmd{
		ExecutionId:          executionId,
		ExecutionIsCompleted: executionIsCompleted,
	}

	cmd.NeedsActiveExecutionCmd = executioncmd.NeedsActiveExecutionCmd{
		IExecutionCmd: cmd,
		ExecutionId:   executionId,
		Ctx:           ctx,
		Transactional: true,
	}

	for _, option := range options {
		option(&cmd.NeedsActiveExecutionCmd)
	}

	return cmd
}
