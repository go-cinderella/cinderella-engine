package marshal

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type SequenceFlow struct {
	FlowNode
	XMLName             xml.Name `xml:"bpmn:sequenceFlow"`
	SourceRef           string   `xml:"sourceRef,attr"`
	TargetRef           string   `xml:"targetRef,attr"`
	ConditionExpression string   `xml:"bpmn:conditionExpression,omitempty"`
}

func NewSequenceFlow(sourceRef, targetRef string) SequenceFlow {
	return SequenceFlow{
		FlowNode: FlowNode{
			DefaultBaseElement: DefaultBaseElement{
				Id:   "Flow_" + gonanoid.MustGenerate(constant.Alphabet, constant.ColumnRandomCodeLength),
				Name: "",
			},
		},
		XMLName: xml.Name{
			Space: constant.BPMN2_NAMESPACE,
			Local: constant.ELEMENT_SEQUENCE_FLOW,
		},
		SourceRef:           sourceRef,
		TargetRef:           targetRef,
		ConditionExpression: "",
	}
}
