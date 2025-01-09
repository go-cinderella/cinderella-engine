package utils

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/expression"
	"github.com/go-cinderella/cinderella-engine/engine/expression/spel"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/spf13/cast"
	"maps"
	"strings"
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

func IsExpr(input string) bool {
	return strings.HasPrefix(input, "${") && strings.HasSuffix(input, "}")
}

func GetStringSliceFromExpression(variables map[string]interface{}, input string) []string {
	expressionManager := contextutil.GetExpressionManager()
	context := expressionManager.EvaluationContext()

	if len(variables) > 0 {
		context.SetVariables(maps.Clone(variables))
	}

	expression := expressionManager.CreateExpression(input)
	value := expression.GetValueContext(&context)

	b := cast.ToString(value)
	return strings.Split(b, ",")
}
