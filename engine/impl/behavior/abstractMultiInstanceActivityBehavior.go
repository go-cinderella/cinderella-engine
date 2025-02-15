package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"strings"
)

type MultiInstanceActivity interface {
	delegate.FlowElement
	model.LoopCharacteristicsGetter
}

const (
	NUMBER_OF_INSTANCES            string = "nrOfInstances"
	NUMBER_OF_ACTIVE_INSTANCES     string = "nrOfActiveInstances"
	NUMBER_OF_COMPLETED_INSTANCES  string = "nrOfCompletedInstances"
	DELETE_REASON_END              string = "MI_END"
	CollectionElementIndexVariable string = "loopCounter"
)

type MultiInstanceSupportBehavior interface {
	delegate.TriggerableActivityBehavior
	SetMultiInstanceActivityBehavior(multiInstanceActivityBehavior multiInstanceActivityBehavior)
}

// AbstractMultiInstanceActivityBehavior Implementation of the multi-instance functionality as described in the BPMN 2.0 spec.
type AbstractMultiInstanceActivityBehavior struct {
	flowNodeActivityBehavior

	Impl multiInstanceActivityBehavior
	// Instance members
	Activity              MultiInstanceActivity
	InnerActivityBehavior MultiInstanceSupportBehavior
}

func (f AbstractMultiInstanceActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	return f.InnerActivityBehavior.Trigger(execution)
}

func (f AbstractMultiInstanceActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	_, ok, err := f.getLocalLoopVariable(execution, CollectionElementIndexVariable)
	if err != nil {
		return err
	}
	if !ok {
		nrOfInstances, err := f.Impl.createInstances(execution)
		if err != nil {
			return err
		}
		if nrOfInstances == 0 {
			if err = f.cleanupMiRoot(execution); err != nil {
				return err
			}
		}
		return nil
	}

	return f.InnerActivityBehavior.Execute(execution)
}

func (f AbstractMultiInstanceActivityBehavior) leave(execution delegate.DelegateExecution) error {
	return f.cleanupMiRoot(execution)
}

func (f AbstractMultiInstanceActivityBehavior) getLocalLoopVariable(execution delegate.DelegateExecution, variableName string) (value int, ok bool, err error) {
	localVariables, err := execution.GetVariablesLocal()
	if err != nil {
		return 0, false, err
	}
	if lo.HasKey(localVariables, variableName) {
		return cast.ToInt(localVariables[variableName]), true, nil
	}

	if !execution.IsMultiInstanceRoot() {
		parentExecution, err := execution.GetParent()
		if err != nil {
			return 0, false, err
		}
		localVariables, err = parentExecution.GetVariablesLocal()
		if err != nil {
			return 0, false, err
		}
		if lo.HasKey(localVariables, variableName) {
			return cast.ToInt(localVariables[variableName]), true, nil
		}

		if !parentExecution.IsMultiInstanceRoot() {
			superExecution, err := parentExecution.GetParent()
			if err != nil {
				return 0, false, err
			}
			localVariables, err = superExecution.GetVariablesLocal()
			if err != nil {
				return 0, false, err
			}
			if lo.HasKey(localVariables, variableName) {
				return cast.ToInt(localVariables[variableName]), true, nil
			}
		}

	}

	return 0, false, nil
}

func (f AbstractMultiInstanceActivityBehavior) cleanupMiRoot(execution delegate.DelegateExecution) error {
	// Delete multi instance root and all child executions.
	// Create a fresh execution to continue

	multiInstanceRootExecution, err := f.getMultiInstanceRootExecution(execution)
	if err != nil {
		return err
	}

	executionEntityManager := entitymanager.GetExecutionEntityManager()

	var children []entitymanager.ExecutionEntity

	children, err = executionEntityManager.CollectChildren(multiInstanceRootExecution.GetExecutionId())
	if err != nil {
		return err
	}

	actions, err := executionEntityManager.GetTopKValueFromChildExecutions(children, "action", 1)
	if err != nil {
		return err
	}

	if len(actions) > 0 {
		variables := make(map[string]interface{})
		variables["action"] = actions[0]

		if err = execution.SetProcessVariables(variables); err != nil {
			return err
		}
	}

	if err = executionEntityManager.DeleteChildExecutions(children, lo.ToPtr(DELETE_REASON_END)); err != nil {
		return err
	}
	if err = executionEntityManager.DeleteRelatedDataForExecution(multiInstanceRootExecution.GetExecutionId(), lo.ToPtr(DELETE_REASON_END)); err != nil {
		return err
	}
	if err = executionEntityManager.DeleteExecution(multiInstanceRootExecution.GetExecutionId()); err != nil {
		return err
	}

	parentExecution, err := multiInstanceRootExecution.GetParent()
	if err != nil {
		return err
	}

	newExecution := entitymanager.CreateChildExecution(parentExecution)
	newExecution.SetCurrentFlowElement(execution.GetCurrentFlowElement())

	if err = executionEntityManager.CreateExecution(&newExecution); err != nil {
		return err
	}

	return f.flowNodeActivityBehavior.leave(&newExecution)
}

func (f AbstractMultiInstanceActivityBehavior) getMultiInstanceRootExecution(execution delegate.DelegateExecution) (delegate.DelegateExecution, error) {
	var multiInstanceRootExecution delegate.DelegateExecution
	var err error
	currentExecution := execution

	for currentExecution != nil && multiInstanceRootExecution == nil && stringutils.IsNotEmpty(currentExecution.GetParentId()) {
		if currentExecution.IsMultiInstanceRoot() {
			multiInstanceRootExecution = currentExecution
		} else {
			currentExecution, err = currentExecution.GetParent()
			if err != nil {
				return nil, err
			}
		}
	}

	return multiInstanceRootExecution, nil
}

func (f AbstractMultiInstanceActivityBehavior) resolveNrOfInstances(execution delegate.DelegateExecution) (int, error) {
	instances, err := f.resolveAndValidateCollection(execution)
	if err != nil {
		return 0, err
	}

	return len(instances), nil
}

func (f AbstractMultiInstanceActivityBehavior) resolveAndValidateCollection(execution delegate.DelegateExecution) ([]string, error) {
	loopCharacteristics := f.Activity.GetLoopCharacteristics()
	if loopCharacteristics == nil {
		return nil, nil
	}

	collection := loopCharacteristics.Collection
	if stringutils.IsEmpty(collection) {
		return nil, nil
	}

	if utils.IsExpr(collection) {
		variables, err := execution.GetVariables()
		if err != nil {
			return nil, err
		}
		output := utils.GetStringSliceFromExpression(variables, collection)
		return output, nil
	} else {
		output := strings.Split(collection, ",")
		return output, nil
	}
}

func (f AbstractMultiInstanceActivityBehavior) setLoopVariable(execution delegate.DelegateExecution, variableName string, value interface{}) error {
	return execution.SetVariableLocal(variableName, value)
}

func (f AbstractMultiInstanceActivityBehavior) getLoopVariable(execution delegate.DelegateExecution, variableName string) (int, error) {
	value, err := f.getLoopVariableInstance(execution, variableName)
	if err != nil {
		return 0, err
	}

	return cast.ToInt(value), nil
}

func (f AbstractMultiInstanceActivityBehavior) getLoopVariableInstance(execution delegate.DelegateExecution, variableName string) (value interface{}, err error) {
	value, ok, err := execution.GetVariableLocal(variableName)
	if err != nil {
		return nil, err
	}

	parent, err := execution.GetParent()
	if err != nil {
		return nil, err
	}

	for !ok && parent != nil {
		value, ok, err = parent.GetVariableLocal(variableName)
		if err != nil {
			return nil, err
		}
		parent, err = parent.GetParent()
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (f AbstractMultiInstanceActivityBehavior) ExecuteOriginalBehavior(execution delegate.DelegateExecution, multiInstanceRootExecution delegate.DelegateExecution, loopCounter int) error {
	instances, err := f.resolveAndValidateCollection(execution)
	if err != nil {
		return err
	}
	value := instances[loopCounter]

	loopCharacteristics := f.Activity.GetLoopCharacteristics()
	if err = f.setLoopVariable(execution, loopCharacteristics.ElementVariable, value); err != nil {
		return err
	}

	execution.SetCurrentFlowElement(f.Activity)
	contextutil.GetAgenda().PlanContinueMultiInstanceOperation(execution, multiInstanceRootExecution, loopCounter)
	return nil
}

func (f AbstractMultiInstanceActivityBehavior) completionConditionSatisfied(multiInstanceRootExecution delegate.DelegateExecution) (bool, error) {
	loopCharacteristics := f.Activity.GetLoopCharacteristics()
	completionCondition := loopCharacteristics.CompletionCondition

	if stringutils.IsNotEmpty(completionCondition) {
		variables, err := multiInstanceRootExecution.GetVariablesLocal()
		if err != nil {
			return false, err
		}

		return utils.IsTrue(variables, completionCondition), nil
	}

	return false, nil
}
