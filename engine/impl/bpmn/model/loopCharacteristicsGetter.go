package model

type LoopCharacteristicsGetter interface {
	GetLoopCharacteristics() *MultiInstanceLoopCharacteristics
	HasMultiInstanceLoopCharacteristics() bool
}
