package model

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"reflect"
)

var _ delegate.BaseHandlerType = (*ServiceTask)(nil)
var _ delegate.BaseElement = (*ServiceTask)(nil)
var _ delegate.FlowElement = (*ServiceTask)(nil)

type ServiceTask struct {
	Task
	XMLName           xml.Name          `xml:"serviceTask"`
	TaskType          string            `xml:"type,attr"`
	ExtensionElements *ExtensionElement `xml:"extensionElements"`
}

func (serviceTask ServiceTask) GetType() string {
	return constant.ELEMENT_TASK_SERVICE
}

func (serviceTask *ServiceTask) Equal(otherServiceTask interface{}) bool {
	if otherServiceTask == nil {
		return serviceTask == nil
	}

	other, ok := otherServiceTask.(*ServiceTask)
	if !ok {
		that2, ok := otherServiceTask.(ServiceTask)
		if ok {
			other = &that2
		} else {
			return false
		}
	}

	if other == nil {
		return serviceTask == nil
	} else if serviceTask == nil {
		return false
	}

	if serviceTask.DefaultBaseElement != other.DefaultBaseElement {
		return false
	}

	if serviceTask.TaskType != other.TaskType {
		return false
	}

	if !reflect.DeepEqual(serviceTask.ExtensionElements, other.ExtensionElements) {
		return false
	}

	return true
}

func (serviceTask ServiceTask) Clone() delegate.FlowElement {
	serviceTaskCopy := serviceTask
	return &serviceTaskCopy
}
