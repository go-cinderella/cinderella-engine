package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type TriggerExecutionOperation struct {
	AbstractOperation
}

func (trigger TriggerExecutionOperation) Run() (err error) {
	element := trigger.Execution.GetCurrentFlowElement()
	behavior := element.GetBehavior()
	operation := behavior.(delegate.TriggerableActivityBehavior)
	operation.Trigger(trigger.Execution)
	return err
}
