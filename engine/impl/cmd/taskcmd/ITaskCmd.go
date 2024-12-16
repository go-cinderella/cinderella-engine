package taskcmd

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
)

type ITaskCmd interface {
	TaskExecute(commandContext engine.Context, entity entitymanager.TaskEntity) (interface{}, error)
}
