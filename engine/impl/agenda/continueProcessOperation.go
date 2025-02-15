package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ContinueProcessOperation struct {
	AbstractOperation
}

func (cont *ContinueProcessOperation) Run() (err error) {
	element := cont.Execution.GetCurrentFlowElement()
	if element != nil {
		if element.GetOutgoing() != nil {
			err = cont.continueThroughFlowNode(element)
		} else {
			err = cont.continueThroughSequenceFlow(element)
		}
	}
	return err
}

func (cont *ContinueProcessOperation) continueThroughSequenceFlow(sequenceFlow delegate.FlowElement) error {
	execution := cont.Execution
	executionEntityManager := entitymanager.GetExecutionEntityManager()

	dataManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	if err := dataManager.RecordSequenceFlow(execution); err != nil {
		return err
	}

	if err := executionEntityManager.DeleteRelatedDataForExecution(execution.GetExecutionId(), nil); err != nil {
		return err
	}

	// 这里不能用defer函数删旧数据，因为加了外键约束
	if err := executionEntityManager.DeleteExecution(execution.GetExecutionId()); err != nil {
		return err
	}

	executionEntity := entitymanager.CreateExecution(execution)
	executionEntity.SetCurrentFlowElement(sequenceFlow.GetTargetFlowElement())

	if err := executionEntityManager.CreateExecution(&executionEntity); err != nil {
		return err
	}

	cont.GetAgenda().PlanContinueProcessOperation(&executionEntity)

	return nil
}

func (cont *ContinueProcessOperation) hasMultiInstanceRootExecution(execution delegate.DelegateExecution, element delegate.FlowElement) (bool, error) {
	currentExecution, err := execution.GetParent()
	if err != nil {
		return false, err
	}

	for currentExecution != nil {
		if currentExecution.IsMultiInstanceRoot() && element.GetId() == currentExecution.GetCurrentActivityId() {
			return true, nil
		}

		currentExecution, err = currentExecution.GetParent()
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

func (cont *ContinueProcessOperation) createMultiInstanceRootExecution(execution delegate.DelegateExecution) (delegate.DelegateExecution, error) {
	parentExecution, err := execution.GetParent()
	if err != nil {
		return nil, err
	}

	flowElement := execution.GetCurrentFlowElement()

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	if err = executionEntityManager.DeleteRelatedDataForExecution(execution.GetExecutionId(), nil); err != nil {
		return nil, err
	}
	if err = executionEntityManager.DeleteExecution(execution.GetExecutionId()); err != nil {
		return nil, err
	}

	multiInstanceRootExecution := entitymanager.CreateChildExecution(parentExecution)
	multiInstanceRootExecution.SetCurrentFlowElement(flowElement)
	multiInstanceRootExecution.SetMultiInstanceRoot(true)

	if err = executionEntityManager.CreateExecution(&multiInstanceRootExecution); err != nil {
		return nil, err
	}

	return &multiInstanceRootExecution, nil
}

func (cont *ContinueProcessOperation) executeActivityBehavior(activityBehavior delegate.ActivityBehavior, element delegate.FlowElement) error {
	return activityBehavior.Execute(cont.Execution)
}

func (cont *ContinueProcessOperation) executeMultiInstanceSynchronous(element delegate.FlowElement) error {
	ok, err := cont.hasMultiInstanceRootExecution(cont.Execution, element)
	if err != nil {
		return err
	}

	if !ok {
		cont.Execution, err = cont.createMultiInstanceRootExecution(cont.Execution)
		if err != nil {
			return err
		}
	}

	behavior := element.GetBehavior()
	if behavior == nil {
		panic("Expected an activity behavior in flow node " + element.GetId())
	}
	return cont.executeActivityBehavior(behavior, element)
}

func (cont *ContinueProcessOperation) executeSynchronous(element delegate.FlowElement) error {
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
		err = cont.executeActivityBehavior(behavior, element)
	} else {
		cont.GetAgenda().PlanTakeOutgoingSequenceFlowsOperation(cont.Execution, true)
	}

	return err
}

func (cont *ContinueProcessOperation) continueThroughFlowNode(element delegate.FlowElement) error {
	getter, ok := element.(model.LoopCharacteristicsGetter)
	if ok && getter.HasMultiInstanceLoopCharacteristics() {
		return cont.executeMultiInstanceSynchronous(element)
	}

	return cont.executeSynchronous(element)
}
