package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/zlogger"
)

var _ multiInstanceActivityBehavior = (*SequentialMultiInstanceBehavior)(nil)

type SequentialMultiInstanceBehavior struct {
	AbstractMultiInstanceActivityBehavior
}

func (s SequentialMultiInstanceBehavior) createInstances(multiInstanceRootExecution delegate.DelegateExecution) (int, error) {
	nrOfInstances, err := s.resolveNrOfInstances(multiInstanceRootExecution)
	if err != nil {
		return 0, err
	}

	if nrOfInstances == 0 {
		return 0, nil
	}

	execution := entitymanager.CreateChildExecution(multiInstanceRootExecution)
	execution.SetCurrentFlowElement(multiInstanceRootExecution.GetCurrentFlowElement())

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	if err = executionEntityManager.CreateExecution(&execution); err != nil {
		return 0, err
	}

	if err = s.setLoopVariable(multiInstanceRootExecution, NUMBER_OF_INSTANCES, nrOfInstances); err != nil {
		return 0, err
	}
	if err = s.setLoopVariable(multiInstanceRootExecution, NUMBER_OF_COMPLETED_INSTANCES, 0); err != nil {
		return 0, err
	}
	if err = s.setLoopVariable(multiInstanceRootExecution, NUMBER_OF_ACTIVE_INSTANCES, 1); err != nil {
		return 0, err
	}

	if err = s.ExecuteOriginalBehavior(&execution, multiInstanceRootExecution, 0); err != nil {
		return 0, err
	}

	return nrOfInstances, nil
}

func (s SequentialMultiInstanceBehavior) leave(execution delegate.DelegateExecution) error {
	multiInstanceRootExecution, err := s.getMultiInstanceRootExecution(execution)
	if err != nil {
		return err
	}

	loopCounter, err := s.getLoopVariable(execution, CollectionElementIndexVariable)
	if err != nil {
		return err
	}
	loopCounter = loopCounter + 1

	nrOfInstances, err := s.getLoopVariable(multiInstanceRootExecution, NUMBER_OF_INSTANCES)
	if err != nil {
		return err
	}

	nrOfCompletedInstances, err := s.getLoopVariable(multiInstanceRootExecution, NUMBER_OF_COMPLETED_INSTANCES)
	if err != nil {
		return err
	}

	nrOfCompletedInstances = nrOfCompletedInstances + 1

	nrOfActiveInstances, err := s.getLoopVariable(multiInstanceRootExecution, NUMBER_OF_ACTIVE_INSTANCES)
	if err != nil {
		return err
	}

	zlogger.Info().Msgf("Multi-instance %s(%s). Details: loopCounter=%d, nrOrCompletedInstances=%d,nrOfActiveInstances=%d,nrOfInstances=%d",
		execution.GetCurrentFlowElement().GetName(), execution.GetCurrentFlowElement().GetId(), loopCounter, nrOfCompletedInstances, nrOfActiveInstances, nrOfInstances)

	if err = s.setLoopVariable(multiInstanceRootExecution, NUMBER_OF_COMPLETED_INSTANCES, nrOfCompletedInstances); err != nil {
		return err
	}

	completeConditionSatisfied, err := s.completionConditionSatisfied(multiInstanceRootExecution)
	if err != nil {
		return err
	}

	if loopCounter >= nrOfInstances || completeConditionSatisfied {
		// Call logic to leave activity instead of continuing multi-instance
		return s.AbstractMultiInstanceActivityBehavior.leave(execution)
	} else {
		return s.continueSequentialMultiInstance(execution, loopCounter, multiInstanceRootExecution)
	}
}

func (s SequentialMultiInstanceBehavior) continueSequentialMultiInstance(execution delegate.DelegateExecution, loopCounter int, multiInstanceRootExecution delegate.DelegateExecution) error {
	dataManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	err := dataManager.RecordActEndByExecutionId(execution.GetCurrentFlowElement(), execution, nil)
	if err != nil {
		return err
	}

	variablesLocal, err := execution.GetVariablesLocal()
	if err != nil {
		return err
	}
	delete(variablesLocal, NUMBER_OF_INSTANCES)
	delete(variablesLocal, NUMBER_OF_COMPLETED_INSTANCES)
	delete(variablesLocal, NUMBER_OF_ACTIVE_INSTANCES)

	err = execution.RemoveVariablesLocal(lo.Keys(variablesLocal))
	if err != nil {
		return err
	}

	return s.ExecuteOriginalBehavior(execution, multiInstanceRootExecution, loopCounter)
}

// ContinueSequentialMultiInstance 继续执行顺序多实例
// 提供给外部命令调用的公开方法
func (s SequentialMultiInstanceBehavior) ContinueSequentialMultiInstance(execution delegate.DelegateExecution, loopCounter int, multiInstanceRootExecution delegate.DelegateExecution) error {
	return s.continueSequentialMultiInstance(execution, loopCounter, multiInstanceRootExecution)
}

// GetLoopVariable 获取循环变量
// 提供给外部命令调用的公开方法
func (s SequentialMultiInstanceBehavior) GetLoopVariable(execution delegate.DelegateExecution, variableName string) (int, error) {
	return s.getLoopVariable(execution, variableName)
}
