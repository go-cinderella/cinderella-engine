package model

type Activity struct {
	FlowNode
	MultiInstance *MultiInstanceLoopCharacteristics `xml:"multiInstanceLoopCharacteristics"`
}

func (a Activity) GetLoopCharacteristics() *MultiInstanceLoopCharacteristics {
	return a.MultiInstance
}
