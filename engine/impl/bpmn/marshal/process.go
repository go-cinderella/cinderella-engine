package marshal

import (
	"encoding/xml"
)

type Process struct {
	FlowNode
	XMLName                 xml.Name                 `xml:"bpmn:process"`
	IsExecutable            string                   `xml:"isExecutable,attr"`
	StartEvents             []StartEvent             `xml:"startEvent"`
	EndEvents               []EndEvent               `xml:"endEvent"`
	UserTasks               []UserTask               `xml:"userTask"`
	SequenceFlows           []SequenceFlow           `xml:"sequenceFlow"`
	IntermediateCatchEvents []IntermediateCatchEvent `xml:"intermediateCatchEvent"`
}