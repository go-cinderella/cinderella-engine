package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

type AbstractBpmnParseHandler struct {
	ParseHandler
}

func (abstractBpmnParse AbstractBpmnParseHandler) GetHandledTypes() []string {
	types := make([]string, 0)
	types = append(types, abstractBpmnParse.ParseHandler.GetHandledType())
	return types
}
func (abstractBpmnParse AbstractBpmnParseHandler) Parse(bpmnParse *BpmnParse, element delegate.BaseElement) {
	abstractBpmnParse.ExecuteParse(bpmnParse, element)

	if getter, ok := element.(activityGetter); ok && getter.GetActivity().GetLoopCharacteristics() != nil {
		abstractBpmnParse.createMultiInstanceLoopCharacteristics(bpmnParse, getter.GetActivity())
	}
}

func (abstractBpmnParse AbstractBpmnParseHandler) createMultiInstanceLoopCharacteristics(bpmnParse *BpmnParse, activity model.Activity) {
	loopCharacteristics := activity.GetLoopCharacteristics()
	activityBehavior := abstractBpmnParse.createMultiInstanceActivityBehavior(activity, loopCharacteristics, bpmnParse)
	activity.SetBehavior(activityBehavior)
}

func (abstractBpmnParse AbstractBpmnParseHandler) createMultiInstanceActivityBehavior(activity model.Activity, loopCharacteristics *model.MultiInstanceLoopCharacteristics, bpmnParse *BpmnParse) delegate.ActivityBehavior {
	activityBehavior := activity.GetBehavior()
	triggerableActivityBehavior := activityBehavior.(delegate.TriggerableActivityBehavior)
	if loopCharacteristics.IsSequential {
		return bpmnParse.ActivityBehaviorFactory.CreateSequentialMultiInstanceBehavior(activity, triggerableActivityBehavior)
	} else {
		return bpmnParse.ActivityBehaviorFactory.CreateParallelMultiInstanceBehavior(activity, triggerableActivityBehavior)
	}
}
