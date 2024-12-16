package parser

import (
	. "encoding/xml"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ProcessParser struct {
}

func (ProcessParser ProcessParser) Parse(decoder *Decoder, token StartElement, model *BpmnModel) *Process {
	process := Process{FlowNode: FlowNode{BaseHandlerType: delegate.BaseHandlerType(Process{})}, FlowElementList: make([]delegate.FlowElement, 0), FlowElementMap: make(map[string]delegate.FlowElement, 0)}
	attrs := token.Attr
	tem := make(map[string]string, 0)
	for _, attr := range attrs {
		tem[attr.Name.Local] = attr.Value
	}
	process.Id = tem["id"]
	process.Name = tem["name"]
	process.IsExecutable = tem["isExecutable"]

	//这里不能用这个方法：decoder.DecodeElement(&process, &token)
	
	model.AddProcess(&process)
	return &process
}
