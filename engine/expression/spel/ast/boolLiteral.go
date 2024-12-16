package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

type BoolLiteral struct {
	*Literal
}

func (l BoolLiteral) GetValueInternal(expressionState ExpressionState) TypedValue {
	return l.Value
}
