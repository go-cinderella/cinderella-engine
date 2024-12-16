package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"strings"
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
}

func (flow *FlowNode) SetIncoming(f []delegate.FlowElement) {
	flow.IncomingFlow = f
}
func (flow *FlowNode) SetOutgoing(f []delegate.FlowElement) {
	flow.OutgoingFlow = f
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

func (flow FlowNode) String() string {
	var sb strings.Builder
	sb.WriteString("FlowNode")
	sb.WriteString("{")
	sb.WriteString("DefaultBaseElement: ")
	sb.WriteString(flow.DefaultBaseElement.String())
	sb.WriteString("}")
	return sb.String()
}
