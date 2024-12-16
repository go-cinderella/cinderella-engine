package model

import "strings"

type ExtensionElement struct {
	TaskListener []ActivitiListener `xml:"taskListener"`
}

func (receiver ExtensionElement) String() string {
	var sb strings.Builder
	sb.WriteString("ExtensionElement")
	sb.WriteString("{")
	for i, element := range receiver.TaskListener {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(element.String())
	}
	sb.WriteString("}")
	return sb.String()
}
