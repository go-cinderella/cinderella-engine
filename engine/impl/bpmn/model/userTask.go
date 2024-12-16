package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

var _ delegate.BaseHandlerType = (*UserTask)(nil)
var _ delegate.BaseElement = (*UserTask)(nil)
var _ delegate.FlowElement = (*UserTask)(nil)

type UserTask struct {
	FlowNode
	Assignee          *string                           `xml:"assignee,attr"`
	FormKey           *string                           `xml:"formKey,attr"`
	CandidateUsers    *string                           `xml:"candidateUsers,attr"`
	CandidateGroups   *string                           `xml:"candidateGroups,attr"`
	DueDate           *string                           `xml:"dueDate,attr"`
	MultiInstance     *MultiInstanceLoopCharacteristics `xml:"multiInstanceLoopCharacteristics"`
	ExtensionElements *ExtensionElement                 `xml:"extensionElements"`
}

func (user UserTask) GetType() string {
	return constant.ELEMENT_TASK_USER
}

func (user UserTask) GetCandidateUsers() *string {
	return user.CandidateUsers
}

func (user UserTask) GetCandidateGroups() *string {
	return user.CandidateGroups
}

func (user UserTask) GetAssignee() *string {
	return user.Assignee
}

func (user UserTask) String() string {
	var sb strings.Builder
	sb.WriteString("UserTask")
	sb.WriteString("{")
	sb.WriteString("FlowNode: ")
	sb.WriteString(user.FlowNode.String())

	if user.Assignee != nil {
		sb.WriteString(", ")
		sb.WriteString("Assignee: ")
		sb.WriteString(*user.Assignee)
	}
	if user.CandidateUsers != nil {
		sb.WriteString(", ")
		sb.WriteString("CandidateUsers: ")
		sb.WriteString(*user.CandidateUsers)
	}
	if user.CandidateGroups != nil {
		sb.WriteString(", ")
		sb.WriteString("CandidateGroups: ")
		sb.WriteString(*user.CandidateGroups)
	}
	if user.DueDate != nil {
		sb.WriteString(", ")
		sb.WriteString("DueDate: ")
		sb.WriteString(*user.DueDate)
	}
	if user.MultiInstance != nil {
		sb.WriteString(", ")
		sb.WriteString("MultiInstance: ")
		sb.WriteString(user.MultiInstance.String())
	}
	if user.ExtensionElements != nil {
		sb.WriteString(", ")
		sb.WriteString("ExtensionElements: ")
		sb.WriteString(user.ExtensionElements.String())
	}
	sb.WriteString("}")
	return sb.String()
}

func (user *UserTask) Equal(otherUser interface{}) bool {
	if otherUser == nil {
		return user == nil
	}

	other, ok := otherUser.(*UserTask)
	if !ok {
		that2, ok := otherUser.(UserTask)
		if ok {
			other = &that2
		} else {
			return false
		}
	}

	if other == nil {
		return user == nil
	} else if user == nil {
		return false
	}

	if user.DefaultBaseElement != other.DefaultBaseElement {
		return false
	}

	if cast.ToString(user.Assignee) != cast.ToString(other.Assignee) {
		return false
	}

	if cast.ToString(user.FormKey) != cast.ToString(other.FormKey) {
		return false
	}

	if cast.ToString(user.CandidateUsers) != cast.ToString(other.CandidateUsers) {
		return false
	}

	if cast.ToString(user.CandidateGroups) != cast.ToString(other.CandidateGroups) {
		return false
	}

	if cast.ToString(user.DueDate) != cast.ToString(other.DueDate) {
		return false
	}

	if !reflect.DeepEqual(user.MultiInstance, other.MultiInstance) {
		return false
	}

	if !reflect.DeepEqual(user.ExtensionElements, other.ExtensionElements) {
		return false
	}

	return true
}
