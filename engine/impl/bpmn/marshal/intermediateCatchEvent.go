package marshal

import (
	"encoding/xml"
)

// IntermediateCatchEvent 中间抛出事件
type IntermediateCatchEvent struct {
	FlowNode
	XMLName                    xml.Name                    `xml:"bpmn:intermediateCatchEvent"`
	ConditionalEventDefinition *ConditionalEventDefinition `xml:"bpmn:conditionalEventDefinition"`
	FormButtonEventDefinition  *FormButtonEventDefinition  `xml:"bpmn:extensionElements>flowable:formButtonEventDefinition"`
}
