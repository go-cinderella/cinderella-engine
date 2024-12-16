package expression

import "github.com/go-cinderella/cinderella-engine/engine/expression/spel"

var _ Expression = (*CompositeStringExpression)(nil)

type CompositeStringExpression struct {
	ExpressionString string
	Expressions      []Expression
}

func (c *CompositeStringExpression) GetValueContext(context spel.EvaluationContext) interface{} {
	return c.GetValue()
}

func (c *CompositeStringExpression) GetExpressionString() string {
	return c.ExpressionString
}

func (c *CompositeStringExpression) GetValue() interface{} {
	//s := ""

	return "c.literalValue"
}
