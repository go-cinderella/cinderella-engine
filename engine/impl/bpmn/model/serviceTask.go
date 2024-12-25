package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ delegate.BaseHandlerType = (*ServiceTask)(nil)
var _ delegate.BaseElement = (*ServiceTask)(nil)
var _ delegate.FlowElement = (*ServiceTask)(nil)

type ServiceTask struct {
	FlowNode
	TaskType          string            `xml:"type,attr"`
	ExtensionElements *ExtensionElement `xml:"extensionElements"`
}

func (serviceTask ServiceTask) GetType() string {
	return constant.ELEMENT_TASK_SERVICE
}
