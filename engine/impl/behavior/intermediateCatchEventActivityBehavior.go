package behavior

import "github.com/go-cinderella/cinderella-engine/engine/impl/delegate"

var _ delegate.ActivityBehavior = (*IntermediateCatchEventActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*IntermediateCatchEventActivityBehavior)(nil)

type IntermediateCatchEventActivityBehavior struct {
	abstractBpmnActivityBehavior
}

func (f IntermediateCatchEventActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	return f.leaveIntermediateCatchEvent(execution)
}

func (f IntermediateCatchEventActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	return nil
}

func (f IntermediateCatchEventActivityBehavior) leaveIntermediateCatchEvent(execution delegate.DelegateExecution) error {
	// EventGateway TODO

	return f.leave(execution)
}
