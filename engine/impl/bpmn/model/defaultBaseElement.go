package model

import "strings"

type DefaultBaseElement struct {
	Id   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

func (d DefaultBaseElement) GetId() string {
	return d.Id
}

func (d DefaultBaseElement) GetName() string {
	return d.Name
}

func (d DefaultBaseElement) String() string {
	var sb strings.Builder
	sb.WriteString("DefaultBaseElement")
	sb.WriteString("{")
	sb.WriteString("Id: ")
	sb.WriteString(d.GetId())
	sb.WriteString(", ")
	sb.WriteString("Name: ")
	sb.WriteString(d.GetName())
	sb.WriteString("}")
	return sb.String()
}
