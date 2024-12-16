package converter

import (
	. "encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ParallelGatewayXMLConverter struct {
	BpmnXMLConverter
}

func (parallelGateway ParallelGatewayXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_GATEWAY_PARALLEL
}
func (parallelGateway ParallelGatewayXMLConverter) ConvertXMLToElement(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process) delegate.BaseElement {
	parallel := ParallelGateway{Gateway{FlowNode: FlowNode{BaseHandlerType: delegate.BaseHandlerType(ParallelGateway{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}}
	decoder.DecodeElement(&parallel, &token)
	return &parallel
}
