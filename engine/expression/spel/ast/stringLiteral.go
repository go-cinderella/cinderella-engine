package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

/*
*
 */
type StringLiteral struct {
	*Literal
}

func (l StringLiteral) GetValueInternal(expressionState ExpressionState) TypedValue {
	return l.Value
}
