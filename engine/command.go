package engine

import (
	"context"
)

// Command can be directly executed but cannot be executed by ComandExecutor in another command,
// otherwise CommandContext of parent command may be deleted mistakenly by interceptors.
type Command interface {
	Execute(commandContext Context) (interface{}, error)
	Context() context.Context
	IsTransactional() bool
}
