package model

type StartEvent struct {
	FlowNode
	Initiator string `xml:"initiator,attr"`
	FormKey   string `xml:"formKey,attr"`
}

func (start StartEvent) GetType() string {
	return "startEvent"
}
