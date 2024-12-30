package model

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ delegate.BaseHandlerType = (*IntermediateCatchEvent)(nil)
var _ delegate.BaseElement = (*IntermediateCatchEvent)(nil)
var _ delegate.FlowElement = (*IntermediateCatchEvent)(nil)

// IntermediateCatchEvent 中间抛出事件
type IntermediateCatchEvent struct {
	FlowNode
	XMLName                    xml.Name                    `xml:"intermediateCatchEvent"`
	MessageEventDefinition     *MessageEventDefinition     `xml:"messageEventDefinition"`
	TimerEventDefinition       *TimerEventDefinition       `xml:"timerEventDefinition"`
	LinkEventDefinition        *LinkEventDefinition        `xml:"linkEventDefinition"`
	ConditionalEventDefinition *ConditionalEventDefinition `xml:"conditionalEventDefinition"`
	FormButtonEventDefinition  *FormButtonEventDefinition  `xml:"extensionElements>formButtonEventDefinition"`
}

func (event IntermediateCatchEvent) GetType() string {
	return constant.ELEMENT_EVENT_CATCH
}
