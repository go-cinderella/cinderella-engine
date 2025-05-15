package model

import (
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/shopspring/decimal"
	"reflect"
)

var _ delegate.MainNodeSupported = (*Activity)(nil)

type Activity struct {
	FlowNode
	IsMainNode        *bool                             `xml:"isMainNode,attr"`
	MainNodeWeight    *decimal.Decimal                  `xml:"mainNodeWeight,attr"`
	MainNodeIndex     *int                              `xml:"mainNodeIndex,attr"`
	MainNodeProgress  *decimal.Decimal                  `xml:"mainNodeProgress,attr"`
	MultiInstance     *MultiInstanceLoopCharacteristics `xml:"multiInstanceLoopCharacteristics"`
	ExtensionElements *ExtensionElement                 `xml:"extensionElements"`
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

func (a *Activity) ActivityEqual(otherActivity interface{}) bool {
	other, ok := otherActivity.(*Activity)
	if !ok {
		that2, ok := otherActivity.(Activity)
		if ok {
			other = &that2
		} else {
			return false
		}
	}

	if other == nil {
		return a == nil
	} else if a == nil {
		return false
	}

	if a.DefaultBaseElement != other.DefaultBaseElement {
		return false
	}

	if !reflect.DeepEqual(a.MultiInstance, other.MultiInstance) {
		return false
	}

	if a.GetIsMainNode() != other.GetIsMainNode() {
		return false
	}

	if !a.GetMainNodeWeight().Equal(other.GetMainNodeWeight()) {
		return false
	}

	if a.GetMainNodeIndex() != other.GetMainNodeIndex() {
		return false
	}

	if !a.GetMainNodeProgress().Equal(other.GetMainNodeProgress()) {
		return false
	}

	if !reflect.DeepEqual(a.ExtensionElements, other.ExtensionElements) {
		return false
	}

	return true
}
