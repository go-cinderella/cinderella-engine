package marshal

import (
	"encoding/xml"
)

type UserTaskExtensionElement struct {
	AssigneeType            string `xml:"flowable:assigneeType"`
	StaticAssigneeVariables string `xml:"flowable:staticAssigneeVariables,omitempty"`
	IdmCandidateGroups      string `xml:"flowable:idmCandidateGroups,omitempty"`
	IdmCandidateUsers       string `xml:"flowable:idmCandidateUsers,omitempty"`
}

type UserTask struct {
	FlowNode
	XMLName           xml.Name                  `xml:"bpmn:userTask"`
	Assignee          *string                   `xml:"flowable:assignee,attr"`
	FormKey           *string                   `xml:"flowable:formKey,attr"`
	CandidateUsers    *string                   `xml:"flowable:candidateUsers,attr"`
	CandidateGroups   *string                   `xml:"flowable:candidateGroups,attr"`
	DueDate           *string                   `xml:"flowable:dueDate,attr"`
	ExtensionElements *UserTaskExtensionElement `xml:"bpmn:extensionElements"`
}
