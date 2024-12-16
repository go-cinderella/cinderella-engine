package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ParseHandler interface {
	GetHandledType() string

	ExecuteParse(bpmnParse *BpmnParse, element delegate.BaseElement)
}
