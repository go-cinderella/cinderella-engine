package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/runnable"
)

type CinderellaEngineAgenda interface {
	IsEmpty() bool

	PlanOperation(operation runnable.Operation)

	GetNextOperation() runnable.Operation

	PlanContinueProcessOperation(execution delegate.DelegateExecution)

	PlanEvaluateConditionalEventsOperation(execution delegate.DelegateExecution)

	//planContinueProcessSynchronousOperation(execution ExecutionEntity)
	//
	//planContinueProcessInCompensation(execution ExecutionEntity)
	//
	//planContinueMultiInstanceOperation(execution ExecutionEntity)

	PlanTakeOutgoingSequenceFlowsOperation(execution delegate.DelegateExecution, evaluateConditions bool)

	PlanEndExecutionOperation(execution delegate.DelegateExecution)

	PlanTriggerExecutionOperation(execution delegate.DelegateExecution)
	//
	//planDestroyScopeOperation(execution ExecutionEntity)
	//
	//planExecuteInactiveBehaviorsOperation()
}
