package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

// map参数
const (
	THIS = "this"
	ROOT = "root"
)

type VariableReference struct {
	*SpelNodeImpl
	Name string
}

func (v VariableReference) GetValueInternal(state ExpressionState) TypedValue {
	return state.LookupVariable(v.Name)
}
