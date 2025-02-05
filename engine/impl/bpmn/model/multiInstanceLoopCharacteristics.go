package model

// MultiInstanceLoopCharacteristics
//
//	nrOfInstances：实例总数
//	nrOfActiveInstances：当前活动的（尚未完成的）实例数量。对于串行多实例来说，这个值始终是 1。
//	nrOfCompletedInstances：已经完成的实例的数量。
//	loopCounter：给定实例在 for-each 循环中的 index。
type MultiInstanceLoopCharacteristics struct {
	// LoopCardinality TODO
	LoopCardinality int `xml:"loopCardinality"`

	IsSequential    bool   `xml:"isSequential,attr"`
	Collection      string `xml:"collection,attr"`
	ElementVariable string `xml:"elementVariable,attr"`

	// if empty, means all
	// e.g.:
	// ${(nrOfCompletedInstances / nrOfInstances) * 100 &gt; 50}
	// ${nrOfCompletedInstances == 5}
	CompletionCondition string `xml:"completionCondition"`
}
