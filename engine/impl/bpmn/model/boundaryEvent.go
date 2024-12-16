package model

type BoundaryEvent struct {
	*FlowNode
	AttachedToRef        string               `xml:"attachedToRef,attr"`
	CancelActivity       string               `xml:"cancelActivity,attr"`
	TimerEventDefinition TimerEventDefinition `xml:"timerEventDefinition"`
}
