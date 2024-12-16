package engine

import (
	"github.com/go-cinderella/cinderella-engine/engine/expression"
	"github.com/go-cinderella/cinderella-engine/engine/expression/spel"
)

type ExpressionManager interface {
	CreateExpression(expression string) expression.Expression
	EvaluationContext() spel.StandardEvaluationContext
}
