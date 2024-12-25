package model

import "strings"

type FieldExtension struct {
	FieldName   string `xml:"name,attr"`
	StringValue string `xml:"string"`
	Expression  string `xml:"expression"`
}

func (receiver FieldExtension) String() string {
	var sb strings.Builder
	sb.WriteString("FieldExtension")
	sb.WriteString("{")

	sb.WriteString("FieldName: ")
	sb.WriteString(receiver.FieldName)
	sb.WriteString(", ")

	sb.WriteString("StringValue: ")
	sb.WriteString(receiver.StringValue)
	sb.WriteString(", ")

	sb.WriteString("Expression: ")
	sb.WriteString(receiver.Expression)

	sb.WriteString("}")
	return sb.String()
}
