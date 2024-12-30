package model

import "encoding/xml"

type EndEvent struct {
	FlowNode
	XMLName xml.Name `xml:"endEvent"`
}

func (endEvent EndEvent) GetType() string {
	return "endEvent"
}
