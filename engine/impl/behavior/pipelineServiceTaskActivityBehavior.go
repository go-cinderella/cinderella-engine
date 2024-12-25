package behavior

import (
	"context"
	"encoding/json"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"strings"
)

var _ delegate.ActivityBehavior = (*PipelineServiceTaskActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*PipelineServiceTaskActivityBehavior)(nil)

type PipelineServiceTaskActivityBehavior struct {
	abstractBpmnActivityBehavior
	ServiceTask model.ServiceTask
	ProcessKey  string
}

func (pipeline PipelineServiceTaskActivityBehavior) Execute(execution delegate.DelegateExecution) error {
	httpClient := contextutil.GetHttpClient()
	req := httpClient.NewRequest()

	// set context
	ctx := db.DB().Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}
	req.SetContext(ctx)

	// set request body
	requestBody := make(map[string]interface{})
	extensionElements := pipeline.ServiceTask.ExtensionElements

	field := extensionElements.GetFieldByName("requestBody")
	if stringutils.IsNotEmpty(field.FieldName) {
		fieldValue := make(map[string]interface{})
		if err := json.Unmarshal([]byte(field.Expression), &fieldValue); err != nil {
			return errors.WithStack(err)
		}
		for k, v := range fieldValue {
			if v == nil {
				continue
			}

			parameterValue := cast.ToString(v)
			if strings.HasPrefix(parameterValue, "${") && strings.HasSuffix(parameterValue, "}") {
				expressionManager := contextutil.GetExpressionManager()
				context := expressionManager.EvaluationContext()

				variable := execution.GetProcessVariable()
				if len(variable) > 0 {
					context.SetVariables(variable)
				}

				expression := expressionManager.CreateExpression(parameterValue)
				value := expression.GetValueContext(&context)
				b, ok := value.(string)
				if ok && stringutils.IsNotEmpty(b) {
					requestBody[k] = b
				}
			} else {
				requestBody[k] = v
			}
		}
	}

	req.SetBody(requestBody)

	// set request url
	field = extensionElements.GetFieldByName("requestUrl")
	if stringutils.IsEmpty(field.FieldName) {
		return errors.WithStack(errors.New("requestUrl is empty"))
	}

	requestUrl := field.StringValue
	resp, err := req.Post(requestUrl)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.IsError() {
		return errors.WithStack(errors.New(string(resp.Body())))
	}

	/**
	{
	   "code": "-1",
	   "message": "不支持的Git触发类型",
	   "data": null
	}

	*/
	result := make(map[string]interface{})
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return errors.WithStack(err)
	}

	if cast.ToString(result["code"]) != "200" {
		return errors.WithStack(errors.New(result["message"].(string)))
	}

	pipeline.leave(execution)
	return nil
}
