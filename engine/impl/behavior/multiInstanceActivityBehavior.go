package behavior

import "github.com/go-cinderella/cinderella-engine/engine/impl/delegate"

// MultiInstanceActivityBehavior TODO
type MultiInstanceActivityBehavior struct {
	flowNodeActivityBehavior
}

func (f MultiInstanceActivityBehavior) leave(execution delegate.DelegateExecution) error {
	panic("implement me")
}
