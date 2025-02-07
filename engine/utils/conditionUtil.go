package utils

import (
	"errors"
	"github.com/expr-lang/expr"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/unionj-cloud/toolkit/zlogger"
	"maps"
	"strings"
)

type ConditionUtil struct {
}

func HasTrueCondition(sequenceFlow model.SequenceFlow, execution delegate.DelegateExecution) bool {
	var conditionExpression = sequenceFlow.ConditionExpression
	if stringutils.IsNotEmpty(conditionExpression) && IsExpr(conditionExpression) {
		code := Trim(conditionExpression)
		variables, err := execution.GetVariables()
		if err != nil {
			zlogger.Error().Err(err).Msg("failed to compile condition expression")
			return false
		}

		program, err := expr.Compile(code, expr.Env(variables))
		if err != nil {
			zlogger.Error().Err(err).Msg("failed to compile condition expression")
			return false
		}

		output, err := expr.Run(program, variables)
		if err != nil {
			zlogger.Error().Err(err).Msg("failed to evaluate condition expression")
			return false
		}

		b, ok := output.(bool)
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

func IsTrue(variables map[string]interface{}, input string) bool {
	output, err := evaluate(variables, input)
	if err != nil {
		zlogger.Error().Err(err).Msg("failed to evaluate condition expression")
		return false
	}

	return cast.ToBool(output)
}

func Trim(input string) string {
	return strings.TrimSuffix(strings.TrimPrefix(input, "${"), "}")
}

func evaluate(variables map[string]interface{}, input string) (interface{}, error) {
	if !IsExpr(input) {
		return nil, errors.New(`not a valid expression`)
	}
	env := maps.Clone(variables)
	program, err := expr.Compile(Trim(input), expr.Env(env))
	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func GetStringSliceFromExpression(variables map[string]interface{}, input string) []string {
	output, err := evaluate(variables, input)
	if err != nil {
		zlogger.Error().Err(err).Msg("failed to evaluate condition expression")
		return nil
	}

	b := cast.ToString(output)
	return stringutils.Split(b, ",")
}

func GetStringFromExpression(variables map[string]interface{}, input string) string {
	result := GetStringSliceFromExpression(variables, input)
	if len(result) == 0 {
		return ""
	}
	return result[0]
}
