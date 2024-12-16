package agenda

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/agenda"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type AbstractOperation struct {
	Agenda    agenda.CinderellaEngineAgenda
	ctx       engine.Context
	Execution delegate.DelegateExecution
}

func (abstractOperation *AbstractOperation) GetAgenda() agenda.CinderellaEngineAgenda {
	return abstractOperation.Agenda
}

func (abstractOperation *AbstractOperation) GetContext() engine.Context {
	return abstractOperation.ctx
}

//func (abstractOperation *AbstractOperation) GetCommandExecutor() engine.Executor {
//	return abstractOperation.context.commandExecutor
//}

func (abstractOperation *AbstractOperation) SetAgenda(cinderellaEngineAgenda agenda.CinderellaEngineAgenda) {
	abstractOperation.Agenda = cinderellaEngineAgenda
}

//func (abstractOperation *AbstractOperation) SetCommandExecutor(commandExecutor engine.Executor) {
//	abstractOperation.context.commandExecutor = commandExecutor
//}

type OperationOption func(*AbstractOperation)

func WithContext(ctx engine.Context) OperationOption {
	return func(operation *AbstractOperation) {
		operation.ctx = ctx
	}
}

//func WithCommandExecutor(commandExecutor engine.Executor) OperationOption {
//	return func(operation *AbstractOperation) {
//		operation.SetCommandExecutor(commandExecutor)
//	}
//}

func NewAbstractOperation(execution delegate.DelegateExecution, cinderellaEngineAgenda agenda.CinderellaEngineAgenda, options ...OperationOption) AbstractOperation {
	abstractOperation := AbstractOperation{
		Execution: execution,
		Agenda:    cinderellaEngineAgenda,
	}

	for _, option := range options {
		option(&abstractOperation)
	}

	return abstractOperation
}
