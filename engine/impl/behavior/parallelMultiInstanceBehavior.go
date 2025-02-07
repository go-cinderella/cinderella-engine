package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ multiInstanceActivityBehavior = (*ParallelMultiInstanceBehavior)(nil)

type ParallelMultiInstanceBehavior struct {
	AbstractMultiInstanceActivityBehavior
}

func (p ParallelMultiInstanceBehavior) createInstances(execution delegate.DelegateExecution) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (p ParallelMultiInstanceBehavior) leave(execution delegate.DelegateExecution) error {
	//TODO implement me
	panic("implement me")
}
