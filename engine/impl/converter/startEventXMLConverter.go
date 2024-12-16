package converter

import (
	. "encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type StartEventXMLConverter struct {
	BpmnXMLConverter
}

func (start StartEventXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_EVENT_START
}
func (start StartEventXMLConverter) ConvertXMLToElement(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process) delegate.BaseElement {
	startEvent := StartEvent{FlowNode: FlowNode{BaseHandlerType: delegate.BaseHandlerType(StartEvent{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}
	decoder.DecodeElement(&startEvent, &token)
	activeProcess.InitialFlowElement = &startEvent
	return &startEvent
}
