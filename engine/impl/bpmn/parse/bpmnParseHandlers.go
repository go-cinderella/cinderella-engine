package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type BpmnParseHandlers struct {
	ParseHandlers map[string][]BpmnParseHandler
}

func NewBpmnParseHandlers() BpmnParseHandlers {
	return BpmnParseHandlers{ParseHandlers: make(map[string][]BpmnParseHandler, 0)}
}
func (bpmnParseHandlers BpmnParseHandlers) AddHandlers(handlers []BpmnParseHandler) {
	for _, handler := range handlers {
		bpmnParseHandlers.AddHandler(handler)
	}
}

func (bpmnParseHandlers BpmnParseHandlers) AddHandler(bpmnParseHandler BpmnParseHandler) {
	handledTypes := bpmnParseHandler.GetHandledTypes()
	for _, handledType := range handledTypes {
		_, ok := bpmnParseHandlers.ParseHandlers[handledType]
		if !ok {
			parseHandlers := make([]BpmnParseHandler, 0)
			bpmnParseHandlers.ParseHandlers[handledType] = parseHandlers
		}
		bpmnParseHandlers.ParseHandlers[handledType] = append(bpmnParseHandlers.ParseHandlers[handledType], bpmnParseHandler)

	}
}
func (bpmnParseHandlers BpmnParseHandlers) ParseElement(bpmnParse *BpmnParse, element delegate.BaseElement) {
	var handlerType string
	flowElement, ok := element.(delegate.FlowElement)
	if ok {
		bpmnParse.SetCurrentFlowElement(flowElement)
		handlerType = flowElement.GetHandlerType()
	} else {
		handlerType = element.GetHandlerType()
	}

	handlers := bpmnParseHandlers.ParseHandlers[handlerType]
	for _, handler := range handlers {
		handler.Parse(bpmnParse, element)
	}
}
