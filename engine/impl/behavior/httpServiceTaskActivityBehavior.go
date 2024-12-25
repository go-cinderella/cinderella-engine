package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ delegate.ActivityBehavior = (*HttpServiceTaskActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*HttpServiceTaskActivityBehavior)(nil)

type HttpServiceTaskActivityBehavior struct {
	abstractBpmnActivityBehavior
	ServiceTask model.ServiceTask
	ProcessKey  string
}

// Execute TODO
func (user HttpServiceTaskActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	return nil
}
