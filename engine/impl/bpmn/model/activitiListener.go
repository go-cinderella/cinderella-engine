package model

import "strings"

type ActivitiListener struct {
	EventType string `xml:"event,attr"`
}

func (receiver ActivitiListener) String() string {
	var sb strings.Builder
	sb.WriteString("ActivitiListener")
	sb.WriteString("{")
	sb.WriteString("EventType: ")
	sb.WriteString(receiver.EventType)
	sb.WriteString("}")
	return sb.String()
}
