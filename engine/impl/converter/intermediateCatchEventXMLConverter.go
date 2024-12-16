package converter

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type IntermediateCatchEventXMLConverter struct {
	BpmnXMLConverter
}

func (converter IntermediateCatchEventXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_EVENT_CATCH
}

func (converter IntermediateCatchEventXMLConverter) ConvertXMLToElement(decoder *xml.Decoder, token xml.StartElement, bpmnModel *model.BpmnModel, activeProcess *model.Process) delegate.BaseElement {
	intermediateCatchEvent := model.IntermediateCatchEvent{FlowNode: model.FlowNode{BaseHandlerType: delegate.BaseHandlerType(model.IntermediateCatchEvent{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}
	decoder.DecodeElement(&intermediateCatchEvent, &token)
	return &intermediateCatchEvent
}
