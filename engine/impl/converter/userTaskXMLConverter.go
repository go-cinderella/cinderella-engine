package converter

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type UserTaskXMLConverter struct {
	BpmnXMLConverter
}

func (user UserTaskXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_TASK_USER
}

func (user UserTaskXMLConverter) ConvertXMLToElement(decoder *xml.Decoder, token xml.StartElement, bpmnModel *model.BpmnModel, activeProcess *model.Process) delegate.BaseElement {
	userTask := model.UserTask{FlowNode: model.FlowNode{BaseHandlerType: delegate.BaseHandlerType(model.UserTask{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}
	decoder.DecodeElement(&userTask, &token)
	return &userTask
}
