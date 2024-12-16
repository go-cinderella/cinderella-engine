package parse

import (
	. "github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/parse/factory"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	log "github.com/sirupsen/logrus"
)

type BpmnParse struct {
	Name                    string
	Byte                    []byte
	DeploymentEntity        DeploymentEntity
	ActivityBehaviorFactory factory.ActivityBehaviorFactory
	BpmnParserHandlers      BpmnParseHandlers
	BpmnModel               *BpmnModel
	currentFlowElement      delegate.FlowElement
	CurrentProcess          *Process
}

func (bpmnParse BpmnParse) SourceInputStream(byte []byte) BpmnParse {
	bpmnParse.setStreamSource(byte)
	return bpmnParse
}

func (bpmnParse BpmnParse) Deployment(deployment DeploymentEntity) BpmnParse {
	bpmnParse.DeploymentEntity = deployment
	return bpmnParse
}

func (bpmnParse BpmnParse) SourceName(name string) BpmnParse {
	bpmnParse.Name = name
	return bpmnParse
}
func (bpmnParse *BpmnParse) setStreamSource(byte []byte) {
	if byte == nil {
		log.Error("invalid: multiple sources ")
		panic("invalid: multiple sources ")
	}
	bpmnParse.Byte = byte
}
func (bpmnParse *BpmnParse) Execute() {
	xmlConverter := converter.BpmnXMLConverter{}
	bpmnModel := xmlConverter.ConvertToBpmnModel(bpmnParse.Byte)
	bpmnParse.BpmnModel = bpmnModel
	bpmnParse.applyParseHandlers()
}

func (bpmnParse *BpmnParse) applyParseHandlers() {
	for _, process := range bpmnParse.BpmnModel.GetMainProcess() {
		bpmnParse.CurrentProcess = process
		bpmnParse.BpmnParserHandlers.ParseElement(bpmnParse, process)
	}
}

func (bpmnParse *BpmnParse) SetCurrentFlowElement(currentFlowElement delegate.FlowElement) {
	bpmnParse.currentFlowElement = currentFlowElement
}

func (bpmnParse *BpmnParse) GetCurrentFlowElement() delegate.FlowElement {
	return bpmnParse.currentFlowElement
}

func (bpmnParse *BpmnParse) ProcessFlowElements(flowElements []delegate.FlowElement) {
	for _, element := range flowElements {
		bpmnParse.BpmnParserHandlers.ParseElement(bpmnParse, element)
	}
}
