package converter

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ServiceTaskXMLConverter struct {
	BpmnXMLConverter
}

func (user ServiceTaskXMLConverter) GetXMLElementName() string {
	return constant.ELEMENT_TASK_SERVICE
}

func (user ServiceTaskXMLConverter) ConvertXMLToElement(decoder *xml.Decoder, token xml.StartElement, bpmnModel *model.BpmnModel, activeProcess *model.Process) delegate.BaseElement {
	serviceTask := &model.ServiceTask{FlowNode: model.FlowNode{BaseHandlerType: delegate.BaseHandlerType(model.ServiceTask{}), IncomingFlow: make([]delegate.FlowElement, 0), OutgoingFlow: make([]delegate.FlowElement, 0)}}
	decoder.DecodeElement(serviceTask, &token)
	return serviceTask
}
