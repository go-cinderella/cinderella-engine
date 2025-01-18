package cmd

import (
	"context"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	bpmn_model "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/converter"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/stringutils"
	"strings"
)

var _ engine.Command = (*GetNextFormButtonsCmd)(nil)

type GetNextFormButtonsCmd struct {
	Ctx                 context.Context
	ProcessDefinitionId string
	TaskDefinitionKey   string
	Transactional       bool
}

func (g GetNextFormButtonsCmd) IsTransactional() bool {
	return g.Transactional
}

type UserTaskFormButton struct {
	ActionValue     *string
	ActionName      *string
	CandidateGroups []string
	OpenConfirm     *bool
	OpenForm        *bool
	FormKey         *string
}

func (g GetNextFormButtonsCmd) Execute(commandContext engine.Context) (interface{}, error) {
	result, err := GetDeploymentResourceCmd{ProcessDefinitionId: g.ProcessDefinitionId, Ctx: g.Ctx}.Execute(commandContext)
	if err != nil {
		return nil, err
	}

	resourceEntity := result.(entitymanager.ResourceEntity)
	bpmnXMLConverter := converter.BpmnXMLConverter{}
	bpmnModel := bpmnXMLConverter.ConvertToBpmnModel(resourceEntity.GetBytes())
	process := bpmnModel.GetProcess()

	userTask := process.GetFlowElement(g.TaskDefinitionKey)
	outgoings := userTask.GetOutgoing()
	targetFlowElements := lo.Map[delegate.FlowElement, delegate.FlowElement](outgoings, func(item delegate.FlowElement, index int) delegate.FlowElement {
		return item.GetTargetFlowElement()
	})
	targetFlowElements = lo.Filter(targetFlowElements, func(item delegate.FlowElement, index int) bool {
		intermediateCatchEvent, ok := item.(*bpmn_model.IntermediateCatchEvent)
		if !ok {
			return false
		}
		return intermediateCatchEvent.FormButtonEventDefinition != nil
	})
	conditionalEvents := lo.Map[delegate.FlowElement, *bpmn_model.IntermediateCatchEvent](targetFlowElements, func(item delegate.FlowElement, index int) *bpmn_model.IntermediateCatchEvent {
		return item.(*bpmn_model.IntermediateCatchEvent)
	})
	mainForm := process.FormKey
	return lo.Map[*bpmn_model.IntermediateCatchEvent, *UserTaskFormButton](conditionalEvents, func(item *bpmn_model.IntermediateCatchEvent, index int) *UserTaskFormButton {
		formKey := item.FormButtonEventDefinition.FormKey
		if stringutils.IsEmpty(formKey) {
			formKey = mainForm
		}
		return &UserTaskFormButton{
			ActionValue:     &item.Id,
			ActionName:      &item.Name,
			CandidateGroups: lo.Ternary(stringutils.IsNotEmpty(item.FormButtonEventDefinition.CandidateGroups), strings.Split(item.FormButtonEventDefinition.CandidateGroups, ","), []string{}),
			OpenConfirm:     &item.FormButtonEventDefinition.OpenConfirm,
			OpenForm:        &item.FormButtonEventDefinition.OpenForm,
			FormKey:         &formKey,
		}
	}), nil
}

func (g GetNextFormButtonsCmd) Context() context.Context {
	return g.Ctx
}
