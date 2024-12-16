package converter

import (
	. "encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ExclusiveGatewayXMLConverter struct {
	BpmnXMLConverter
}

func (exclusiveGateway ExclusiveGatewayXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_GATEWAY_EXCLUSIVE
}
func (exclusiveGateway ExclusiveGatewayXMLConverter) ConvertXMLToElement(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process) delegate.BaseElement {
	exclusive := ExclusiveGateway{Gateway{FlowNode: FlowNode{BaseHandlerType: delegate.BaseHandlerType(ExclusiveGateway{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}}
	decoder.DecodeElement(&exclusive, &token)
	return &exclusive
}
