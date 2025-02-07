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

	if getter, ok := element.(model.LoopCharacteristicsGetter); ok && getter.GetLoopCharacteristics() != nil {
		abstractBpmnParse.createMultiInstanceLoopCharacteristics(bpmnParse, getter)
	}
}

func (abstractBpmnParse AbstractBpmnParseHandler) createMultiInstanceLoopCharacteristics(bpmnParse *BpmnParse, getter model.LoopCharacteristicsGetter) {
	loopCharacteristics := getter.GetLoopCharacteristics()
	activity := getter.(delegate.FlowElement)
	activityBehavior := abstractBpmnParse.createMultiInstanceActivityBehavior(activity, loopCharacteristics, bpmnParse)
	activity.SetBehavior(activityBehavior)
}

func (abstractBpmnParse AbstractBpmnParseHandler) createMultiInstanceActivityBehavior(activity delegate.FlowElement, loopCharacteristics *model.MultiInstanceLoopCharacteristics, bpmnParse *BpmnParse) delegate.ActivityBehavior {
	clonedActivity := activity.(delegate.Cloneable).Clone()

	activityBehavior := clonedActivity.GetBehavior()

	if loopCharacteristics.IsSequential {
		return bpmnParse.ActivityBehaviorFactory.CreateSequentialMultiInstanceBehavior(clonedActivity, activityBehavior)
	} else {
		return bpmnParse.ActivityBehaviorFactory.CreateParallelMultiInstanceBehavior(clonedActivity, activityBehavior)
	}
}
