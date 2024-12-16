package utils

import (
	"github.com/go-cinderella/cinderella-engine/engine/expression"
	"github.com/go-cinderella/cinderella-engine/engine/expression/spel"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
)

var (
	context = spel.StandardEvaluationContext{}
	parser  = expression.SpelExpressionParser{}
)

type ConditionUtil struct {
}

func HasTrueCondition(sequenceFlow model.SequenceFlow, execution delegate.DelegateExecution) bool {
	var conditionExpression = sequenceFlow.ConditionExpression
	if conditionExpression != "" {
		variable := execution.GetProcessVariable()
		context.SetVariables(variable)
		valueContext := parser.ParseExpression(conditionExpression).GetValueContext(&context)
		b, ok := valueContext.(bool)
		if ok {
			return b
		}
		return false
	} else {
		return true
	}

}
