package cfg

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/idgenerator"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/parse"
	"github.com/go-resty/resty/v2"
)

var _ engine.ProcessEngineConfiguration = (*ProcessEngineConfigurationImpl)(nil)

type ProcessEngineConfigurationImpl struct {
	commandInterceptors   []engine.Interceptor
	commandExecutor       engine.Executor
	commandContextFactory engine.ICommandContextFactory
	bpmnParseHandlers     []parse.BpmnParseHandler

	bpmnDeployer engine.Deployer

	idGenerator idgenerator.IDGenerator

	deploymentSettings       map[string]interface{}
	expressionManagerFactory engine.ExpressionManagerFactory

	httpClient *resty.Client
}

func (processEngineConfiguration ProcessEngineConfigurationImpl) GetHttpClient() *resty.Client {
	return processEngineConfiguration.httpClient
}

func (processEngineConfiguration ProcessEngineConfigurationImpl) GetExpressionManagerFactory() engine.ExpressionManagerFactory {
	return processEngineConfiguration.expressionManagerFactory
}

func (processEngineConfiguration ProcessEngineConfigurationImpl) GetDeploymentSettings() map[string]interface{} {
	return processEngineConfiguration.deploymentSettings
}

func (processEngineConfiguration ProcessEngineConfigurationImpl) GetBpmnDeployer() engine.Deployer {
	return processEngineConfiguration.bpmnDeployer
}

func (processEngineConfiguration ProcessEngineConfigurationImpl) GetIDGenerator() idgenerator.IDGenerator {
	return processEngineConfiguration.idGenerator
}

func (processEngineConfiguration ProcessEngineConfigurationImpl) GetCommandExecutor() engine.Executor {
	return processEngineConfiguration.commandExecutor
}

func initInterceptorChain(interceptors []engine.Interceptor) engine.Interceptor {
	if len(interceptors) > 0 {
		for i := 0; i < len(interceptors)-1; i++ {
			interceptor := interceptors[i]
			interceptor.SetNext(interceptors[i+1])
		}
	}
	return interceptors[0]
}

func getDefaultBpmnParseHandlers() []parse.BpmnParseHandler {
	handlers := make([]parse.BpmnParseHandler, 0)

	handlers = append(handlers, parse.ProcessParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ParseHandler(parse.ProcessParseHandler{})}}})

	handlers = append(handlers, parse.StartEventParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ParseHandler(parse.StartEventParseHandler{})}}})

	handlers = append(handlers, parse.UserTaskParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ParseHandler(parse.UserTaskParseHandler{})}}})

	handlers = append(handlers, parse.SequenceFlowParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ParseHandler(parse.SequenceFlowParseHandler{})}}})

	handlers = append(handlers, parse.ExclusiveGatewayParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ParseHandler(parse.ExclusiveGatewayParseHandler{})}}})

	handlers = append(handlers, parse.InclusiveGatewayParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ParseHandler(parse.InclusiveGatewayParseHandler{})}}})

	handlers = append(handlers, parse.ParallelGatewayParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ParseHandler(parse.ParallelGatewayParseHandler{})}}})

	handlers = append(handlers, parse.IntermediateCatchEventParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.IntermediateCatchEventParseHandler{}}}})

	handlers = append(handlers, parse.ConditionalEventDefinitionParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ConditionalEventDefinitionParseHandler{}}}})

	handlers = append(handlers, parse.ServiceTaskParseHandler{AbstractActivityBpmnParseHandler: parse.AbstractActivityBpmnParseHandler{AbstractBpmnParseHandler: parse.AbstractBpmnParseHandler{ParseHandler: parse.ServiceTaskParseHandler{}}}})

	return handlers
}
