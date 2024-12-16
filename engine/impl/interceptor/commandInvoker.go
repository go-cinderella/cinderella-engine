package interceptor

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
)

type CommandInvoker struct {
	Next engine.Interceptor
}

func (commandInvoker CommandInvoker) Execute(command engine.Command) (result interface{}, err error) {
	context, err := contextutil.GetCommandContext()
	if err != nil {
		return nil, err
	}
	result, err = command.Execute(context)
	if err != nil {
		return nil, err
	}
	err = executeOperations(context)
	return result, err
}

func executeOperations(context engine.Context) (err error) {
	for !context.GetAgenda().IsEmpty() {
		err = context.GetAgenda().GetNextOperation().Run()
		if err != nil {
			return err
		}
	}
	return err
}

func (commandInvoker *CommandInvoker) SetNext(next engine.Interceptor) {
	commandInvoker.Next = next
}
