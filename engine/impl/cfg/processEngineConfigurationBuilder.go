package cfg

import (
	"github.com/go-cinderella/cinderella-engine/engine"
	"github.com/go-cinderella/cinderella-engine/engine/idgenerator"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/parse"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/parse/deployer"
	"github.com/go-cinderella/cinderella-engine/engine/impl/bpmn/parse/factory"
	"github.com/go-cinderella/cinderella-engine/engine/impl/expr"
	"github.com/go-cinderella/cinderella-engine/engine/impl/interceptor"
	"github.com/go-cinderella/cinderella-engine/engine/utils"
	"github.com/go-cinderella/cinderella-engine/engine/uuidgenerator"
	"github.com/go-resty/resty/v2"
)

type ProcessEngineConfigurationBuilder struct {
	commandInvoker        engine.Interceptor
	commandInterceptors   []engine.Interceptor
	commandExecutor       engine.Executor
	commandContextFactory engine.ICommandContextFactory

	bpmnParseHandlers []parse.BpmnParseHandler

	bpmnDeployer engine.Deployer

	idGenerator idgenerator.IDGenerator

	deploymentSettings map[string]interface{}

	instance *ProcessEngineConfigurationImpl

	expressionManagerFactory engine.ExpressionManagerFactory

	httpClient *resty.Client
}

func NewProcessEngineConfigurationBuilder() *ProcessEngineConfigurationBuilder {
	return &ProcessEngineConfigurationBuilder{
		instance: &ProcessEngineConfigurationImpl{},
	}
}

func (builder *ProcessEngineConfigurationBuilder) ExpressionManagerFactory(expressionManagerFactory engine.ExpressionManagerFactory) *ProcessEngineConfigurationBuilder {
	builder.expressionManagerFactory = expressionManagerFactory
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) IdGenerator(idGenerator idgenerator.IDGenerator) *ProcessEngineConfigurationBuilder {
	builder.idGenerator = idGenerator
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) CommandContextFactory(commandContextFactory engine.ICommandContextFactory) *ProcessEngineConfigurationBuilder {
	builder.commandContextFactory = commandContextFactory
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) CommandInvoker(commandInvoker engine.Interceptor) *ProcessEngineConfigurationBuilder {
	builder.commandInvoker = commandInvoker
	return builder
}

// CommandInterceptor custom command interceptor
func (builder *ProcessEngineConfigurationBuilder) CommandInterceptor(interceptor engine.Interceptor) *ProcessEngineConfigurationBuilder {
	builder.commandInterceptors = append(builder.commandInterceptors, interceptor)
	return builder
}

// CommandInterceptors custom command interceptors
func (builder *ProcessEngineConfigurationBuilder) CommandInterceptors(interceptor ...engine.Interceptor) *ProcessEngineConfigurationBuilder {
	builder.commandInterceptors = append(builder.commandInterceptors, interceptor...)
	return builder
}

// BpmnParseHandler custom bpmn parse handler
func (builder *ProcessEngineConfigurationBuilder) BpmnParseHandler(bpmnParseHandler parse.BpmnParseHandler) *ProcessEngineConfigurationBuilder {
	builder.bpmnParseHandlers = append(builder.bpmnParseHandlers, bpmnParseHandler)
	return builder
}

// BpmnParseHandlers custom bpmn parse handlers
func (builder *ProcessEngineConfigurationBuilder) BpmnParseHandlers(bpmnParseHandler ...parse.BpmnParseHandler) *ProcessEngineConfigurationBuilder {
	builder.bpmnParseHandlers = append(builder.bpmnParseHandlers, bpmnParseHandler...)
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) CommandExecutor(commandExecutor engine.Executor) *ProcessEngineConfigurationBuilder {
	builder.commandExecutor = commandExecutor
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) BpmnDeployer(bpmnDeployer engine.Deployer) *ProcessEngineConfigurationBuilder {
	builder.bpmnDeployer = bpmnDeployer
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) DeploymentSettings(deploymentSettings map[string]interface{}) *ProcessEngineConfigurationBuilder {
	builder.deploymentSettings = deploymentSettings
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) HttpClient(httpClient *resty.Client) *ProcessEngineConfigurationBuilder {
	builder.httpClient = httpClient
	return builder
}

func (builder *ProcessEngineConfigurationBuilder) Get() *ProcessEngineConfigurationImpl {
	return builder.instance
}

func (builder *ProcessEngineConfigurationBuilder) Build() *ProcessEngineConfigurationImpl {
	result := builder.instance
	result.commandInterceptors = builder.commandInterceptors
	result.commandExecutor = builder.commandExecutor
	result.commandContextFactory = builder.commandContextFactory
	result.bpmnParseHandlers = builder.bpmnParseHandlers
	result.bpmnDeployer = builder.bpmnDeployer
	result.idGenerator = builder.idGenerator
	result.deploymentSettings = builder.deploymentSettings
	result.expressionManagerFactory = builder.expressionManagerFactory
	result.httpClient = builder.httpClient

	if result.idGenerator == nil {
		result.idGenerator = uuidgenerator.UUIDGenerator{}
	}

	if result.commandContextFactory == nil {
		result.commandContextFactory = interceptor.CommandContextFactory{
			ProcessEngineConfiguration: result,
		}
	}

	var interceptors []engine.Interceptor
	interceptors = append(interceptors, &interceptor.CommandContextInterceptor{CommandContextFactory: result.commandContextFactory})
	interceptors = append(interceptors, &interceptor.TransactionContextInterceptor{})

	if len(result.commandInterceptors) > 0 {
		interceptors = append(interceptors, result.commandInterceptors...)
	}

	interceptors = append(interceptors, &interceptor.CommandInvoker{})
	result.commandInterceptors = interceptors

	if result.commandExecutor == nil {
		first := initInterceptorChain(result.commandInterceptors)
		result.commandExecutor = CommandExecutorImpl{First: first}
	}

	if result.bpmnDeployer == nil {
		bpmnParser := parse.BpmnParser{}
		bpmnParser.BpmnParseFactory = parse.DefaultBpmnParseFactory{}
		bpmnParser.ActivityBehaviorFactory = factory.DefaultActivityBehaviorFactory{}

		parseHandlers := parse.BpmnParseHandlers{ParseHandlers: make(map[string][]parse.BpmnParseHandler)}
		parseHandlers.AddHandlers(getDefaultBpmnParseHandlers())

		if len(result.bpmnParseHandlers) > 0 {
			parseHandlers.AddHandlers(result.bpmnParseHandlers)
		}

		bpmnParser.BpmnParserHandlers = parseHandlers

		bpmnDeployer := deployer.BpmnDeployer{
			ParsedDeploymentBuilderFactory: parse.ParsedDeploymentBuilderFactory{
				BpmnParser: &bpmnParser,
			},
		}
		result.bpmnDeployer = &bpmnDeployer
	}

	if result.expressionManagerFactory == nil {
		result.expressionManagerFactory = func() engine.ExpressionManager {
			return expr.DefaultExpressionManager{}
		}
	}

	if result.httpClient == nil {
		result.httpClient = utils.NewDefaultHttpClient()
	}

	return result
}
