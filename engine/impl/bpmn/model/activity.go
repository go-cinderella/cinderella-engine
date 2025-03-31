package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/shopspring/decimal"
)

var _ delegate.MainNodeSupported = (*Activity)(nil)

type Activity struct {
	FlowNode
	IsMainNode       *bool                             `xml:"isMainNode,attr"`
	MainNodeWeight   *decimal.Decimal                  `xml:"mainNodeWeight,attr"`
	MainNodeIndex    *int                              `xml:"mainNodeIndex,attr"`
	MainNodeProgress *decimal.Decimal                  `xml:"mainNodeProgress,attr"`
	MultiInstance    *MultiInstanceLoopCharacteristics `xml:"multiInstanceLoopCharacteristics"`
}

func (a Activity) GetLoopCharacteristics() *MultiInstanceLoopCharacteristics {
	return a.MultiInstance
}

func (a Activity) HasMultiInstanceLoopCharacteristics() bool {
	return a.GetLoopCharacteristics() != nil
}

func (a Activity) GetIsMainNode() bool {
	return a.IsMainNode != nil && *a.IsMainNode
}

func (a Activity) GetMainNodeWeight() decimal.Decimal {
	if a.MainNodeWeight == nil {
		return decimal.Zero
	}
	return *a.MainNodeWeight
}

func (a Activity) GetMainNodeIndex() int {
	if a.MainNodeIndex == nil {
		return 0
	}
	return *a.MainNodeIndex
}

func (a Activity) GetMainNodeProgress() decimal.Decimal {
	if a.MainNodeProgress == nil {
		return decimal.Zero
	}
	return *a.MainNodeProgress
}
