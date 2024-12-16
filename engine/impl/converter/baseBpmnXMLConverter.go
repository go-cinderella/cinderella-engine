package converter

import (
	. "encoding/xml"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type BaseBpmnXMLConverter interface {
	GetXMLElementName() string

	convertToBpmnModel(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process)

	ConvertXMLToElement(decoder *Decoder, token StartElement, model *BpmnModel, activeProcess *Process) delegate.BaseElement
}
