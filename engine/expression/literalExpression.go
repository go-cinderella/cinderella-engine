package expression

import "github.com/go-cinderella/cinderella-engine/engine/expression/spel"

var _ Expression = (*LiteralExpression)(nil)

type LiteralExpression struct {
	literalValue string
}

func (l *LiteralExpression) GetValueContext(context spel.EvaluationContext) interface{} {
	return l.GetValue()
}

func (l *LiteralExpression) GetExpressionString() string {
	return l.literalValue
}

func (l *LiteralExpression) GetValue() interface{} {
	return l.literalValue
}
