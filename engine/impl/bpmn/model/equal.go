package model

type IEqual interface {
	Equal(other interface{}) bool
}
