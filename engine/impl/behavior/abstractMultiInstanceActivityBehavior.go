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

type IMultiInstanceActivity interface {
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

// AbstractMultiInstanceActivityBehavior Implementation of the multi-instance functionality as described in the BPMN 2.0 spec.
type AbstractMultiInstanceActivityBehavior struct {
	flowNodeActivityBehavior

	Impl IMultiInstanceActivityBehavior
	// Instance members
	Activity              IMultiInstanceActivity
	InnerActivityBehavior delegate.TriggerableActivityBehavior

	IsSequential        bool
	Collection          string
	ElementVariable     string
	CompletionCondition string
}

func (f AbstractMultiInstanceActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	return nil
}

func (f AbstractMultiInstanceActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	_, ok, err := f.getLocalLoopVariable(execution, CollectionElementIndexVariable)
	if err != nil {
		return err
	}
	if !ok {
		nrOfInstances, err := f.Impl.CreateInstances(execution)
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
	return nil
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
	executionEntityManager.DeleteChildExecution(multiInstanceRootExecution.GetExecutionId(), lo.ToPtr(DELETE_REASON_END))
	executionEntityManager.DeleteRelatedDataForExecution(multiInstanceRootExecution.GetExecutionId(), lo.ToPtr(DELETE_REASON_END))
	executionEntityManager.DeleteExecution(multiInstanceRootExecution.GetExecutionId())

	flowElement := multiInstanceRootExecution.GetCurrentFlowElement()
	parentExecution, err := multiInstanceRootExecution.GetParent()
	if err != nil {
		return err
	}

	newExecution, err := executionEntityManager.CreateChildExecution(parentExecution)
	if err != nil {
		return err
	}
	newExecution.SetCurrentFlowElement(flowElement)

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
		variables, err := execution.GetProcessVariables()
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

func (f AbstractMultiInstanceActivityBehavior) SetLoopVariable(execution delegate.DelegateExecution, variableName string, value interface{}) error {
	return execution.SetVariableLocal(variableName, value)
}

func (f AbstractMultiInstanceActivityBehavior) ExecuteOriginalBehavior(execution delegate.DelegateExecution, multiInstanceRootExecution delegate.DelegateExecution, loopCounter int) error {
	instances, err := f.resolveAndValidateCollection(execution)
	if err != nil {
		return err
	}
	value := instances[loopCounter]

	loopCharacteristics := f.Activity.GetLoopCharacteristics()
	if err = f.SetLoopVariable(execution, loopCharacteristics.ElementVariable, value); err != nil {
		return err
	}

	execution.SetCurrentFlowElement(f.Activity)
	contextutil.GetAgenda().PlanContinueMultiInstanceOperation(execution, multiInstanceRootExecution, loopCounter)
	return nil
}
