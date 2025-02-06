package agenda

import (
	"container/list"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/agenda"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/runnable"
)

var _ agenda.CinderellaEngineAgenda = (*DefaultCinderellaEngineAgenda)(nil)

type DefaultCinderellaEngineAgenda struct {
	operations list.List
	ctx        engine.Context
}

func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) PlanContinueMultiInstanceOperation(execution delegate.DelegateExecution, multiInstanceRootExecution delegate.DelegateExecution, loopCounter int) {
	abstractOperation := NewAbstractOperation(execution, defaultCinderellaEngineAgenda, WithContext(defaultCinderellaEngineAgenda.ctx))
	defaultCinderellaEngineAgenda.PlanOperation(&ContinueMultiInstanceOperation{AbstractOperation: abstractOperation, MultiInstanceRootExecution: multiInstanceRootExecution, LoopCounter: loopCounter})
}

func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) GetContext() engine.Context {
	return defaultCinderellaEngineAgenda.ctx
}

func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) SetContext(ctx engine.Context) {
	defaultCinderellaEngineAgenda.ctx = ctx
}

// 判断是否为空
func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) IsEmpty() bool {
	return defaultCinderellaEngineAgenda.operations.Len() == 0
}

// 设置后续操作
func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) PlanOperation(operation runnable.Operation) {
	defaultCinderellaEngineAgenda.operations.PushBack(operation)
}

func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) GetNextOperation() runnable.Operation {
	value := defaultCinderellaEngineAgenda.operations.Front()
	defaultCinderellaEngineAgenda.operations.Remove(value)
	return value.Value.(runnable.Operation)
}

// 连线继续执行
func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) PlanContinueProcessOperation(execution delegate.DelegateExecution) {
	abstractOperation := NewAbstractOperation(execution, defaultCinderellaEngineAgenda, WithContext(defaultCinderellaEngineAgenda.ctx))
	defaultCinderellaEngineAgenda.PlanOperation(&ContinueProcessOperation{AbstractOperation: abstractOperation})
}

// 任务出口执行
func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) PlanTriggerExecutionOperation(execution delegate.DelegateExecution) {
	abstractOperation := NewAbstractOperation(execution, defaultCinderellaEngineAgenda, WithContext(defaultCinderellaEngineAgenda.ctx))
	defaultCinderellaEngineAgenda.PlanOperation(&TriggerExecutionOperation{AbstractOperation: abstractOperation})
}

// 连线出口设置
func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) PlanTakeOutgoingSequenceFlowsOperation(execution delegate.DelegateExecution, valuateConditions bool) {
	abstractOperation := NewAbstractOperation(execution, defaultCinderellaEngineAgenda, WithContext(defaultCinderellaEngineAgenda.ctx))
	defaultCinderellaEngineAgenda.PlanOperation(&TakeOutgoingSequenceFlowsOperation{AbstractOperation: abstractOperation, EvaluateConditions: valuateConditions})
}

// 任务结束
func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) PlanEndExecutionOperation(execution delegate.DelegateExecution) {
	abstractOperation := NewAbstractOperation(execution, defaultCinderellaEngineAgenda, WithContext(defaultCinderellaEngineAgenda.ctx))
	defaultCinderellaEngineAgenda.PlanOperation(&EndExecutionOperation{AbstractOperation: abstractOperation})
}

func (defaultCinderellaEngineAgenda *DefaultCinderellaEngineAgenda) PlanEvaluateConditionalEventsOperation(execution delegate.DelegateExecution) {
	abstractOperation := NewAbstractOperation(execution, defaultCinderellaEngineAgenda, WithContext(defaultCinderellaEngineAgenda.ctx))
	defaultCinderellaEngineAgenda.PlanOperation(&EvaluateConditionalEventsOperation{AbstractOperation: abstractOperation})
}
