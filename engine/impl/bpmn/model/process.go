package model

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type Process struct {
	FlowNode
	XMLName                 xml.Name                 `xml:"process"`
	IsExecutable            string                   `xml:"isExecutable,attr"`
	FormKey                 string                   `xml:"formKey,attr"` // main form key/default form key
	StartEvents             []StartEvent             `xml:"startEvent"`
	EndEvents               []EndEvent               `xml:"endEvent"`
	UserTasks               []UserTask               `xml:"userTask"`
	ServiceTasks            []ServiceTask            `xml:"serviceTask"`
	SequenceFlows           []SequenceFlow           `xml:"sequenceFlow"`
	ExclusiveGateways       []ExclusiveGateway       `xml:"exclusiveGateway"`
	InclusiveGateways       []InclusiveGateway       `xml:"inclusiveGateway"`
	ParallelGateways        []ParallelGateway        `xml:"parallelGateway"`
	BoundaryEvents          []BoundaryEvent          `xml:"boundaryEvent"`
	IntermediateCatchEvents []IntermediateCatchEvent `xml:"intermediateCatchEvent"`
	SubProcesses            []SubProcess             `xml:"subProcess"`
	FlowElementList         []delegate.FlowElement
	InitialFlowElement      delegate.FlowElement
	FlowElementMap          map[string]delegate.FlowElement

	// TODO
	//IntermediateTrowEvent        []IntermediateThrowEvent `xml:"intermediateThrowEvent"`
	//EventBasedGateway            []EventBasedGateway      `xml:"eventBasedGateway"`
}

func (process Process) GetFlowElement(flowElementId string) delegate.FlowElement {
	return process.FlowElementMap[flowElementId]
}

func (process *Process) AddFlowElement(element delegate.FlowElement) {
	process.FlowElementList = append(process.FlowElementList, element)
	process.FlowElementMap[element.GetId()] = element
}

func (process Process) GetType() string {
	return "process"
}
