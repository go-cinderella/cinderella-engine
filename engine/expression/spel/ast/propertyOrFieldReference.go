package ast

import (
	. "github.com/go-cinderella/cinderella-engine/engine/expression/spel"
	"reflect"
)

// 对象属性处理
type PropertyOrFieldReference struct {
	*SpelNodeImpl
	NullSafe                            bool
	Name                                string
	OriginalPrimitiveExitTypeDescriptor string
	CachedReadAccessor                  PropertyAccessor
}

type ValueRef interface {
	GetValue() TypedValue
}

type AccessorLValue struct {
	ref                    PropertyOrFieldReference
	contextObject          TypedValue
	evalContext            EvaluationContext
	autoGrowNullReferences bool
}

func (this AccessorLValue) GetValue() TypedValue {
	return this.ref.getValueInternal(this.contextObject, this.evalContext, this.autoGrowNullReferences)
}

func (this PropertyOrFieldReference) GetValueRef(state ExpressionState) ValueRef {
	return AccessorLValue{ref: this, contextObject: state.GetActiveContextObject(), evalContext: state.RelatedContext,
		autoGrowNullReferences: false}
}

func (this PropertyOrFieldReference) GetValueInternal(state ExpressionState) TypedValue {
	return this.getValueInternal(state.GetActiveContextObject(), state.GetEvaluationContext(), false)
}

func (this PropertyOrFieldReference) getValueInternal(contextObject TypedValue, evalContext EvaluationContext, isAutoGrowNullReferences bool) TypedValue {
	return this.readProperty(contextObject, evalContext, this.Name)
}

func (this PropertyOrFieldReference) readProperty(contextObject TypedValue, evalContext EvaluationContext, name string) TypedValue {
	accessors := getPropertyAccessorsToTry(contextObject.Value, evalContext.GetPropertyAccessors())
	for _, accessor := range accessors {
		this.CachedReadAccessor = accessor
		return accessor.Read(evalContext, contextObject.Value, name)
	}
	return TypedValue{}
}

func getPropertyAccessorsToTry(contextObject interface{}, propertyAccessors []PropertyAccessor) []PropertyAccessor {
	var targetType reflect.Kind
	if contextObject != nil {
		targetType = reflect.TypeOf(contextObject).Kind()
	}
	var specificAccessors []PropertyAccessor
	var generalAccessors []PropertyAccessor
	for _, accessor := range propertyAccessors {
		classes := accessor.GetSpecificTargetClasses()
		if classes == nil {
			generalAccessors = append(generalAccessors, accessor)
		} else {
			//是否是子类
			if classes == targetType {
				generalAccessors = append(specificAccessors, accessor)
			} else {
				specificAccessors = append(generalAccessors, accessor)
			}
		}
	}
	var resolvers []PropertyAccessor
	resolvers = generalAccessors
	return resolvers
}
