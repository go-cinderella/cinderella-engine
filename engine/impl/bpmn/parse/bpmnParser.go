package parse

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/parse/factory"
)

type BpmnParser struct {
	ActivityBehaviorFactory factory.ActivityBehaviorFactory
	BpmnParserHandlers      BpmnParseHandlers
	BpmnParseFactory        BpmnParseFactory
}

func (bpmnParser BpmnParser) CreateParse() BpmnParse {
	return bpmnParser.BpmnParseFactory.CreateBpmnParse(bpmnParser)
}

func (bpmnParser BpmnParser) SetActivityBehaviorFactory(activityBehaviorFactory factory.ActivityBehaviorFactory) {
	bpmnParser.ActivityBehaviorFactory = activityBehaviorFactory
}

func (bpmnParser BpmnParser) GetActivityBehaviorFactory() factory.ActivityBehaviorFactory {
	return bpmnParser.ActivityBehaviorFactory
}

func (bpmnParser BpmnParser) GetBpmnParserHandlers() BpmnParseHandlers {
	return bpmnParser.BpmnParserHandlers
}

func (bpmnParser BpmnParser) SetBpmnParserHandlers(bpmnParseHandlers BpmnParseHandlers) {
	bpmnParser.BpmnParserHandlers = bpmnParseHandlers
}
