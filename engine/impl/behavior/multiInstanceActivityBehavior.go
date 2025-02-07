package behavior

import "github.com/go-cinderella/cinderella-engine/engine/impl/delegate"

type multiInstanceActivityBehavior interface {
	createInstances(execution delegate.DelegateExecution) (int, error)

	leave(execution delegate.DelegateExecution) error
}
