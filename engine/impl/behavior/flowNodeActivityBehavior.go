package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ delegate.ActivityBehavior = (*flowNodeActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*flowNodeActivityBehavior)(nil)

type flowNodeActivityBehavior struct {
	bpmnActivityBehavior
}

func (f flowNodeActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (f flowNodeActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	return f.leave(execution)
}

func (f flowNodeActivityBehavior) leave(execution delegate.DelegateExecution) error {
	f.bpmnActivityBehavior.performDefaultOutgoingBehavior(execution)
	return nil
}
