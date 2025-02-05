package factory

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var _ ActivityBehaviorFactory = (*DefaultActivityBehaviorFactory)(nil)

type DefaultActivityBehaviorFactory struct {
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateSequentialMultiInstanceBehavior(activity model.Activity, innerActivityBehavior delegate.TriggerableActivityBehavior) behavior.SequentialMultiInstanceBehavior {
	return behavior.SequentialMultiInstanceBehavior{
		MultiInstanceActivityBehavior: behavior.MultiInstanceActivityBehavior{
			Activity:              activity,
			InnerActivityBehavior: innerActivityBehavior,
		},
	}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateParallelMultiInstanceBehavior(activity model.Activity, innerActivityBehavior delegate.TriggerableActivityBehavior) behavior.ParallelMultiInstanceBehavior {
	return behavior.ParallelMultiInstanceBehavior{
		MultiInstanceActivityBehavior: behavior.MultiInstanceActivityBehavior{
			Activity:              activity,
			InnerActivityBehavior: innerActivityBehavior,
		},
	}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateHttpActivityBehavior(serviceTask model.ServiceTask, key string) behavior.HttpServiceTaskActivityBehavior {
	return behavior.HttpServiceTaskActivityBehavior{ServiceTask: serviceTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreatePipelineActivityBehavior(serviceTask model.ServiceTask, key string) behavior.PipelineServiceTaskActivityBehavior {
	return behavior.PipelineServiceTaskActivityBehavior{ServiceTask: serviceTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateIntermediateCatchEventActivityBehavior(intermediateCatchEvent model.IntermediateCatchEvent) behavior.IntermediateCatchEventActivityBehavior {
	return behavior.IntermediateCatchEventActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition model.ConditionalEventDefinition) behavior.IntermediateCatchConditionalEventActivityBehavior {
	return behavior.NewIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition)
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateUserTaskActivityBehavior(userTask model.UserTask, key string) behavior.UserTaskActivityBehavior {
	return behavior.UserTaskActivityBehavior{UserTask: userTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateAutoUserTaskActivityBehavior(userTask model.UserTask, key string) behavior.UserAutoTaskActivityBehavior {
	return behavior.UserAutoTaskActivityBehavior{UserTask: userTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateExclusiveGatewayActivityBehavior(exclusiveGateway model.ExclusiveGateway) behavior.ExclusiveGatewayActivityBehavior {
	return behavior.ExclusiveGatewayActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateInclusiveGatewayActivityBehavior(inclusiveGateway model.InclusiveGateway) behavior.InclusiveGatewayActivityBehavior {
	return behavior.InclusiveGatewayActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateParallelGatewayActivityBehavior(inclusiveGateway model.ParallelGateway) behavior.ParallelGatewayActivityBehavior {
	return behavior.ParallelGatewayActivityBehavior{}
}
