package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

// 表达式对象
type SpelNode interface {
	GetValue(expressionState ExpressionState) interface{}

	GetValueInternal(expressionState ExpressionState) TypedValue

	GetValueRef(state ExpressionState) ValueRef

	GetStartPosition() int

	GetEndPosition() int
}
