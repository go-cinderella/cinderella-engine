package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type StartEventParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (startEventParseHandler StartEventParseHandler) GetHandledType() string {
	return model.StartEvent{}.GetType()
}

func (startEventParseHandler StartEventParseHandler) ExecuteParse(bpmnParse *BpmnParse, flow delegate.BaseElement) {

}
