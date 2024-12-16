package factory

import (
	. "github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	. "github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
)

var _ ActivityBehaviorFactory = (*DefaultActivityBehaviorFactory)(nil)

type DefaultActivityBehaviorFactory struct {
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateIntermediateCatchEventActivityBehavior(intermediateCatchEvent IntermediateCatchEvent) IntermediateCatchEventActivityBehavior {
	return IntermediateCatchEventActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition ConditionalEventDefinition) IntermediateCatchConditionalEventActivityBehavior {
	return NewIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition)
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateUserTaskActivityBehavior(userTask UserTask, key string) UserTaskActivityBehavior {
	return UserTaskActivityBehavior{UserTask: userTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateAutoUserTaskActivityBehavior(userTask UserTask, key string) UserAutoTaskActivityBehavior {
	return UserAutoTaskActivityBehavior{UserTask: userTask, ProcessKey: key}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateExclusiveGatewayActivityBehavior(exclusiveGateway ExclusiveGateway) ExclusiveGatewayActivityBehavior {
	return ExclusiveGatewayActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateInclusiveGatewayActivityBehavior(inclusiveGateway InclusiveGateway) InclusiveGatewayActivityBehavior {
	return InclusiveGatewayActivityBehavior{}
}

func (defaultActivityBehaviorFactory DefaultActivityBehaviorFactory) CreateParallelGatewayActivityBehavior(inclusiveGateway ParallelGateway) ParallelGatewayActivityBehavior {
	return ParallelGatewayActivityBehavior{}
}
