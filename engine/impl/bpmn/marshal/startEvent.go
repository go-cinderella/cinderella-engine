package marshal

import "encoding/xml"

type StartEvent struct {
	FlowNode
	XMLName   xml.Name `xml:"bpmn:startEvent"`
	Initiator string   `xml:"flowable:initiator,attr,omitempty"`
	FormKey   string   `xml:"flowable:formKey,attr,omitempty"`
}
