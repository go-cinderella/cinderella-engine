package ast

import . "github.com/go-cinderella/cinderella-engine/engine/expression/spel"

type LiteralValue struct {
	*SpelNodeImpl
}

func (l *LiteralValue) GetLiteralValue() TypedValue {
	return TypedValue{}
}
