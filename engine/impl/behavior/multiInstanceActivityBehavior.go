package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

const (
	NUMBER_OF_INSTANCES           string = "nrOfInstances"
	NUMBER_OF_ACTIVE_INSTANCES    string = "nrOfActiveInstances"
	NUMBER_OF_COMPLETED_INSTANCES string = "nrOfCompletedInstances"
)

// MultiInstanceActivityBehavior Implementation of the multi-instance functionality as described in the BPMN 2.0 spec.
type MultiInstanceActivityBehavior struct {
	flowNodeActivityBehavior

	// Instance members
	Activity              model.Activity
	InnerActivityBehavior delegate.TriggerableActivityBehavior
}

func (f MultiInstanceActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	return nil
}

func (f MultiInstanceActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	return nil
}

func (f MultiInstanceActivityBehavior) leave(execution delegate.DelegateExecution) error {
	return nil
}

func (f MultiInstanceActivityBehavior) getLocalLoopVariable(execution delegate.DelegateExecution, variableName string) (value int, err error) {
	return 0, nil

}
