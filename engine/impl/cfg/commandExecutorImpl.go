package cfg

import "github.com/go-cinderella/cinderella-engine/engine"

var _ engine.Executor = (*CommandExecutorImpl)(nil)

var commandExecutorImpl CommandExecutorImpl

type CommandExecutorImpl struct {
	First engine.Interceptor
}

func SetCommandExecutorImpl(commandExecutor CommandExecutorImpl) {
	commandExecutorImpl = commandExecutor
}

func GetCommandExecutorImpl() CommandExecutorImpl {
	return commandExecutorImpl
}

func (comm CommandExecutorImpl) Exe(conf engine.Command) (interface{}, error) {
	return comm.First.Execute(conf)
}
