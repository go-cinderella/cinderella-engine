package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"time"
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

	// 这里不能用defer函数删旧数据，因为加了外键约束
	if err := executionEntityManager.DeleteExecution(execution.GetExecutionId()); err != nil {
		return err
	}

	executionEntity := entitymanager.ExecutionEntity{
		ProcessInstanceId:   execution.GetProcessInstanceId(),
		ProcessDefinitionId: execution.GetProcessDefinitionId(),
		StartTime:           time.Now().UTC(),
	}

	flowElement := sequenceFlow.GetTargetFlowElement()
	executionEntity.SetCurrentFlowElement(flowElement)

	if err := executionEntityManager.CreateExecution(&executionEntity); err != nil {
		return err
	}

	cont.GetAgenda().PlanContinueProcessOperation(&executionEntity)

	return nil
}

func (cont *ContinueProcessOperation) continueThroughFlowNode(element delegate.FlowElement) error {
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
	} else {
		cont.GetAgenda().PlanTakeOutgoingSequenceFlowsOperation(cont.Execution, true)
	}
	return err
}
