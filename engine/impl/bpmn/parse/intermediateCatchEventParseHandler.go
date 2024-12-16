package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type IntermediateCatchEventParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (eventParseHandler IntermediateCatchEventParseHandler) GetHandledType() string {
	return model.IntermediateCatchEvent{}.GetType()
}

func (eventParseHandler IntermediateCatchEventParseHandler) ExecuteParse(bpmnParse *BpmnParse, baseElement delegate.BaseElement) {
	intermediateCatchEvent := baseElement.(*model.IntermediateCatchEvent)
	var eventDefinition delegate.BaseElement
	if intermediateCatchEvent.ConditionalEventDefinition != nil {
		eventDefinition = intermediateCatchEvent.ConditionalEventDefinition
	}

	bpmnParse.BpmnParserHandlers.ParseElement(bpmnParse, eventDefinition)

}
