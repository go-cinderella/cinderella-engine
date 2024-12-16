package converter

import (
	. "encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type InclusiveGatewayXMLConverter struct {
	BpmnXMLConverter
}

func (inclusiveGateway InclusiveGatewayXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_GATEWAY_INCLUSIVE
}
func (inclusiveGateway InclusiveGatewayXMLConverter) ConvertXMLToElement(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process) delegate.BaseElement {
	inclusive := InclusiveGateway{Gateway{FlowNode: FlowNode{BaseHandlerType: delegate.BaseHandlerType(InclusiveGateway{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}}
	decoder.DecodeElement(&inclusive, &token)
	return &inclusive
}
