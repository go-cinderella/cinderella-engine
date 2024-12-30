package model

import "encoding/xml"

type SequenceFlow struct {
	FlowNode
	XMLName             xml.Name `xml:"sequenceFlow"`
	SourceRef           string   `xml:"sourceRef,attr"`
	TargetRef           string   `xml:"targetRef,attr"`
	ConditionExpression string   `xml:"conditionExpression"`
}

func (sequenceFlow SequenceFlow) GetType() string {
	return "sequenceFlow"
}
