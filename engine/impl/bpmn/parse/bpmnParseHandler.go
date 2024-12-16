package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type BpmnParseHandler interface {
	GetHandledTypes() []string

	Parse(bpmnParse *BpmnParse, flow delegate.BaseElement)
}
