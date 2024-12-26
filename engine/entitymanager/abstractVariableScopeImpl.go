package entitymanager

import (
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/variable"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"reflect"
	"sync"
)

var usedVariablesCache *sync.Pool

type AbstractVariableScopeImpl interface {
	GetSourceActivityExecution() ExecutionEntity
}

func init() {
	usedVariablesCache = &sync.Pool{
		New: func() interface{} {
			return nil
		},
	}
}

type VariableScopeImpl struct {
	AbstractVariableScopeImpl
	ExecutionEntity ExecutionEntity
}

func IsIntegral(val float64) bool {
	return val == float64(int(val))
}

// SetVariable 保存流程变量
func (variableScope VariableScopeImpl) SetVariable(execution delegate.DelegateExecution, variables map[string]interface{}) error {
	variableManager := variable.GetVariableManager()
	variableTypes := variableManager.VariableTypes
	idGenerator := contextutil.GetIDGenerator()

	for k, v := range variables {
		if v == nil {
			continue
		}

		kind := reflect.TypeOf(v).Kind()
		if kind == reflect.Float64 {
			if IsIntegral(v.(float64)) {
				v = cast.ToInt(v)
				kind = reflect.Int
			}
		}

		variableType := variableTypes.GetVariableType(kind.String())
		if variableType == nil {
			continue
		}

		vari := variable.Variable{}
		vari.ID_, _ = idGenerator.NextID()
		vari.Name_ = k
		vari.Type_ = variableType.GetTypeName()
		vari.SetValue(v, variableType)
		vari.ProcInstID_ = lo.ToPtr(execution.GetProcessInstanceId())
		vari.ExecutionID_ = lo.ToPtr(execution.GetProcessInstanceId())

		if err := variableEntityManager.UpsertVariable(vari); err != nil {
			return err
		}
	}

	return nil
}
func (variableScope VariableScopeImpl) SetVariableLocal(variables map[string]interface{}) error {
	variable := usedVariablesCache.Get()
	if variable == nil {
		usedVariablesCache.Put(variables)
	} else {
		variablesCache := variable.(map[string]interface{})
		for k, v := range variables {
			variablesCache[k] = v
		}
		usedVariablesCache.Put(variablesCache)
	}
	return nil
}

func (variableScope VariableScopeImpl) GetVariableLocal() (variables map[string]interface{}) {
	variable := usedVariablesCache.Get()
	if variable == nil {
		return nil
	} else {
		variablesCache := variable.(map[string]interface{})
		return variablesCache
	}
	return nil
}
