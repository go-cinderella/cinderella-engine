package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type AbstractActivityBpmnParseHandler struct {
	AbstractBpmnParseHandler
}

func (abstractBpmnParse AbstractActivityBpmnParseHandler) Parse(bpmnParse *BpmnParse, element delegate.BaseElement) {
	abstractBpmnParse.AbstractBpmnParseHandler.Parse(bpmnParse, element)
}
