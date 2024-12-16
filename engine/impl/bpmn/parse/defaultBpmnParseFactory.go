package parse

type DefaultBpmnParseFactory struct {
}

func (DefaultBpmnParseFactory) CreateBpmnParse(bpmnParser BpmnParser) BpmnParse {
	return BpmnParse{ActivityBehaviorFactory: bpmnParser.ActivityBehaviorFactory,
		BpmnParserHandlers: bpmnParser.BpmnParserHandlers}
}
