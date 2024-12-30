package model

import "encoding/xml"

type StartEvent struct {
	FlowNode
	XMLName   xml.Name `xml:"startEvent"`
	Initiator string   `xml:"initiator,attr"`
	FormKey   string   `xml:"formKey,attr"`
}

func (start StartEvent) GetType() string {
	return "startEvent"
}
