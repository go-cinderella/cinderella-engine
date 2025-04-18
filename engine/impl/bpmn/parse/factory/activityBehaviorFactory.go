package factory

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/behavior"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ActivityBehaviorFactory interface {
	CreateUserTaskActivityBehavior(userTask model.UserTask, key string) *behavior.UserTaskActivityBehavior

	CreateSequentialMultiInstanceBehavior(activity delegate.FlowElement, innerActivityBehavior delegate.ActivityBehavior) *behavior.SequentialMultiInstanceBehavior

	CreateParallelMultiInstanceBehavior(activity delegate.FlowElement, innerActivityBehavior delegate.ActivityBehavior) *behavior.ParallelMultiInstanceBehavior

	CreateAutoUserTaskActivityBehavior(userTask model.UserTask, key string) *behavior.UserAutoTaskActivityBehavior

	CreateExclusiveGatewayActivityBehavior(exclusiveGateway model.ExclusiveGateway) *behavior.ExclusiveGatewayActivityBehavior

	CreateInclusiveGatewayActivityBehavior(inclusiveGateway model.InclusiveGateway) *behavior.InclusiveGatewayActivityBehavior

	CreateParallelGatewayActivityBehavior(inclusiveGateway model.ParallelGateway) *behavior.ParallelGatewayActivityBehavior

	CreateIntermediateCatchConditionalEventActivityBehavior(conditionalEventDefinition model.ConditionalEventDefinition) *behavior.IntermediateCatchConditionalEventActivityBehavior

	CreateIntermediateCatchEventActivityBehavior(intermediateCatchEvent model.IntermediateCatchEvent) *behavior.IntermediateCatchEventActivityBehavior

	CreateHttpActivityBehavior(serviceTask model.ServiceTask, key string) *behavior.HttpServiceTaskActivityBehavior

	CreatePipelineActivityBehavior(serviceTask model.ServiceTask, key string) *behavior.PipelineServiceTaskActivityBehavior
}
