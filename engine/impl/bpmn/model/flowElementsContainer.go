package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type FlowElementsContainer interface {
	AddFlowElement(element delegate.FlowElement)
}
