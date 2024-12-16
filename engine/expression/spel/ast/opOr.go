package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
	. "github.com/go-cinderella/cinderella-engine/engine/expression/support"
)

type OpOr struct {
	*Operator
}

func (o *OpOr) GetValueInternal(expressionState ExpressionState) TypedValue {
	if getBooleanValue(expressionState, o.getLeftOperand()) {
		value := BooleanTypedValue{}
		return value.ForValue(true)
	}
	booleanValue := getBooleanValue(expressionState, o.getRightOperand())
	value := BooleanTypedValue{}
	return value.ForValue(booleanValue)
}
