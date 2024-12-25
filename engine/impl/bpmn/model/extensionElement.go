package model

import "strings"

type ExtensionElement struct {
	TaskListener    []ActivitiListener `xml:"taskListener"`
	FieldExtensions []FieldExtension   `xml:"field"`
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
	for i, element := range receiver.FieldExtensions {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(element.String())
	}
	sb.WriteString("}")
	return sb.String()
}

func (receiver ExtensionElement) GetFieldByName(fieldName string) FieldExtension {
	for _, element := range receiver.FieldExtensions {
		if element.FieldName == fieldName {
			return element
		}
	}
	return FieldExtension{}
}
