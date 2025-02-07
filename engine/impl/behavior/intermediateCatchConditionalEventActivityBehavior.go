package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/stringutils"
)

var _ delegate.ActivityBehavior = (*IntermediateCatchConditionalEventActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*IntermediateCatchConditionalEventActivityBehavior)(nil)

type IntermediateCatchConditionalEventActivityBehavior struct {
	IntermediateCatchEventActivityBehavior
	conditionalEventDefinition model.ConditionalEventDefinition
}

func NewIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition model.ConditionalEventDefinition) IntermediateCatchConditionalEventActivityBehavior {
	return IntermediateCatchConditionalEventActivityBehavior{conditionalEventDefinition: conditionalEventDefinition}
}

func (behavior *IntermediateCatchConditionalEventActivityBehavior) ConditionalEventDefinition() model.ConditionalEventDefinition {
	return behavior.conditionalEventDefinition
}

func (behavior *IntermediateCatchConditionalEventActivityBehavior) SetConditionalEventDefinition(conditionalEventDefinition model.ConditionalEventDefinition) {
	behavior.conditionalEventDefinition = conditionalEventDefinition
}

func (behavior IntermediateCatchConditionalEventActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	conditionExpression := behavior.conditionalEventDefinition.Condition
	if stringutils.IsEmpty(conditionExpression) {
		return nil
	}

	variables, err := execution.GetVariables()
	if err != nil {
		return err
	}
	if len(variables) == 0 {
		return nil
	}

	if utils.IsTrue(variables, conditionExpression) {
		return behavior.leaveIntermediateCatchEvent(execution)
	}

	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	if err := historicActivityInstanceEntityManager.DeleteByExecutionId(execution); err != nil {
		return err
	}

	executionEntityManager := entitymanager.GetExecutionEntityManager()
	if err := executionEntityManager.DeleteExecution(execution.GetExecutionId()); err != nil {
		return err
	}

	currentFlowElement := execution.GetCurrentFlowElement()
	incomings := currentFlowElement.GetIncoming()

	lo.ForEachWhile(incomings, func(item delegate.FlowElement, index int) (goon bool) {
		if err = historicActivityInstanceEntityManager.DeleteByProcessInstanceId(execution.GetProcessInstanceId(), item.GetId()); err != nil {
			return false
		}

		return true
	})

	if err != nil {
		return err
	}

	return nil
}

func (behavior IntermediateCatchConditionalEventActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	return nil
}
