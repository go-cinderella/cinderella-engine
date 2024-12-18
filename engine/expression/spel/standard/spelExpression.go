package standard

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel/ast"
)

type SpelExpression struct {
	Expression        string
	Configuration     SpelParserConfiguration
	EvaluationContext EvaluationContext
	Ast               SpelNode
}

func (e SpelExpression) GetExpressionString() string {
	return e.Expression
}

func (e *SpelExpression) GetValue() interface{} {
	context := e.getEvaluationContext()
	state := ExpressionState{RelatedContext: context}
	return e.Ast.GetValue(state)
}

func (e *SpelExpression) GetValueContext(context EvaluationContext) interface{} {
	if context == nil {
		panic("context is required")
	}
	state := ExpressionState{RelatedContext: context}
	typedResultValue := e.Ast.GetValueInternal(state)
	return typedResultValue.Value
}

func (e *SpelExpression) getEvaluationContext() EvaluationContext {
	if e.EvaluationContext == nil {
		context := StandardEvaluationContext{}
		e.EvaluationContext = &context
	}
	return e.EvaluationContext
}
