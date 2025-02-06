package behavior

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ IMultiInstanceActivityBehavior = (*ParallelMultiInstanceBehavior)(nil)

type ParallelMultiInstanceBehavior struct {
	AbstractMultiInstanceActivityBehavior
}

func (p ParallelMultiInstanceBehavior) CreateInstances(execution delegate.DelegateExecution) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (p ParallelMultiInstanceBehavior) Leave(execution delegate.DelegateExecution) error {
	//TODO implement me
	panic("implement me")
}
