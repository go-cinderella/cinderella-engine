package expression

import "github.com/go-cinderella/cinderella-engine/engine/expression/spel"

// 获取值
type Expression interface {
	GetExpressionString() string

	GetValue() interface{}

	GetValueContext(context spel.EvaluationContext) interface{}
}
