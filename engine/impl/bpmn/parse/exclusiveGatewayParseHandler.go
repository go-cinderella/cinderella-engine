package parse

import (
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ExclusiveGatewayParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (exclusiveGatewayParseHandler ExclusiveGatewayParseHandler) GetHandledType() string {
	return ExclusiveGateway{}.GetType()
}

func (exclusiveGatewayParseHandler ExclusiveGatewayParseHandler) ExecuteParse(bpmnParse *BpmnParse, baseElement delegate.BaseElement) {
	exclusiveGateway := baseElement.(*ExclusiveGateway)
	exclusiveGateway.SetBehavior(bpmnParse.ActivityBehaviorFactory.CreateExclusiveGatewayActivityBehavior(*exclusiveGateway))
}
