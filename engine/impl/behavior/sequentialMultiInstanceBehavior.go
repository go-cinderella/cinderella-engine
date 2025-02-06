package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ IMultiInstanceActivityBehavior = (*SequentialMultiInstanceBehavior)(nil)

type SequentialMultiInstanceBehavior struct {
	AbstractMultiInstanceActivityBehavior
}

func (s SequentialMultiInstanceBehavior) CreateInstances(multiInstanceRootExecution delegate.DelegateExecution) (int, error) {
	nrOfInstances, err := s.resolveNrOfInstances(multiInstanceRootExecution)
	if err != nil {
		return 0, err
	}

	if nrOfInstances == 0 {
		return 0, nil
	}

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	execution, err := executionEntityManager.CreateChildExecution(multiInstanceRootExecution)
	if err != nil {
		return 0, err
	}
	execution.SetCurrentFlowElement(multiInstanceRootExecution.GetCurrentFlowElement())

	if err = s.SetLoopVariable(multiInstanceRootExecution, NUMBER_OF_INSTANCES, nrOfInstances); err != nil {
		return 0, err
	}
	if err = s.SetLoopVariable(multiInstanceRootExecution, NUMBER_OF_COMPLETED_INSTANCES, 0); err != nil {
		return 0, err
	}
	if err = s.SetLoopVariable(multiInstanceRootExecution, NUMBER_OF_ACTIVE_INSTANCES, 1); err != nil {
		return 0, err
	}

	if err = s.ExecuteOriginalBehavior(&execution, multiInstanceRootExecution, 0); err != nil {
		return 0, err
	}

	return nrOfInstances, nil
}

func (s SequentialMultiInstanceBehavior) Leave(execution delegate.DelegateExecution) error {
	//TODO implement me
	panic("implement me")
}
