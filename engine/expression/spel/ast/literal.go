package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

type GetLiteralValue interface {
	GetLiteralValue() TypedValue
}

type Literal struct {
	*SpelNodeImpl
	OriginalValue string
	Value         TypedValue
}
