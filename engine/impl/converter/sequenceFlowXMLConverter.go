package converter

import (
	. "encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type SequenceFlowXMLConverter struct {
	BpmnXMLConverter
}

func (sequence SequenceFlowXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_SEQUENCE_FLOW
}

func (sequence SequenceFlowXMLConverter) ConvertXMLToElement(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process) delegate.BaseElement {
	sequenceFlow := SequenceFlow{FlowNode: FlowNode{BaseHandlerType: delegate.BaseHandlerType(SequenceFlow{})}}
	decoder.DecodeElement(&sequenceFlow, &token)
	return &sequenceFlow
}
