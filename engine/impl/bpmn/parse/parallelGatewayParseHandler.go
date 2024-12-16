package parse

import (
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ParallelGatewayParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (parallelGatewayParseHandler ParallelGatewayParseHandler) GetHandledType() string {
	return ParallelGateway{}.GetType()
}

func (parallelGatewayParseHandler ParallelGatewayParseHandler) ExecuteParse(bpmnParse *BpmnParse, baseElement delegate.BaseElement) {
	parallelGateway := baseElement.(*ParallelGateway)
	parallelGateway.SetBehavior(bpmnParse.ActivityBehaviorFactory.CreateParallelGatewayActivityBehavior(*parallelGateway))
}
