package marshal

import "encoding/xml"

type EndEvent struct {
	FlowNode
	XMLName xml.Name `xml:"bpmn:endEvent"`
}
