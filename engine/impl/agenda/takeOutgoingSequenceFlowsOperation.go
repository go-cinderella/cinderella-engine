package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
)

type TakeOutgoingSequenceFlowsOperation struct {
	AbstractOperation
	EvaluateConditions bool
}

func (take TakeOutgoingSequenceFlowsOperation) Run() (err error) {
	currentFlowElement := take.Execution.GetCurrentFlowElement()
	if currentFlowElement.GetOutgoing() != nil {
		err = take.handleFlowNode()
	} else {
		take.handleSequenceFlow()
	}
	return err
}

// 处理节点
func (take TakeOutgoingSequenceFlowsOperation) handleFlowNode() (err error) {
	execution := take.Execution
	currentFlowElement := execution.GetCurrentFlowElement()
	executionEntityManager := entitymanager.GetExecutionEntityManager()

	if err = take.handleActivityEndByExecutionId(currentFlowElement); err != nil {
		return
	}

	if err = executionEntityManager.DeleteRelatedDataForExecution(execution.GetExecutionId(), nil); err != nil {
		return
	}

	if err = executionEntityManager.DeleteExecution(execution.GetExecutionId()); err != nil {
		return
	}

	var defaultSequenceFlowId string

	gateway, ok := currentFlowElement.(*model.Gateway)
	if ok {
		defaultSequenceFlowId = gateway.DefaultFlow
	}

	var outgoingSequenceFlows []delegate.FlowElement

	flowElements := currentFlowElement.GetOutgoing()
	for _, flowElement := range flowElements {
		sequenceFlow := (flowElement).(*model.SequenceFlow)
		if !take.EvaluateConditions || utils.HasTrueCondition(*sequenceFlow, execution) {
			outgoingSequenceFlows = append(outgoingSequenceFlows, flowElement)
		}
	}

	if len(outgoingSequenceFlows) == 0 {
		for _, flowElement := range flowElements {
			if defaultSequenceFlowId == flowElement.GetId() {
				outgoingSequenceFlows = append(outgoingSequenceFlows, flowElement)
			}
		}
	}

	if len(outgoingSequenceFlows) == 0 {
		if len(flowElements) == 0 {
			take.GetAgenda().PlanEndExecutionOperation(execution)
		}
		return nil
	}

	for _, outgoing := range outgoingSequenceFlows {

		executionEntity := entitymanager.CreateExecution(execution)
		executionEntity.SetCurrentFlowElement(outgoing)

		if err = executionEntityManager.CreateExecution(&executionEntity); err != nil {
			return err
		}

		take.GetAgenda().PlanContinueProcessOperation(&executionEntity)
	}

	return nil
}

// 处理连线
func (take TakeOutgoingSequenceFlowsOperation) handleSequenceFlow() {
	take.GetAgenda().PlanContinueProcessOperation(take.Execution)
}

func (take TakeOutgoingSequenceFlowsOperation) handleActivityEndByExecutionId(element delegate.FlowElement) (err error) {
	dataManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	err = dataManager.RecordActEndByExecutionId(element, take.Execution, nil)
	return err
}
