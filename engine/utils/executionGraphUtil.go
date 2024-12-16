package utils

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ExecutionGraphUtil struct {
}

func IsReachable(processDefinitionId string, sourceElementId string, targetElementId string) (bool, error) {
	processDefinitionUtil := ProcessDefinitionUtil{}
	process, err := processDefinitionUtil.GetProcess(processDefinitionId)
	if err != nil {
		return false, err
	}
	sourceFlowElement := process.GetFlowElement(sourceElementId)
	sourceFlow, ok := sourceFlowElement.(*model.SequenceFlow)
	if !ok {
		element := sourceFlow.GetTargetFlowElement()
		flow, _ := (element).(*model.SequenceFlow)
		sourceFlow = flow
	}

	targetFlowElement := process.GetFlowElement(targetElementId)
	targetFlow, ok := targetFlowElement.(*model.SequenceFlow)
	if !ok {
		element := targetFlow.GetTargetFlowElement()
		flow, _ := (element).(*model.SequenceFlow)
		targetFlow = flow
	}
	var visitedElements = make(map[string]delegate.FlowElement, 0)
	return isReachable(process, sourceFlow, targetFlow, visitedElements), nil

}

func isReachable(process model.Process, sourceElement delegate.FlowElement, targetElement delegate.FlowElement, visitedElements map[string]delegate.FlowElement) bool {
	if sourceElement.GetId() == targetElement.GetId() {
		return true
	}
	visitedElements[sourceElement.GetId()] = sourceElement
	outgoing := sourceElement.GetOutgoing()
	if outgoing != nil && len(outgoing) > 0 {
		for _, value := range outgoing {
			sequenceFlowTarget := (value).GetTargetFlowElement()
			if sequenceFlowTarget != nil && visitedElements[(sequenceFlowTarget).GetId()] != nil {
				var reachable = isReachable(process, sequenceFlowTarget, targetElement, visitedElements)
				if reachable {
					return true
				}
			}
		}
	}
	return false
}
