package converter

import (
	. "encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type EndEventXMLConverter struct {
	BpmnXMLConverter
}

func (endEvent EndEventXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_EVENT_END
}
func (endEvent EndEventXMLConverter) ConvertXMLToElement(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process) delegate.BaseElement {
	end := EndEvent{FlowNode: FlowNode{BaseHandlerType: delegate.BaseHandlerType(EndEvent{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}
	decoder.DecodeElement(&end, &token)
	return &end
}
