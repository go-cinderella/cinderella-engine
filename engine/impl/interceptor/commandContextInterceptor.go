package interceptor

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
)

type CommandContextInterceptor struct {
	Next                  engine.Interceptor
	CommandContextFactory engine.ICommandContextFactory
}

func (commandContextInterceptor CommandContextInterceptor) Execute(command engine.Command) (interface{}, error) {
	cContext, err := contextutil.GetCommandContext()
	if err != nil {
		cContext = commandContextInterceptor.CommandContextFactory.CreateCommandContext()
	}
	contextutil.SetCommandContext(cContext)
	defer func() {
		contextutil.RemoveCommandContext()
		contextutil.Clear()
	}()
	return commandContextInterceptor.Next.Execute(command)
}

func (commandContextInterceptor *CommandContextInterceptor) SetNext(next engine.Interceptor) {
	commandContextInterceptor.Next = next
}
