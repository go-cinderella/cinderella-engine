package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type ServiceTaskParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (eventParseHandler ServiceTaskParseHandler) GetHandledType() string {
	return model.ServiceTask{}.GetType()
}

func (eventParseHandler ServiceTaskParseHandler) ExecuteParse(bpmnParse *BpmnParse, baseElement delegate.BaseElement) {
	var behavior delegate.ActivityBehavior

	serviceTask := baseElement.(*model.ServiceTask)
	switch serviceTask.TaskType {
	case constant.SERVICE_TASK_HTTP:
		behavior = bpmnParse.ActivityBehaviorFactory.CreateHttpActivityBehavior(*serviceTask, bpmnParse.Name)
	case constant.SERVICE_TASK_PIPELINE:
		behavior = bpmnParse.ActivityBehaviorFactory.CreatePipelineActivityBehavior(*serviceTask, bpmnParse.Name)
	default:
		panic("Unknown service task type")
	}

	serviceTask.SetBehavior(behavior)
}
