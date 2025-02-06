package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ContinueMultiInstanceOperation struct {
	AbstractOperation

	MultiInstanceRootExecution delegate.DelegateExecution
	LoopCounter                int
}

func (cont *ContinueMultiInstanceOperation) Run() (err error) {
	element := cont.Execution.GetCurrentFlowElement()
	if element != nil {
		if element.GetOutgoing() != nil {
			err = cont.continueThroughMultiInstanceFlowNode(element)
		}
	}
	return err
}

func (cont *ContinueMultiInstanceOperation) continueThroughMultiInstanceFlowNode(element delegate.FlowElement) error {
	cont.setLoopCounterVariable()

	var err error

	dataManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	if err = dataManager.RecordActivityStart(cont.Execution); err != nil {
		return err
	}

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	if err = executionEntityManager.RecordBusinessStatus(cont.Execution); err != nil {
		return err
	}

	behavior := element.GetBehavior()
	if behavior != nil {
		err = behavior.Execute(cont.Execution)
	}
	return err
}

func (cont *ContinueMultiInstanceOperation) setLoopCounterVariable() error {
	return cont.Execution.SetVariableLocal(behavior.CollectionElementIndexVariable, cont.LoopCounter)
}
