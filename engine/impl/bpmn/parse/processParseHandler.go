package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ProcessParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (processParseHandler ProcessParseHandler) GetHandledType() string {
	return model.Process{}.GetType()
}

func (processParseHandler ProcessParseHandler) ExecuteParse(bpmnParse *BpmnParse, flow delegate.BaseElement) {
	process := flow.(*model.Process)
	bpmnParse.ProcessFlowElements(process.FlowElementList)
}
