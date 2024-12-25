package engine

import (
	"github.com/go-cinderella/cinderella-engine/engine/idgenerator"
	"github.com/go-resty/resty/v2"
)

type ExpressionManagerFactory func() ExpressionManager

type ProcessEngineConfiguration interface {
	GetCommandExecutor() Executor
	GetIDGenerator() idgenerator.IDGenerator
	GetDeploymentSettings() map[string]interface{}
	GetBpmnDeployer() Deployer
	GetExpressionManagerFactory() ExpressionManagerFactory
	GetHttpClient() *resty.Client
}
