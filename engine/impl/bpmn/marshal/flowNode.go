package marshal

// FlowNode 父类实现体
type FlowNode struct {
	DefaultBaseElement

	Incoming []string `xml:"bpmn:incoming"`
	Outgoing []string `xml:"bpmn:outgoing"`
}
