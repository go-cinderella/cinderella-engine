package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ConditionalEventDefinitionParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (conditionalEventDefinitionParseHandler ConditionalEventDefinitionParseHandler) GetHandledType() string {
	return model.ConditionalEventDefinition{}.GetHandlerType()
}

func (conditionalEventDefinitionParseHandler ConditionalEventDefinitionParseHandler) ExecuteParse(bpmnParse *BpmnParse, baseElement delegate.BaseElement) {
	currentFlowElement := bpmnParse.GetCurrentFlowElement()
	intermediateCatchEvent, ok := currentFlowElement.(*model.IntermediateCatchEvent)
	if !ok {
		return
	}

	conditionalEventDefinition := baseElement.(*model.ConditionalEventDefinition)
	intermediateCatchEvent.SetBehavior(bpmnParse.ActivityBehaviorFactory.CreateIntermediateCatchConditionalEventActivityBehavior(*conditionalEventDefinition))
}
