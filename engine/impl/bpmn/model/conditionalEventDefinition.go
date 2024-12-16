package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ delegate.BaseElement = (*ConditionalEventDefinition)(nil)

type ConditionalEventDefinition struct {
	DefaultBaseElement
	Condition string `xml:"condition"`
}

func (c ConditionalEventDefinition) GetHandlerType() string {
	return "conditionalEventDefinition"
}
