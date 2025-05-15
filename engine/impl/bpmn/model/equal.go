package model

type IEqual interface {
	ActivityEqual(other interface{}) bool
}
