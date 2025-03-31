package delegate

import "github.com/shopspring/decimal"

type MainNodeSupported interface {
	GetIsMainNode() bool
	GetMainNodeProgress() decimal.Decimal
	GetMainNodeWeight() decimal.Decimal
	GetMainNodeIndex() int
}
