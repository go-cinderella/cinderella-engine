package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/zlogger"
)

var _ multiInstanceActivityBehavior = (*ParallelMultiInstanceBehavior)(nil)

type ParallelMultiInstanceBehavior struct {
	AbstractMultiInstanceActivityBehavior
}

func (p ParallelMultiInstanceBehavior) createInstances(execution delegate.DelegateExecution) (int, error) {
	nrOfInstances, err := p.resolveNrOfInstances(execution)
	if err != nil {
		return 0, err
	}

	if nrOfInstances <= 0 {
		zlogger.Info().Msgf("无效的实例数量: 必须是非负整数, 但得到 %d", nrOfInstances)
		return 0, nil
	}

	if err = p.setLoopVariable(execution, NUMBER_OF_INSTANCES, nrOfInstances); err != nil {
		return 0, err
	}
	if err = p.setLoopVariable(execution, NUMBER_OF_COMPLETED_INSTANCES, 0); err != nil {
		return 0, err
	}
	if err = p.setLoopVariable(execution, NUMBER_OF_ACTIVE_INSTANCES, nrOfInstances); err != nil {
		return 0, err
	}

	// 先创建所有的执行实例
	concurrentExecutions := make([]delegate.DelegateExecution, 0, nrOfInstances)
	executionEntityManager := entitymanager.GetExecutionEntityManager()

	for loopCounter := 0; loopCounter < nrOfInstances; loopCounter++ {
		// 创建并行执行实例
		concurrentExecution := entitymanager.CreateChildExecution(execution)
		concurrentExecution.SetCurrentFlowElement(execution.GetCurrentFlowElement())

		if err := executionEntityManager.CreateExecution(&concurrentExecution); err != nil {
			return 0, err
		}

		// 设置循环计数器变量
		if err := p.setLoopVariable(&concurrentExecution, CollectionElementIndexVariable, loopCounter); err != nil {
			return 0, err
		}

		concurrentExecutions = append(concurrentExecutions, &concurrentExecution)
		zlogger.Info().Msgf("并行多实例 %s(%s) 初始化. 详情: loopCounter=%d, nrOfCompletedInstances=0, nrOfActiveInstances=%d, nrOfInstances=%d",
			execution.GetCurrentFlowElement().GetName(), execution.GetCurrentFlowElement().GetId(), loopCounter, nrOfInstances, nrOfInstances)
	}

	// 先创建所有执行实例后，再执行原始行为
	for loopCounter, concurrentExecution := range concurrentExecutions {
		if err := p.ExecuteOriginalBehavior(concurrentExecution, execution, loopCounter); err != nil {
			return 0, err
		}
	}

	// 将父执行设置为非活动状态，因为所有工作都在子执行中进行
	if len(concurrentExecutions) > 0 {
		if executionEntity, ok := execution.(*entitymanager.ExecutionEntity); ok {
			executionEntity.SetActive(false)
		}
	}

	return nrOfInstances, nil
}

func (p ParallelMultiInstanceBehavior) leave(execution delegate.DelegateExecution) error {
	multiInstanceRootExecution, err := p.getMultiInstanceRootExecution(execution)
	if err != nil {
		return err
	}

	nrOfInstances, err := p.getLoopVariable(multiInstanceRootExecution, NUMBER_OF_INSTANCES)
	if err != nil {
		return err
	}

	if nrOfInstances == 0 {
		// 空集合，直接离开
		return p.AbstractMultiInstanceActivityBehavior.leave(execution)
	}

	loopCounter, err := p.getLoopVariable(execution, CollectionElementIndexVariable)
	if err != nil {
		return err
	}

	nrOfCompletedInstances, err := p.getLoopVariable(multiInstanceRootExecution, NUMBER_OF_COMPLETED_INSTANCES)
	if err != nil {
		return err
	}
	nrOfCompletedInstances = nrOfCompletedInstances + 1

	nrOfActiveInstances, err := p.getLoopVariable(multiInstanceRootExecution, NUMBER_OF_ACTIVE_INSTANCES)
	if err != nil {
		return err
	}
	nrOfActiveInstances = nrOfActiveInstances - 1

	if err = p.setLoopVariable(multiInstanceRootExecution, NUMBER_OF_COMPLETED_INSTANCES, nrOfCompletedInstances); err != nil {
		return err
	}
	if err = p.setLoopVariable(multiInstanceRootExecution, NUMBER_OF_ACTIVE_INSTANCES, nrOfActiveInstances); err != nil {
		return err
	}

	zlogger.Info().Msgf("并行多实例 %s(%s) 实例完成. 详情: loopCounter=%d, nrOfCompletedInstances=%d, nrOfActiveInstances=%d, nrOfInstances=%d",
		execution.GetCurrentFlowElement().GetName(), execution.GetCurrentFlowElement().GetId(), loopCounter, nrOfCompletedInstances, nrOfActiveInstances, nrOfInstances)

	// 使当前执行变为非活动状态
	if executionEntity, ok := execution.(*entitymanager.ExecutionEntity); ok {
		executionEntity.SetActive(false)
	}

	parentExecution, err := execution.GetParent()
	if err != nil {
		return err
	}

	// 查找所有已完成的并行执行实例
	executionEntityManager := entitymanager.GetExecutionEntityManager()
	children, err := executionEntityManager.CollectChildren(parentExecution.GetExecutionId())
	if err != nil {
		return err
	}

	inactiveExecutions := lo.Filter(children, func(item entitymanager.ExecutionEntity, index int) bool {
		return !item.IsActive()
	})

	// 判断是否所有实例都已完成或满足完成条件
	completionConditionSatisfied, err := p.completionConditionSatisfied(multiInstanceRootExecution)
	if err != nil {
		return err
	}

	if len(inactiveExecutions) >= nrOfInstances || completionConditionSatisfied {
		// 移除所有仍然活跃的子执行实例（如果完成条件满足）
		activeExecutions := lo.Filter(children, func(item entitymanager.ExecutionEntity, index int) bool {
			return item.IsActive()
		})

		for _, activeExecution := range activeExecutions {
			zlogger.Info().Msgf("执行 %s 仍然活跃，但多实例已完成。移除此执行", activeExecution.GetExecutionId())
			if err := executionEntityManager.DeleteRelatedDataForExecution(activeExecution.GetExecutionId(), lo.ToPtr(DELETE_REASON_END)); err != nil {
				return err
			}
			if err := executionEntityManager.DeleteExecution(activeExecution.GetExecutionId()); err != nil {
				return err
			}
		}

		// 清理并继续执行流程
		return p.AbstractMultiInstanceActivityBehavior.leave(execution)
	}

	return nil
}
