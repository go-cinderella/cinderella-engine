package model

type SequenceFlow struct {
	FlowNode
	//Id                  string   `xml:"id,attr"`
	SourceRef           string `xml:"sourceRef,attr"`
	TargetRef           string `xml:"targetRef,attr"`
	ConditionExpression string `xml:"conditionExpression"`
}

func (sequenceFlow SequenceFlow) GetType() string {
	return "sequenceFlow"
}
