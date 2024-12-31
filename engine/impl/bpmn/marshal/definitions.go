package marshal

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
)

type Definitions struct {
	XMLName         xml.Name  `xml:"bpmn:definitions"`
	Text            string    `xml:",chardata"`
	Xsi             string    `xml:"xmlns:xsi,attr"`
	Bpmn            string    `xml:"xmlns:bpmn,attr"`
	Bpmndi          string    `xml:"xmlns:bpmndi,attr"`
	Dc              string    `xml:"xmlns:dc,attr"`
	Di              string    `xml:"xmlns:di,attr"`
	Flowable        string    `xml:"xmlns:flowable,attr"`
	ID              string    `xml:"id,attr"`
	TargetNamespace string    `xml:"targetNamespace,attr"`
	Process         []Process `xml:"process"`
}

func NewDefinitions(processId string) Definitions {
	return Definitions{
		XMLName: xml.Name{
			Space: constant.BPMN2_NAMESPACE,
			Local: constant.ELEMENT_DEFINITIONS,
		},
		Xsi:             constant.XSI_NAMESPACE,
		Bpmn:            constant.BPMN2_NAMESPACE,
		Bpmndi:          constant.BPMNDI_NAMESPACE,
		Dc:              constant.OMGDC_NAMESPACE,
		Di:              constant.OMGDI_NAMESPACE,
		Flowable:        "http://flowable.org/bpmn",
		ID:              "Definitions_" + processId,
		TargetNamespace: "http://bpmn.io/schema/bpmn",
	}
}
