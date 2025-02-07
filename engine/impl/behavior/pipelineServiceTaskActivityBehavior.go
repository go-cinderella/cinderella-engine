package behavior

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/go-cinderella/cinderella-engine/engine/contextutil"
	"github.com/go-cinderella/cinderella-engine/engine/db"
	"github.com/go-cinderella/cinderella-engine/engine/entitymanager"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/model"
	"github.com/go-cinderella/cinderella-engine/engine/impl/delegate"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/unionj-cloud/toolkit/stringutils"
	"maps"
	"regexp"
	"strings"
)

var json = sonic.ConfigDefault

var _ delegate.ActivityBehavior = (*PipelineServiceTaskActivityBehavior)(nil)
var _ delegate.TriggerableActivityBehavior = (*PipelineServiceTaskActivityBehavior)(nil)

type PipelineServiceTaskActivityBehavior struct {
	abstractBpmnActivityBehavior
	ServiceTask model.ServiceTask
	ProcessKey  string
}

var rgx = regexp.MustCompile(`\$\{(.*?)\}`)

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
	variables, err := execution.GetVariables()
	if err != nil {
		return err
	}
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
			if utils.IsExpr(parameterValue) {
				output := utils.GetStringFromExpression(variables, parameterValue)
				if stringutils.IsNotEmpty(output) {
					requestBody[k] = output
				}
			} else {
				requestBody[k] = v
			}
		}
	}

	businessParameter := make(map[string]any)
	businessParameter["requestBody"] = maps.Clone(requestBody)

	requestBody["processInstanceId"] = execution.GetProcessInstanceId()
	requestBody["executionId"] = execution.GetExecutionId()

	req.SetBody(requestBody)

	// set request url
	field = extensionElements.GetFieldByName("requestUrl")
	if stringutils.IsEmpty(field.FieldName) {
		return errors.WithStack(errors.New("requestUrl is empty"))
	}

	requestUrl := field.StringValue

	if rgx.MatchString(requestUrl) {
		requestUrl = rgx.ReplaceAllStringFunc(requestUrl, func(s string) string {
			s = strings.TrimPrefix(s, "${")
			s = strings.TrimSuffix(s, "}")
			v, ok := variables[s]
			if !ok {
				return s
			}
			return cast.ToString(v)
		})
	}

	businessParameter["requestUrl"] = requestUrl

	resp, err := req.Post(requestUrl)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.IsError() {
		return errors.WithStack(errors.New(string(resp.Body())))
	}

	result := make(map[string]interface{})
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return errors.WithStack(err)
	}

	if cast.ToString(result["code"]) != "200" {
		return errors.WithStack(errors.New(cast.ToString(result["message"])))
	}

	businessResult := cast.ToString(result["data"])
	businessParameterJson, _ := json.MarshalToString(businessParameter)

	historicActivityInstanceEntityManager := entitymanager.GetHistoricActivityInstanceEntityManager()
	if err = historicActivityInstanceEntityManager.RecordBusinessDataByExecutionId(execution, businessParameterJson, businessResult); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Trigger 普通用户节点处理
func (pipeline PipelineServiceTaskActivityBehavior) Trigger(execution delegate.DelegateExecution) error {
	pipeline.leave(execution)
	return nil
}
