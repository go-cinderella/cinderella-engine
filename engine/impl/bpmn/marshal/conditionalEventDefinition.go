package marshal

import "encoding/xml"

type ConditionalEventDefinition struct {
	DefaultBaseElement
	XMLName   xml.Name `xml:"bpmn:conditionalEventDefinition"`
	Condition string   `xml:"bpmn:condition"`
}
