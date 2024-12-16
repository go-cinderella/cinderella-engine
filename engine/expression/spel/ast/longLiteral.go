package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

type LongLiteral struct {
	*Literal
}

func (l LongLiteral) GetValueInternal(expressionState ExpressionState) TypedValue {
	return l.Value
}
