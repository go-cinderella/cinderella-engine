package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

// Float64 类型
type FloatLiteral struct {
	*Literal
}

func (l FloatLiteral) GetValueInternal(expressionState ExpressionState) TypedValue {
	return l.Value
}
