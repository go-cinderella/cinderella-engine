package behavior

import "github.com/go-cinderella/cinderella-engine/engine/impl/delegate"

type IMultiInstanceActivityBehavior interface {
	CreateInstances(execution delegate.DelegateExecution) (int, error)

	Leave(execution delegate.DelegateExecution) error
}
