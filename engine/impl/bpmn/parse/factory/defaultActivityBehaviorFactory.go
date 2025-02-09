package factory

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ ActivityBehaviorFactory = (*DefaultActivityBehaviorFactory)(nil)

type DefaultActivityBehaviorFactory struct {
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateSequentialMultiInstanceBehavior(activity delegate.FlowElement, innerActivityBehavior delegate.ActivityBehavior) *behavior.SequentialMultiInstanceBehavior {
	seq := &behavior.SequentialMultiInstanceBehavior{
		AbstractMultiInstanceActivityBehavior: behavior.AbstractMultiInstanceActivityBehavior{
			Activity:              activity.(behavior.MultiInstanceActivity),
			InnerActivityBehavior: innerActivityBehavior.(behavior.MultiInstanceSupportBehavior),
		},
	}

	seq.AbstractMultiInstanceActivityBehavior.Impl = seq
	seq.AbstractMultiInstanceActivityBehavior.InnerActivityBehavior.SetMultiInstanceActivityBehavior(seq)

	return seq
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateParallelMultiInstanceBehavior(activity delegate.FlowElement, innerActivityBehavior delegate.ActivityBehavior) *behavior.ParallelMultiInstanceBehavior {
	para := &behavior.ParallelMultiInstanceBehavior{
		AbstractMultiInstanceActivityBehavior: behavior.AbstractMultiInstanceActivityBehavior{
			Activity:              activity.(behavior.MultiInstanceActivity),
			InnerActivityBehavior: innerActivityBehavior.(behavior.MultiInstanceSupportBehavior),
		},
	}

	para.AbstractMultiInstanceActivityBehavior.Impl = para
	para.AbstractMultiInstanceActivityBehavior.InnerActivityBehavior.SetMultiInstanceActivityBehavior(para)

	return para
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateHttpActivityBehavior(serviceTask model.ServiceTask, key string) *behavior.HttpServiceTaskActivityBehavior {
	return &behavior.HttpServiceTaskActivityBehavior{ServiceTask: serviceTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreatePipelineActivityBehavior(serviceTask model.ServiceTask, key string) *behavior.PipelineServiceTaskActivityBehavior {
	return &behavior.PipelineServiceTaskActivityBehavior{ServiceTask: serviceTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateIntermediateCatchEventActivityBehavior(intermediateCatchEvent model.IntermediateCatchEvent) *behavior.IntermediateCatchEventActivityBehavior {
	return &behavior.IntermediateCatchEventActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition model.ConditionalEventDefinition) *behavior.IntermediateCatchConditionalEventActivityBehavior {
	return behavior.NewIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition)
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateUserTaskActivityBehavior(userTask model.UserTask, key string) *behavior.UserTaskActivityBehavior {
	return &behavior.UserTaskActivityBehavior{UserTask: userTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateAutoUserTaskActivityBehavior(userTask model.UserTask, key string) *behavior.UserAutoTaskActivityBehavior {
	return &behavior.UserAutoTaskActivityBehavior{UserTask: userTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateExclusiveGatewayActivityBehavior(exclusiveGateway model.ExclusiveGateway) *behavior.ExclusiveGatewayActivityBehavior {
	return &behavior.ExclusiveGatewayActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateInclusiveGatewayActivityBehavior(inclusiveGateway model.InclusiveGateway) *behavior.InclusiveGatewayActivityBehavior {
	return &behavior.InclusiveGatewayActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateParallelGatewayActivityBehavior(inclusiveGateway model.ParallelGateway) *behavior.ParallelGatewayActivityBehavior {
	return &behavior.ParallelGatewayActivityBehavior{}
}
