package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
	"github.com/go-cinderella/cinderella-engine/engine/expression/support"
)

// 处理大于等于
type OpGE struct {
	*Operator
}

func (o *OpGE) GetValueInternal(expressionState ExpressionState) TypedValue {
	value := support.BooleanTypedValue{}
	left := o.getLeftOperand().GetValueInternal(expressionState).Value
	right := o.getRightOperand().GetValueInternal(expressionState).Value
	checkType := checkType(left, right)
	if !checkType {
		return value.ForValue(checkType)
	}
	o.leftActualDescriptor = o.toDescriptorFromObject(left)
	o.rightActualDescriptor = o.toDescriptorFromObject(right)
	var check bool
	leftV, ok := left.(int)
	if ok {
		rightV := right.(int)
		check = leftV >= rightV
	} else {
		leftV, ok := left.(float64)
		if ok {
			rightV := right.(float64)
			check = leftV >= rightV
		}
	}
	return value.ForValue(check)
}
