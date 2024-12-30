package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/samber/lo"
)

var _ delegate.BaseElement = (*FlowNode)(nil)
var _ delegate.FlowElement = (*FlowNode)(nil)

// 父类实现体
type FlowNode struct {
	delegate.BaseHandlerType
	DefaultBaseElement
	IncomingFlow      []delegate.FlowElement
	OutgoingFlow      []delegate.FlowElement
	SourceFlowElement delegate.FlowElement
	TargetFlowElement delegate.FlowElement
	Behavior          delegate.ActivityBehavior

	Incoming []string `xml:"incoming"`
	Outgoing []string `xml:"outgoing"`
}

func (flow *FlowNode) SetIncoming(f []delegate.FlowElement) {
	flow.IncomingFlow = f
	flow.Incoming = lo.Map[delegate.FlowElement, string](f, func(item delegate.FlowElement, index int) string {
		return item.GetId()
	})
}
func (flow *FlowNode) SetOutgoing(f []delegate.FlowElement) {
	flow.OutgoingFlow = f
	flow.Outgoing = lo.Map[delegate.FlowElement, string](f, func(item delegate.FlowElement, index int) string {
		return item.GetId()
	})
}

func (flow *FlowNode) GetIncoming() []delegate.FlowElement {
	return flow.IncomingFlow
}
func (flow *FlowNode) GetOutgoing() []delegate.FlowElement {
	return flow.OutgoingFlow
}

func (flow *FlowNode) SetSourceFlowElement(f delegate.FlowElement) {
	flow.SourceFlowElement = f
}
func (flow *FlowNode) SetTargetFlowElement(f delegate.FlowElement) {
	flow.TargetFlowElement = f
}

func (flow *FlowNode) GetSourceFlowElement() delegate.FlowElement {
	return flow.SourceFlowElement
}
func (flow *FlowNode) GetTargetFlowElement() delegate.FlowElement {
	return flow.TargetFlowElement
}

func (flow *FlowNode) GetId() string {
	return flow.Id
}

func (flow *FlowNode) GetName() string {
	return flow.Name
}

func (flow *FlowNode) GetBehavior() delegate.ActivityBehavior {
	return flow.Behavior
}

func (flow *FlowNode) SetBehavior(behavior delegate.ActivityBehavior) {
	flow.Behavior = behavior
}

func (flow *FlowNode) GetHandlerType() string {
	return flow.BaseHandlerType.GetType()
}
