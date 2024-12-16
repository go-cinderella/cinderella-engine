package contextutil

import (
	"container/list"
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/internal/errs"
	"github.com/go-cinderella/cinderella-engine/engine/runtime"
	"sync"
)

var commandStackHolder *sync.Map

func init() {
	commandStackHolder = new(sync.Map)
}

func RemoveCommandContext() {
	commandStack := getStack(commandStackHolder)
	commandStack.Remove(commandStack.Front())
}

func SetCommandContext(commandContext engine.Context) {
	getStack(commandStackHolder).PushFront(commandContext)
}

func GetCommandContext() (engine.Context, error) {
	commandStack := getStack(commandStackHolder)
	if commandStack.Len() <= 0 {
		return nil, errs.CinderellaError{}
	}
	return commandStack.Front().Value.(engine.Context), nil
}

func MustGetCommandContext() engine.Context {
	commandStack := getStack(commandStackHolder)
	return commandStack.Front().Value.(engine.Context)
}

func getStack(commandStackHolder *sync.Map) *list.List {
	commandStack, ok := commandStackHolder.Load(runtime.GoroutineId())
	if !ok {
		commandStack = list.New()
		commandStackHolder.Store(runtime.GoroutineId(), commandStack)
	}
	return commandStack.(*list.List)
}

func Clear() {
	commandStack, ok := commandStackHolder.Load(runtime.GoroutineId())
	if !ok {
		return
	}
	commandStack.(*list.List).Init()
	commandStackHolder.Delete(runtime.GoroutineId())
}
