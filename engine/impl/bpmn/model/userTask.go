package model

import (
	"encoding/xml"
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/spf13/cast"
)

var _ delegate.BaseHandlerType = (*UserTask)(nil)
var _ delegate.BaseElement = (*UserTask)(nil)
var _ delegate.FlowElement = (*UserTask)(nil)

type UserTask struct {
	Task
	XMLName         xml.Name `xml:"userTask"`
	Assignee        *string  `xml:"assignee,attr"`
	FormKey         *string  `xml:"formKey,attr"`
	CandidateUsers  *string  `xml:"candidateUsers,attr"`
	CandidateGroups *string  `xml:"candidateGroups,attr"`
	DueDate         *string  `xml:"dueDate,attr"`
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

func (user *UserTask) ActivityEqual(otherUser interface{}) bool {
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

	if !user.Activity.ActivityEqual(other.Activity) {
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

	return true
}

func (user UserTask) Clone() delegate.FlowElement {
	userCopy := user
	return &userCopy
}
