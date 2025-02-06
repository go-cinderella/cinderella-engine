package model

type LoopCharacteristicsGetter interface {
	GetLoopCharacteristics() *MultiInstanceLoopCharacteristics
}
