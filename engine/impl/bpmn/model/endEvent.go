package model

type EndEvent struct {
	FlowNode
}

func (endEvent EndEvent) GetType() string {
	return "endEvent"
}
