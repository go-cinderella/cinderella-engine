package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type bpmnActivityBehavior struct {
}

func (behavior bpmnActivityBehavior) performDefaultOutgoingBehavior(activityExecution delegate.DelegateExecution) {
	behavior.performOutgoingBehavior(activityExecution, true)
}

func (behavior bpmnActivityBehavior) performIgnoreConditionsOutgoingBehavior(activityExecution delegate.DelegateExecution) {
	behavior.performOutgoingBehavior(activityExecution, false)
}

func (behavior bpmnActivityBehavior) performOutgoingBehavior(activityExecution delegate.DelegateExecution, checkConditions bool) {
	contextutil.GetAgenda().PlanTakeOutgoingSequenceFlowsOperation(activityExecution, checkConditions)
}
