package expr

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/expression"
	"github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

var _ engine.ExpressionManager = (*DefaultExpressionManager)(nil)

type DefaultExpressionManager struct {
	parser            expression.SpelExpressionParser
	evaluationContext spel.StandardEvaluationContext
}

func (e DefaultExpressionManager) EvaluationContext() spel.StandardEvaluationContext {
	return e.evaluationContext
}

func (e *DefaultExpressionManager) SetEvaluationContext(evaluationContext spel.StandardEvaluationContext) {
	e.evaluationContext = evaluationContext
}

func (e DefaultExpressionManager) CreateExpression(expression string) expression.Expression {
	return e.parser.ParseExpression(expression)
}
